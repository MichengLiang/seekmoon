package source

import (
	"encoding/json"
	"fmt"

	"github.com/yumiaura/seekmoon/internal/model"
)

func jsonUnmarshal(data []byte, target any) error {
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	return nil
}

func evidenceString(value, source string) model.EvidenceString {
	if value == "" {
		return model.Missing[string](source)
	}
	return model.Present(value, source)
}

func evidenceOptionalString(value *string, source string) model.EvidenceString {
	if value == nil || *value == "" {
		return model.Missing[string](source)
	}
	return model.Present(*value, source)
}

func evidenceStrings(values []string, source string) model.EvidenceStringArray {
	if len(values) == 0 {
		return model.Missing[[]string](source)
	}
	return model.Present(values, source)
}

func stringFromRaw(raw map[string]any, key string) string {
	value, _ := raw[key].(string)
	return value
}

func stringsFromRaw(raw map[string]any, key string) []string {
	value, ok := raw[key]
	if !ok {
		return nil
	}
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		out := make([]string, 0, len(typed))
		for _, item := range typed {
			if s, ok := item.(string); ok {
				out = append(out, s)
			}
		}
		return out
	case string:
		if typed == "" {
			return nil
		}
		return []string{typed}
	default:
		return nil
	}
}

func firstStringFromRaw(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		if value := stringFromRaw(raw, key); value != "" {
			return value
		}
	}
	return ""
}

func firstStringsFromRaw(raw map[string]any, keys ...string) []string {
	for _, key := range keys {
		if value := stringsFromRaw(raw, key); len(value) > 0 {
			return value
		}
	}
	return nil
}

func stringMapFromRaw(raw map[string]any, key string) map[string]string {
	value, ok := raw[key]
	if !ok {
		return nil
	}
	out := map[string]string{}
	switch typed := value.(type) {
	case map[string]string:
		return typed
	case map[string]any:
		for k, v := range typed {
			out[k] = fmt.Sprint(v)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
