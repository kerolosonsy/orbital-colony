package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

const refPrefix = "#/components/schemas/"

type schema struct {
	Ref                  string             `json:"$ref"`
	Type                 string             `json:"type"`
	Required             []string           `json:"required"`
	Properties           map[string]*schema `json:"properties"`
	AdditionalProperties *bool              `json:"additionalProperties"`
	Items                *schema            `json:"items"`
	Enum                 []any              `json:"enum"`
	Minimum              *float64           `json:"minimum"`
	MinLength            *int               `json:"minLength"`
}

type mediaType struct {
	Schema *schema `json:"schema"`
}

type response struct {
	Content map[string]mediaType `json:"content"`
}

type operation struct {
	Responses map[string]response `json:"responses"`
}

type pathItem struct {
	Get    *operation `json:"get"`
	Put    *operation `json:"put"`
	Post   *operation `json:"post"`
	Delete *operation `json:"delete"`
	Patch  *operation `json:"patch"`
}

type document struct {
	OpenAPI    string              `json:"openapi"`
	Paths      map[string]pathItem `json:"paths"`
	Components struct {
		Schemas map[string]*schema `json:"schemas"`
	} `json:"components"`
}

type Spec struct {
	ops     map[string]map[string]*operation // template -> method -> operation
	schemas map[string]*schema
}

func (s *Spec) checkRefs(sc *schema) error {
	if sc == nil {
		return nil
	}
	if sc.Ref != "" {
		if !strings.HasPrefix(sc.Ref, refPrefix) || s.schemas[strings.TrimPrefix(sc.Ref, refPrefix)] == nil {
			return fmt.Errorf("unresolved $ref %q", sc.Ref)
		}
	}
	for _, p := range sc.Properties {
		if err := s.checkRefs(p); err != nil {
			return err
		}
	}
	return s.checkRefs(sc.Items)
}

func LoadSpec(raw []byte) (*Spec, error) {
	var doc document
	if err := json.Unmarshal(raw, &doc); err != nil {
		return nil, err
	}
	if !strings.HasPrefix(doc.OpenAPI, "3.") {
		return nil, errors.New("not an OpenAPI 3 document")
	}
	if len(doc.Paths) == 0 {
		return nil, errors.New("no paths")
	}
	s := &Spec{ops: map[string]map[string]*operation{}, schemas: doc.Components.Schemas}
	for _, sc := range s.schemas {
		if err := s.checkRefs(sc); err != nil {
			return nil, err
		}
	}
	for path, item := range doc.Paths {
		if !strings.HasPrefix(path, "/") {
			return nil, fmt.Errorf("path %q must start with /", path)
		}
		ops := map[string]*operation{}
		for method, op := range map[string]*operation{"get": item.Get, "put": item.Put, "post": item.Post, "delete": item.Delete, "patch": item.Patch} {
			if op == nil {
				continue
			}
			if len(op.Responses) == 0 {
				return nil, fmt.Errorf("%s %s has no responses", method, path)
			}
			for code, resp := range op.Responses {
				n, err := strconv.Atoi(code)
				if err != nil || len(code) != 3 || n < 100 || n > 599 {
					return nil, fmt.Errorf("%s %s: bad status %q", method, path, code)
				}
				for _, m := range resp.Content {
					if err := s.checkRefs(m.Schema); err != nil {
						return nil, err
					}
				}
			}
			ops[method] = op
		}
		s.ops[path] = ops
	}
	return s, nil
}

func templateMatches(template, path string) bool {
	want, got := strings.Split(template, "/"), strings.Split(path, "/")
	if len(want) != len(got) {
		return false
	}
	for i, part := range want {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			if got[i] == "" {
				return false
			}
		} else if part != got[i] {
			return false
		}
	}
	return true
}

func (s *Spec) Check(method, path string, status int, contentType string, body []byte) error {
	var ops map[string]*operation
	for template, o := range s.ops {
		if templateMatches(template, path) {
			ops = o
			break
		}
	}
	if ops == nil {
		return fmt.Errorf("undocumented path %q", path)
	}
	op := ops[strings.ToLower(method)]
	if op == nil {
		return fmt.Errorf("undocumented method %s %s", method, path)
	}
	resp, ok := op.Responses[strconv.Itoa(status)]
	if !ok {
		return fmt.Errorf("undocumented status %d for %s %s", status, method, path)
	}
	if len(resp.Content) == 0 {
		if len(body) != 0 {
			return errors.New("documented without content but the body is not empty")
		}
		return nil
	}
	media, _, err := mime.ParseMediaType(contentType)
	if err != nil || media != "application/json" {
		return fmt.Errorf("content type %q", contentType)
	}
	m, ok := resp.Content[media]
	if !ok {
		return fmt.Errorf("content type %q not documented", media)
	}
	var v any
	if err := json.Unmarshal(body, &v); err != nil {
		return err
	}
	return s.validate(m.Schema, v, "body")
}

