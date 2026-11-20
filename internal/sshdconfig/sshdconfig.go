package sshdconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/keywords"
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
	File struct {
		Items []Item
	}
)

func (f *File) mainEnd() int {
	for i, it := range f.Items {
		if it.Entry == nil &&
			strings.HasPrefix(strings.TrimSpace(strings.ToLower(it.Raw)), "match") {
			return i
		}
	}
	return len(f.Items)
}

func (f *File) Serialize() []byte {
	var b strings.Builder
	for _, it := range f.Items {
		switch {
		case it.Entry == nil:
			b.WriteString(it.Raw)
		case it.Entry.raw != "":
			b.WriteString(it.Entry.raw)
		default:
			b.WriteString(it.Entry.Key)
			b.WriteString(" ")
			b.WriteString(it.Entry.Value)
		}
		b.WriteString("\n")
	}
	return []byte(b.String())
}

func (f *File) Get(key string) (string, error) {
	canon, ok := keywords.Canonical(keywords.SSHD, key)
	if !ok {
		return "", fmt.Errorf("unknown key %q", key)
	}
	for _, it := range f.Items[:f.mainEnd()] {
		if it.Entry != nil && strings.EqualFold(it.Entry.Key, canon) {
			return it.Entry.Value, nil
		}
	}
	return "", fmt.Errorf("%w: %s", ErrNotFound, canon)
}

func (f *File) Set(key, value string) error {
	canon, ok := keywords.Canonical(keywords.SSHD, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	if err := keywords.Validate(keywords.SSHD, key, value); err != nil {
		return err
	}
	end := f.mainEnd()
	for i := range f.Items[:end] {
		it := &f.Items[i]
		if it.Entry != nil && strings.EqualFold(it.Entry.Key, canon) {
			it.Entry.Value = value
			it.Entry.raw = ""
			return nil
		}
	}
	entry := Item{Entry: &Entry{Key: canon, Value: value}}
	f.Items = append(f.Items[:end], append([]Item{entry}, f.Items[end:]...)...)

	return nil
}

func Parse(data []byte) *File {
	f := &File{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for sc.Scan() {
		raw := sc.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			f.Items = append(f.Items, Item{Raw: raw})
			continue
		}
		key, value := splitLine(trimmed)
		if strings.EqualFold(key, "match") {
			f.Items = append(f.Items, Item{Raw: raw})
			continue
		}
		f.Items = append(f.Items, Item{Entry: &Entry{Key: key, Value: value, raw: raw}})
	}
	return f
}
