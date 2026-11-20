package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/clientconfig"
	"github.com/5c077m4n/ssh-config/internal/keywords"
)

func split(addr string) (string, string, error) {
	i := strings.LastIndex(addr, ".")
	if i <= 0 {
		return "", "", fmt.Errorf("expected '[host.]key', got %q", addr)
	}
	return addr[:i], addr[i+1:], nil
}

func defaultConfigPath() string {
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

func usage(flagSet *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "usage: sshc [flags] '[host.]key' [value]")
	flagSet.PrintDefaults()
}

func run(args []string) int {
	flagSet := flag.NewFlagSet("sshc", flag.ContinueOnError)
	configPath := flagSet.String(
		"config-path",
		defaultConfigPath(),
		"path to the SSH client config",
	)
	flagSet.Usage = func() { usage(flagSet) }
	if err := flagSet.Parse(args); err != nil {
		return 2
	}
	args = flagSet.Args()

	if len(args) == 0 || len(args) > 2 {
		flagSet.Usage()
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
		return set(*configPath, target, key, args[1])
	}
	return get(*configPath, target, key)
}

func main() {
	os.Exit(run(os.Args[1:]))
}
