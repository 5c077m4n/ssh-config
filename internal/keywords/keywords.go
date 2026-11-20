// Package keywords for the config options
package keywords

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type (
	Spec interface {
		KeywordName() string
		Validate(value string) error
	}
	SpecMap = map[string]Spec
)

func lookup(table SpecMap, key string) (Spec, bool) {
	spec, ok := table[strings.ToLower(strings.TrimSpace(key))]
	return spec, ok
}

func Canonical(table SpecMap, key string) (string, bool) {
	if spec, ok := lookup(table, key); ok {
		return spec.KeywordName(), true
	}
	return "", false
}

func Names(table SpecMap) []string {
	names := make([]string, 0, len(table))
	for _, spec := range table {
		names = append(names, spec.KeywordName())
	}
	sort.Strings(names)
	return names
}

func Validate(table SpecMap, key, value string) error {
	spec, ok := lookup(table, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	return spec.Validate(value)
}

type typed[T any] struct {
	name  string
	parse func(string) (T, error)
}

func (s typed[T]) KeywordName() string {
	return s.name
}

func (s typed[T]) Validate(value string) error {
	_, err := s.parse(value)
	return err
}

func Bool(name string) Spec {
	return typed[bool]{
		name: name,
		parse: func(value string) (bool, error) {
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
		parse: func(value string) (uint64, error) {
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
		parse: func(value string) (string, error) {
			if strings.TrimSpace(value) == "" {
				return "", fmt.Errorf("value for %s must not be empty", name)
			}
			return value, nil
		},
	}
}

func Enum(name string, allowed []string) Spec {
	return typed[string]{
		name: name,
		parse: func(value string) (string, error) {
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
