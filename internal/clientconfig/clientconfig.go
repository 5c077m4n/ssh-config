// Package clientconfig for mamanging the client's config
package clientconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/keywords"
)

type (
	Entry struct {
		Key   string
		Value string
		raw   string
	}
	Item struct {
		Entry *Entry
		Raw   string
	}
	Block struct {
		Patterns []string
		Match    bool
		Items    []Item
	}
	File struct {
		Blocks []Block
	}
)

var ErrNotFound = errors.New("key not found")

func splitLine(line string) (string, string) {
	i := strings.IndexAny(line, " \t=")
	if i < 0 {
		return line, ""
	}
	value := strings.Trim(strings.TrimSpace(line[i+1:]), `"'`)
	return line[:i], value
}

func containsPattern(patterns []string, target string) bool {
	for _, p := range patterns {
		if strings.EqualFold(p, target) {
			return true
		}
	}
	return false
}

func globMatch(p, s string) bool {
	if p == s {
		return true
	}
	if !strings.ContainsAny(p, "*?") {
		return false
	}
	if p == "" {
		return s == ""
	}
	switch p[0] {
	case '*':
		for i := 0; i <= len(s); i++ {
			if globMatch(p[1:], s[i:]) {
				return true
			}
		}
		return false
	case '?':
		return s != "" && globMatch(p[1:], s[1:])
	default:
		return s != "" && p[0] == s[0] && globMatch(p[1:], s[1:])
	}
}

func matchPattern(pattern, name string) bool {
	return globMatch(strings.ToLower(pattern), strings.ToLower(name))
}

func (b *Block) matches(target string) bool {
	if b.Patterns == nil {
		return true
	}
	for _, p := range b.Patterns {
		if strings.HasPrefix(p, "!") {
			continue
		}
		if matchPattern(p, target) {
			return true
		}
	}
	return false
}

func (f *File) Serialize() []byte {
	var b strings.Builder
	for _, blk := range f.Blocks {
		for _, it := range blk.Items {
			switch {
			case it.Entry == nil:
				b.WriteString(it.Raw)
			case it.Entry.raw != "":
				b.WriteString(it.Entry.raw)
			default:
				b.WriteString("\t")
				b.WriteString(it.Entry.Key)
				b.WriteString(" ")
				b.WriteString(it.Entry.Value)
			}
			b.WriteString("\n")
		}
	}
	return []byte(b.String())
}

func (f *File) Get(target, key string) (string, error) {
	canon, ok := keywords.Canonical(keywords.Client, key)
	if !ok {
		return "", fmt.Errorf("unknown key %q", key)
	}
	for _, blk := range f.Blocks {
		if blk.Match || !blk.matches(target) {
			continue
		}
		for _, it := range blk.Items {
			if it.Entry != nil && strings.EqualFold(it.Entry.Key, canon) {
				return it.Entry.Value, nil
			}
		}
	}
	return "", fmt.Errorf("%w: %s in %s", ErrNotFound, canon, target)
}

func (f *File) Set(target, key, value string) error {
	canon, ok := keywords.Canonical(keywords.Client, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	if err := keywords.Validate(keywords.Client, key, value); err != nil {
		return err
	}
	for bi := range f.Blocks {
		blk := &f.Blocks[bi]
		if blk.Match || !containsPattern(blk.Patterns, target) {
			continue
		}
		for ii := range blk.Items {
			it := &blk.Items[ii]
			if it.Entry != nil && strings.EqualFold(it.Entry.Key, canon) {
				it.Entry.Value = value
				it.Entry.raw = ""
				return nil
			}
		}
		blk.Items = append(blk.Items, Item{Entry: &Entry{Key: canon, Value: value}})
		return nil
	}
	f.Blocks = append(f.Blocks, Block{
		Patterns: strings.Fields(strings.ReplaceAll(target, ",", " ")),
		Items: []Item{
			{Raw: "Host " + target},
			{Entry: &Entry{Key: canon, Value: value}},
		},
	})
	return nil
}

func Parse(data []byte) *File {
	f := &File{Blocks: []Block{{}}}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		raw := sc.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			blk := &f.Blocks[len(f.Blocks)-1]
			blk.Items = append(blk.Items, Item{Raw: raw})
			continue
		}
		key, value := splitLine(trimmed)
		switch strings.ToLower(key) {
		case "host":
			f.Blocks = append(
				f.Blocks,
				Block{Patterns: strings.Fields(value), Items: []Item{{Raw: raw}}},
			)
		case "match":
			f.Blocks = append(f.Blocks, Block{Match: true, Items: []Item{{Raw: raw}}})
		default:
			blk := &f.Blocks[len(f.Blocks)-1]
			blk.Items = append(blk.Items, Item{Entry: &Entry{Key: key, Value: value, raw: raw}})
		}
	}
	return f
}
