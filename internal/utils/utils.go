// Package utils for helpers shared by the config parsers
package utils

import (
	"errors"
	"strconv"
	"strings"
)

type Item struct {
	Raw   string
	Field string
}

var ErrNotFound = errors.New("key not found")

func YesOrNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func NonEmpty(s string) (string, bool) {
	return s, s != ""
}

func BoolField(b *bool) (string, bool) {
	if b == nil {
		return "", false
	}
	return YesOrNo(*b), true
}

func UintField(n *uint64) (string, bool) {
	if n == nil {
		return "", false
	}
	return strconv.FormatUint(*n, 10), true
}

func SplitLine(line string) (string, string) {
	i := strings.IndexAny(line, " \t=")
	if i < 0 {
		return line, ""
	}
	value := strings.Trim(strings.TrimSpace(line[i+1:]), `"'`)
	return line[:i], value
}

func CommentKey(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "#") {
		return ""
	}
	key, _ := SplitLine(strings.TrimSpace(trimmed[1:]))
	return key
}
