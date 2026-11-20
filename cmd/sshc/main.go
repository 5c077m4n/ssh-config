package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/clientconfig"
	"github.com/5c077m4n/ssh-config/internal/keywords"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: sshc '[host.]key' [value]")
}

func run(args []string) int {
	if len(args) == 0 || len(args) > 2 {
		usage()
		return 2
	}
	target, key, err := split(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 2
	}
	if _, ok := keywords.Canonical(keywords.Client, key); !ok {
		fmt.Fprintf(os.Stderr, "sshc: unknown key %q\n", key)
		return 1
	}
	if len(args) == 2 {
		if err := keywords.Validate(keywords.Client, key, args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "sshc:", err)
			return 1
		}
		return set(configPath(), target, key, args[1])
	}
	return get(configPath(), target, key)
}

func split(addr string) (string, string, error) {
	i := strings.LastIndex(addr, ".")
	if i <= 0 {
		return "", "", fmt.Errorf("expected '[host.]key', got %q", addr)
	}
	return addr[:i], addr[i+1:], nil
}

func configPath() string {
	if p := os.Getenv("SSHC_CONFIG"); p != "" {
		return p
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".ssh/config"
	}
	return filepath.Join(home, ".ssh", "config")
}

func get(path, target, key string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 1
	}
	value, err := clientconfig.Parse(data).Get(target, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 1
	}
	fmt.Println(value)
	return 0
}

func set(path, target, key, value string) int {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 1
	}
	f := clientconfig.Parse(data)
	if err := f.Set(target, key, value); err != nil {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 1
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			fmt.Fprintln(os.Stderr, "sshc:", err)
			return 1
		}
	}
	if err := os.WriteFile(path, f.Serialize(), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "sshc:", err)
		return 1
	}
	return 0
}
