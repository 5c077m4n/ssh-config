package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/5c077m4n/ssh-config/internal/keywords"
	"github.com/5c077m4n/ssh-config/internal/serverconfig"
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
	f := &serverconfig.File{}
	if err := f.UnmarshalSSH(data); err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	value, err := f.Get(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	fmt.Println(value)
	return 0
}

func set(path string, opts *serverconfig.ServerOptions) int {
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	f := &serverconfig.File{}
	if err := f.UnmarshalSSH(data); err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	if err := f.Set(opts); err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	b, err := f.MarshalSSH()
	if err != nil {
		fmt.Fprintln(os.Stderr, "sshdc:", err)
		return 1
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
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
	if _, ok := keywords.SSHD.Canonical(key); !ok {
		fmt.Fprintf(os.Stderr, "sshdc: unknown key %q\n", key)
		return 1
	}
	if len(args) == 2 {
		opts := &serverconfig.ServerOptions{}
		if err := opts.Set(key, args[1], filepath.Dir(*configPath)); err != nil {
			fmt.Fprintln(os.Stderr, "sshdc:", err)
			return 1
		}
		return set(*configPath, opts)
	}
	return get(*configPath, key)
}

func main() {
	os.Exit(run(os.Args[1:]))
}
