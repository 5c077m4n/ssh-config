// Package keywords for the config options
package keywords

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type (
	Spec interface {
		KeywordName() string
		Validate(value string) error
		ValidateIn(base, value string) error
	}
	SpecMap map[string]Spec
)

func (specMap SpecMap) lookup(key string) (Spec, bool) {
	spec, ok := specMap[strings.ToLower(strings.TrimSpace(key))]
	return spec, ok
}

func (specMap SpecMap) Canonical(key string) (string, bool) {
	if spec, ok := specMap.lookup(key); ok {
		return spec.KeywordName(), true
	}
	return "", false
}

func (specMap SpecMap) Names() []string {
	names := make([]string, 0, len(specMap))
	for _, spec := range specMap {
		names = append(names, spec.KeywordName())
	}
	sort.Strings(names)
	return names
}

func (specMap SpecMap) Validate(key, value, base string) error {
	spec, ok := specMap.lookup(key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	return spec.ValidateIn(base, value)
}

type typed[T any] struct {
	name  string
	parse func(base, value string) (T, error)
}

func (t typed[T]) KeywordName() string {
	return t.name
}

func (t typed[T]) Validate(value string) error {
	_, err := t.parse("", value)
	return err
}

func (t typed[T]) ValidateIn(base, value string) error {
	_, err := t.parse(base, value)
	return err
}

func Bool(name string) Spec {
	return typed[bool]{
		name: name,
		parse: func(_, value string) (bool, error) {
			if strings.EqualFold(value, "yes") {
				return true, nil
			}
			if strings.EqualFold(value, "no") {
				return false, nil
			}
			return false, fmt.Errorf("value %q for %s must be yes or no", value, name)
		},
	}
}

func Uint(name string) Spec {
	return typed[uint64]{
		name: name,
		parse: func(_, value string) (uint64, error) {
			n, err := strconv.ParseUint(value, 10, 64)
			if err != nil || n < 1 {
				return 0, fmt.Errorf("value %q for %s must be a positive integer", value, name)
			}
			return n, nil
		},
	}
}

func Str(name string) Spec {
	return typed[string]{
		name: name,
		parse: func(_, value string) (string, error) {
			if strings.TrimSpace(value) == "" {
				return "", fmt.Errorf("value for %s must not be empty", name)
			}
			return value, nil
		},
	}
}

func Path(name string) Spec {
	return typed[string]{
		name: name,
		parse: func(base, value string) (string, error) {
			if strings.TrimSpace(value) == "" {
				return "", fmt.Errorf("value for %s must not be empty", name)
			}
			for p := range strings.FieldsSeq(value) {
				// sshd expands % tokens at runtime; cannot validate here
				if strings.EqualFold(p, "none") || strings.ContainsRune(p, '%') {
					continue
				}
				if !filepath.IsAbs(p) {
					if base == "" {
						continue
					}
					p = filepath.Join(base, p)
				}
				if _, err := os.Stat(p); err != nil {
					return "", fmt.Errorf(
						"value %q for %s: path does not exist",
						value,
						name,
					)
				}
			}
			return value, nil
		},
	}
}

func Enum(name string, allowed []string) Spec {
	return typed[string]{
		name: name,
		parse: func(_, value string) (string, error) {
			for _, a := range allowed {
				if strings.EqualFold(value, a) {
					return value, nil
				}
			}
			return "", fmt.Errorf(
				"value %q for %s must be one of %s",
				value,
				name,
				strings.Join(allowed, ", "),
			)
		},
	}
}
