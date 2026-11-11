package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"orbital/colony/labkit"
)

const passIssuer = "orbital-desk"
const passLife = 10 * time.Minute

// A valid hash nobody knows the code for: unknown crew still pays a bcrypt check.
var decoyHash, _ = bcrypt.GenerateFromPassword([]byte("decoy"), bcrypt.MinCost)

type PassDesk struct {
	secret []byte
	crew   map[string][]byte
	clock  labkit.Clock
}

type crewKey struct{}

func LoadSecret(getenv func(string) string) ([]byte, error) {
	raw := getenv("ORBITAL_DESK_SECRET")
	if raw == "" {
		return nil, errors.New("secret missing")
	}
	secret, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return nil, err
	}
	if len(secret) < 32 {
		return nil, errors.New("secret shorter than 32 bytes")
	}
	return secret, nil
}

func NewPassDesk(secret []byte, crew map[string][]byte, clock labkit.Clock) (*PassDesk, error) {
	if len(secret) < 32 || clock == nil {
		return nil, errors.New("invalid desk")
	}
	d := &PassDesk{secret: append([]byte(nil), secret...), crew: map[string][]byte{}, clock: clock}
	for name, hash := range crew {
		d.crew[name] = append([]byte(nil), hash...)
	}
	return d, nil
}

func plain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = io.WriteString(w, body)
}

func (d *PassDesk) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		plain(w, http.StatusMethodNotAllowed, "use POST\n")
		return
	}
	raw, err := io.ReadAll(io.LimitReader(r.Body, 4097))
	if err != nil || len(raw) > 4096 {
		plain(w, http.StatusBadRequest, "bad login\n")
		return
	}
	var in struct {
		Crew string `json:"crew"`
		Code string `json:"code"`
	}
	if json.Unmarshal(raw, &in) != nil || in.Crew == "" || in.Code == "" {
		plain(w, http.StatusBadRequest, "bad login\n")
		return
	}
	hash, known := d.crew[in.Crew]
	if !known {
		hash = decoyHash
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(in.Code))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) || err == nil && !known {
		plain(w, http.StatusUnauthorized, "denied\n")
		return
	}
	if err != nil {
		plain(w, http.StatusInternalServerError, "desk failure\n")
		return
	}
	now := d.clock.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    passIssuer,
		Subject:   in.Crew,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(passLife)),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(d.secret)
	if err != nil {
		plain(w, http.StatusInternalServerError, "desk failure\n")
		return
	}
	body, _ := json.Marshal(map[string]string{"pass": token})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (d *PassDesk) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		raw := strings.TrimPrefix(header, "Bearer ")
		var claims jwt.RegisteredClaims
		ok := strings.HasPrefix(header, "Bearer ") && raw != ""
		if ok {
			_, err := jwt.ParseWithClaims(raw, &claims, func(*jwt.Token) (any, error) { return d.secret, nil },
				jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer(passIssuer), jwt.WithTimeFunc(d.clock.Now),
				jwt.WithExpirationRequired(), jwt.WithIssuedAt())
			ok = err == nil && claims.Subject != ""
		}
		if !ok {
			w.Header().Set("WWW-Authenticate", "Bearer")
			plain(w, http.StatusUnauthorized, "pass required\n")
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), crewKey{}, claims.Subject)))
	})
}

func CrewName(r *http.Request) string {
	name, _ := r.Context().Value(crewKey{}).(string)
	return name
}

func main() {}
