package levels

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"testing"
)

func sandboxConfigParity(t *testing.T, arabic, english map[string]any) {
	t.Helper()
	// Compare the JSON shape the client receives, so typed structs and maps
	// are treated alike and only string leaves may differ.
	if difference := sandboxConfigDifference(reflect.ValueOf(sandboxConfigJSON(arabic)), reflect.ValueOf(sandboxConfigJSON(english)), "config"); difference != "" {
		t.Errorf("sandbox %s differs across languages", difference)
	}
}

func sandboxConfigDifference(arabic, english reflect.Value, path string) string {
	arabic, english = unwrapSandboxConfigValue(arabic), unwrapSandboxConfigValue(english)
	if !arabic.IsValid() || !english.IsValid() {
		return sandboxConfigPresenceDifference(arabic, english, path)
	}
	if arabic.Type() != english.Type() {
		return path
	}
	switch arabic.Kind() {
	case reflect.String:
		return ""
	case reflect.Map:
		return sandboxConfigMapDifference(arabic, english, path)
	case reflect.Slice, reflect.Array:
		return sandboxConfigSliceDifference(arabic, english, path)
	default:
		return sandboxConfigLeafDifference(arabic, english, path)
	}
}

func sandboxConfigPresenceDifference(arabic, english reflect.Value, path string) string {
	if arabic.IsValid() == english.IsValid() {
		return ""
	}
	return path
}

func sandboxConfigLeafDifference(arabic, english reflect.Value, path string) string {
	if reflect.DeepEqual(arabic.Interface(), english.Interface()) {
		return ""
	}
	return path
}

func unwrapSandboxConfigValue(configValue reflect.Value) reflect.Value {
	for configValue.IsValid() && configValue.Kind() == reflect.Interface {
		if configValue.IsNil() {
			return reflect.Value{}
		}
		configValue = configValue.Elem()
	}
	return configValue
}

func sandboxConfigMapDifference(arabic, english reflect.Value, path string) string {
	if arabic.Len() != english.Len() {
		return path
	}
	iter := arabic.MapRange()
	for iter.Next() {
		key := iter.Key()
		keyPath := fmt.Sprintf("%s.%v", path, key.Interface())
		translated := english.MapIndex(key)
		if !translated.IsValid() {
			return keyPath
		}
		if difference := sandboxConfigDifference(iter.Value(), translated, keyPath); difference != "" {
			return difference
		}
	}
	return ""
}

func sandboxConfigSliceDifference(arabic, english reflect.Value, path string) string {
	if arabic.Len() != english.Len() {
		return path
	}
	for index := 0; index < arabic.Len(); index++ {
		if difference := sandboxConfigDifference(arabic.Index(index), english.Index(index), fmt.Sprintf("%s[%d]", path, index)); difference != "" {
			return difference
		}
	}
	return ""
}

func TestSandboxConfigParityAllowsLocalizedStrings(t *testing.T) {
	arabic := map[string]any{"labels": map[string]any{"ar": "التالي", "en": "Next"}, "routes": []any{"rim", "relay", "dock"}, "limit": 3}
	english := map[string]any{"labels": map[string]any{"ar": "التالي", "en": "Advance"}, "routes": []any{"rim", "relay", "dock"}, "limit": 3}
	if difference := sandboxConfigDifference(reflect.ValueOf(arabic), reflect.ValueOf(english), "config"); difference != "" {
		t.Fatalf("localized strings should be allowed to differ at %s", difference)
	}
	english["limit"] = 4
	if difference := sandboxConfigDifference(reflect.ValueOf(arabic), reflect.ValueOf(english), "config"); difference == "" {
		t.Fatal("non-string data differences must be reported")
	}
	delete(english, "limit")
	if difference := sandboxConfigDifference(reflect.ValueOf(arabic), reflect.ValueOf(english), "config"); difference == "" {
		t.Fatal("key differences must be reported")
	}
}

func sandboxRegistrationFiles(t *testing.T) []string {
	t.Helper()
	lessonFiles, err := filepath.Glob("../../web/lessons*.js")
	if err != nil {
		t.Fatal(err)
	}
	chapterFiles, err := filepath.Glob("../../web/sandbox/*.js")
	if err != nil {
		t.Fatal(err)
	}
	return append(lessonFiles, chapterFiles...)
}

func registeredSandboxKinds(t *testing.T) map[string]bool {
	t.Helper()
	registration := regexp.MustCompile(`Lessons\.registerSandbox\s*\(\s*["']([^"']+)["']`)
	registered := map[string]bool{}
	for _, path := range sandboxRegistrationFiles(t) {
		source, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range registration.FindAllSubmatch(source, -1) {
			registered[string(match[1])] = true
		}
	}
	return registered
}

func TestSandboxKindsHaveRenderers(t *testing.T) {
	registered := registeredSandboxKinds(t)
	used := map[string]bool{}
	for _, level := range All {
		for _, step := range level.Steps {
			for _, block := range step.Blocks {
				if block.Sandbox != nil {
					used[block.Sandbox.Kind] = true
				}
			}
		}
	}
	for kind := range used {
		if !registered[kind] {
			t.Errorf("sandbox kind %q has no Lessons.registerSandbox call", kind)
		}
	}
}

func sandboxConfigJSON(config map[string]any) any {
	raw, err := json.Marshal(config)
	if err != nil {
		return config
	}
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return config
	}
	return v
}
