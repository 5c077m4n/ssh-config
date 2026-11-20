package keywords

import (
	"fmt"
	"strconv"
	"strings"
)

type Kind int

const (
	String Kind = iota
	Bool
	Enum
	Uint
)

type Spec struct {
	Name string
	Kind Kind
	Enum []string
}

func lookup(table map[string]Spec, key string) (Spec, bool) {
	spec, ok := table[strings.ToLower(strings.TrimSpace(key))]
	return spec, ok
}

func Canonical(table map[string]Spec, key string) (string, bool) {
	if spec, ok := lookup(table, key); ok {
		return spec.Name, true
	}
	return "", false
}

func Validate(table map[string]Spec, key, value string) error {
	spec, ok := lookup(table, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	switch spec.Kind {
	case Bool:
		if !strings.EqualFold(value, "yes") && !strings.EqualFold(value, "no") {
			return fmt.Errorf("value %q for %s must be yes or no", value, spec.Name)
		}
	case Enum:
		for _, allowed := range spec.Enum {
			if strings.EqualFold(value, allowed) {
				return nil
			}
		}
		return fmt.Errorf(
			"value %q for %s must be one of %s",
			value,
			spec.Name,
			strings.Join(spec.Enum, ", "),
		)
	case Uint:
		n, err := strconv.Atoi(value)
		if err != nil || n < 1 {
			return fmt.Errorf("value %q for %s must be a positive integer", value, spec.Name)
		}
	default:
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("value for %s must not be empty", spec.Name)
		}
	}
	return nil
}
