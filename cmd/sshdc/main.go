package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/5c077m4n/ssh-config/internal/keywords"
	"github.com/5c077m4n/ssh-config/internal/sshdconfig"
)

func usage(fs *flag.FlagSet) {
	fmt.Fprintln(os.Stderr, "usage: sshdc [flags] 'key' [value]")
	fs.PrintDefaults()
}

func defaultConfigPath() string {
	return "/etc/ssh/sshd_config"
}

func get(path, key string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	value, err := sshdconfig.Parse(data).Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	fmt.Println(value)
	return 0
}

func set(path, key, value string) int {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	f := sshdconfig.Parse(data)
	if err := f.Set(key, value); err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	if err := os.WriteFile(path, f.Serialize(), 0o600); err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	return 0
}

func run(args []string) int {
	fs := flag.NewFlagSet("sshdc", flag.ContinueOnError)
	configPath := fs.String("config-path", defaultConfigPath(), "path to the SSH server config")
	fs.Usage = func() { usage(fs) }
	if err := fs.Parse(args); err != nil {
		return 2
	}
	args = fs.Args()

	if len(args) == 0 || len(args) > 2 {
		fs.Usage()
		return 2
	}
	key := args[0]
	if _, ok := keywords.Canonical(keywords.SSHD, key); !ok {
		fmt.Fprintf(os.Stderr, "sshdc: unknown key %q\n", key)
		return 1
	}
	if len(args) == 2 {
		if err := keywords.Validate(keywords.SSHD, key, args[1]); err != nil {
			fmt.Fprintln(os.Stderr, "sshdc:", err)
			return 1
		}
		return set(*configPath, key, args[1])
	}
	return get(*configPath, key)
}

func main() {
	os.Exit(run(os.Args[1:]))
}