func (s *Spec) validate(sc *schema, v any, at string) error {
	if sc == nil {
		return nil
	}
	if sc.Ref != "" {
		return s.validate(s.schemas[strings.TrimPrefix(sc.Ref, refPrefix)], v, at)
	}
	switch sc.Type {
	case "object":
		obj, ok := v.(map[string]any)
		if !ok {
			return fmt.Errorf("%s: want object", at)
		}
		for _, key := range sc.Required {
			if _, ok := obj[key]; !ok {
				return fmt.Errorf("%s: missing %q", at, key)
			}
		}
		for key, value := range obj {
			child, known := sc.Properties[key]
			if !known && sc.AdditionalProperties != nil && !*sc.AdditionalProperties {
				return fmt.Errorf("%s: unexpected %q", at, key)
			}
			if err := s.validate(child, value, at+"."+key); err != nil {
				return err
			}
		}
	case "array":
		items, ok := v.([]any)
		if !ok {
			return fmt.Errorf("%s: want array", at)
		}
		for i, item := range items {
			if err := s.validate(sc.Items, item, fmt.Sprintf("%s[%d]", at, i)); err != nil {
				return err
			}
		}
	case "string":
		text, ok := v.(string)
		if !ok {
			return fmt.Errorf("%s: want string", at)
		}
		if sc.MinLength != nil && utf8.RuneCountInString(text) < *sc.MinLength {
			return fmt.Errorf("%s: shorter than %d", at, *sc.MinLength)
		}
	case "integer", "number":
		f, ok := v.(float64)
		if !ok || sc.Type == "integer" && f != math.Trunc(f) {
			return fmt.Errorf("%s: want %s", at, sc.Type)
		}
	case "boolean":
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("%s: want boolean", at)
		}
	case "null":
		if v != nil {
			return fmt.Errorf("%s: want null", at)
		}
	}
	if f, ok := v.(float64); ok && sc.Minimum != nil && f < *sc.Minimum {
		return fmt.Errorf("%s: below %v", at, *sc.Minimum)
	}
	if len(sc.Enum) > 0 && !slices.Contains(sc.Enum, v) {
		return fmt.Errorf("%s: not in enum", at)
	}
	return nil
}

//go:embed openapi.json
var relayDoc []byte

type Relay struct {
	ID     int      `json:"id"`
	Name   string   `json:"name"`
	Band   string   `json:"band"`
	Online bool     `json:"online"`
	Tags   []string `json:"tags"`
}

type relayStore struct {
	spec   *Spec
	mu     sync.Mutex
	next   int
	relays map[int]Relay
}

func relayJSON(w http.ResponseWriter, status int, v any) {
	raw, _ := json.Marshal(v)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(raw)
}

func (s *relayStore) create(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(io.LimitReader(r.Body, 4097))
	var v any
	var in Relay
	if err != nil || len(raw) > 4096 || json.Unmarshal(raw, &v) != nil ||
		s.spec.validate(&schema{Ref: refPrefix + "RelayInput"}, v, "input") != nil || json.Unmarshal(raw, &in) != nil {
		relayJSON(w, 400, map[string]string{"error": "invalid relay"})
		return
	}
	if in.Tags == nil {
		in.Tags = []string{}
	}
	s.mu.Lock()
	s.next++
	in.ID = s.next
	s.relays[in.ID] = in
	s.mu.Unlock()
	w.Header().Set("Location", "/relays/"+strconv.Itoa(in.ID))
	relayJSON(w, 201, in)
}

func (s *relayStore) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/relays" {
		switch r.Method {
		case "GET":
			s.mu.Lock()
			out := make([]Relay, 0, len(s.relays))
			for _, rl := range s.relays {
				out = append(out, rl)
			}
			s.mu.Unlock()
			slices.SortFunc(out, func(a, b Relay) int { return a.ID - b.ID })
			relayJSON(w, 200, out)
		case "POST":
			s.create(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		return
	}
	segment := strings.TrimPrefix(r.URL.Path, "/relays/")
	if !strings.HasPrefix(r.URL.Path, "/relays/") || strings.Contains(segment, "/") || r.Method != "GET" && r.Method != "DELETE" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(segment)
	s.mu.Lock()
	rl, found := s.relays[id]
	if found && err == nil && strconv.Itoa(id) == segment && r.Method == "DELETE" {
		delete(s.relays, id)
	}
	s.mu.Unlock()
	switch {
	case err != nil || strconv.Itoa(id) != segment || !found:
		relayJSON(w, 404, map[string]string{"error": "unknown relay"})
	case r.Method == "DELETE":
		w.WriteHeader(http.StatusNoContent)
	default:
		relayJSON(w, 200, rl)
	}
}

func NewRelayAPI() http.Handler {
	spec, err := LoadSpec(relayDoc)
	if err != nil {
		panic(err)
	}
	return &relayStore{spec: spec, relays: map[int]Relay{}}
}

func main() {}
