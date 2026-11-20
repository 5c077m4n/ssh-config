// Package clientconfig for mamanging the client's config
package clientconfig

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/keywords"
	"github.com/5c077m4n/ssh-config/internal/utils"
)

type (
	ClientOptions struct {
		AddKeysToAgent               string
		AddressFamily                string
		BatchMode                    *bool
		BindAddress                  string
		BindInterface                string
		CanonicalizeFallbackLocal    *bool
		CanonicalizeHostname         string
		CheckHostIP                  *bool
		Compression                  *bool
		ConnectTimeout               *uint64
		ControlMaster                string
		ControlPath                  string
		ForwardAgent                 *bool
		ForwardX11                   *bool
		ForwardX11Trusted            *bool
		HostKeyAlias                 string
		HostName                     string
		IdentitiesOnly               *bool
		IdentityAgent                string
		IdentityFile                 string
		KbdInteractiveAuthentication *bool
		LocalForward                 string
		LogLevel                     string
		PasswordAuthentication       *bool
		Port                         *uint64
		PreferredAuthentications     string
		ProxyCommand                 string
		ProxyJump                    string
		PubkeyAuthentication         *bool
		RemoteForward                string
		SendEnv                      string
		ServerAliveCountMax          *uint64
		ServerAliveInterval          *uint64
		SetEnv                       string
		StrictHostKeyChecking        string
		TCPKeepAlive                 *bool
		Tunnel                       string
		UpdateHostKeys               *bool
		User                         string
		VerifyHostKeyDNS             string
		VisualHostKey                *bool
	}
	Item  = utils.Item
	Block struct {
		Patterns []string
		Match    bool
		Items    []Item
		Options  *ClientOptions
	}
	File struct {
		Blocks []Block
	}
	ClientKey string
)

var clientKeysOrdered = []ClientKey{
	"AddKeysToAgent",
	"AddressFamily",
	"BatchMode",
	"BindAddress",
	"BindInterface",
	"CanonicalizeFallbackLocal",
	"CanonicalizeHostname",
	"CheckHostIP",
	"Compression",
	"ConnectTimeout",
	"ControlMaster",
	"ControlPath",
	"ForwardAgent",
	"ForwardX11",
	"ForwardX11Trusted",
	"HostKeyAlias",
	"HostName",
	"IdentitiesOnly",
	"IdentityAgent",
	"IdentityFile",
	"KbdInteractiveAuthentication",
	"LocalForward",
	"LogLevel",
	"PasswordAuthentication",
	"Port",
	"PreferredAuthentications",
	"ProxyCommand",
	"ProxyJump",
	"PubkeyAuthentication",
	"RemoteForward",
	"SendEnv",
	"ServerAliveCountMax",
	"ServerAliveInterval",
	"SetEnv",
	"StrictHostKeyChecking",
	"TCPKeepAlive",
	"Tunnel",
	"UpdateHostKeys",
	"User",
	"VerifyHostKeyDNS",
	"VisualHostKey",
}

func (o *ClientOptions) setField(canon, value string) error {
	switch canon {
	case "AddKeysToAgent":
		o.AddKeysToAgent = value
	case "AddressFamily":
		o.AddressFamily = value
	case "BatchMode":
		b := strings.EqualFold(value, "yes")
		o.BatchMode = &b
	case "BindAddress":
		o.BindAddress = value
	case "BindInterface":
		o.BindInterface = value
	case "CanonicalizeFallbackLocal":
		b := strings.EqualFold(value, "yes")
		o.CanonicalizeFallbackLocal = &b
	case "CanonicalizeHostname":
		o.CanonicalizeHostname = value
	case "CheckHostIP":
		b := strings.EqualFold(value, "yes")
		o.CheckHostIP = &b
	case "Compression":
		b := strings.EqualFold(value, "yes")
		o.Compression = &b
	case "ConnectTimeout":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.ConnectTimeout = &n
	case "ControlMaster":
		o.ControlMaster = value
	case "ControlPath":
		o.ControlPath = value
	case "ForwardAgent":
		b := strings.EqualFold(value, "yes")
		o.ForwardAgent = &b
	case "ForwardX11":
		b := strings.EqualFold(value, "yes")
		o.ForwardX11 = &b
	case "ForwardX11Trusted":
		b := strings.EqualFold(value, "yes")
		o.ForwardX11Trusted = &b
	case "HostKeyAlias":
		o.HostKeyAlias = value
	case "HostName":
		o.HostName = value
	case "IdentitiesOnly":
		b := strings.EqualFold(value, "yes")
		o.IdentitiesOnly = &b
	case "IdentityAgent":
		o.IdentityAgent = value
	case "IdentityFile":
		o.IdentityFile = value
	case "KbdInteractiveAuthentication":
		b := strings.EqualFold(value, "yes")
		o.KbdInteractiveAuthentication = &b
	case "LocalForward":
		o.LocalForward = value
	case "LogLevel":
		o.LogLevel = value
	case "PasswordAuthentication":
		b := strings.EqualFold(value, "yes")
		o.PasswordAuthentication = &b
	case "Port":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.Port = &n
	case "PreferredAuthentications":
		o.PreferredAuthentications = value
	case "ProxyCommand":
		o.ProxyCommand = value
	case "ProxyJump":
		o.ProxyJump = value
	case "PubkeyAuthentication":
		b := strings.EqualFold(value, "yes")
		o.PubkeyAuthentication = &b
	case "RemoteForward":
		o.RemoteForward = value
	case "SendEnv":
		o.SendEnv = value
	case "ServerAliveCountMax":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.ServerAliveCountMax = &n
	case "ServerAliveInterval":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.ServerAliveInterval = &n
	case "SetEnv":
		o.SetEnv = value
	case "StrictHostKeyChecking":
		o.StrictHostKeyChecking = value
	case "TCPKeepAlive":
		b := strings.EqualFold(value, "yes")
		o.TCPKeepAlive = &b
	case "Tunnel":
		o.Tunnel = value
	case "UpdateHostKeys":
		b := strings.EqualFold(value, "yes")
		o.UpdateHostKeys = &b
	case "User":
		o.User = value
	case "VerifyHostKeyDNS":
		o.VerifyHostKeyDNS = value
	case "VisualHostKey":
		b := strings.EqualFold(value, "yes")
		o.VisualHostKey = &b
	default:
		return fmt.Errorf("unknown key %q", canon)
	}
	return nil
}

func (o *ClientOptions) Set(key, value string) error {
	canon, ok := keywords.Canonical(keywords.Client, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	if err := keywords.Validate(keywords.Client, key, value); err != nil {
		return err
	}
	return o.setField(canon, value)
}

func (o *ClientOptions) field(canon string) (string, bool) {
	if o == nil {
		return "", false
	}
	switch canon {
	case "AddKeysToAgent":
		return utils.NonEmpty(o.AddKeysToAgent)
	case "AddressFamily":
		return utils.NonEmpty(o.AddressFamily)
	case "BatchMode":
		return utils.BoolField(o.BatchMode)
	case "BindAddress":
		return utils.NonEmpty(o.BindAddress)
	case "BindInterface":
		return utils.NonEmpty(o.BindInterface)
	case "CanonicalizeFallbackLocal":
		return utils.BoolField(o.CanonicalizeFallbackLocal)
	case "CanonicalizeHostname":
		return utils.NonEmpty(o.CanonicalizeHostname)
	case "CheckHostIP":
		return utils.BoolField(o.CheckHostIP)
	case "Compression":
		return utils.BoolField(o.Compression)
	case "ConnectTimeout":
		return utils.UintField(o.ConnectTimeout)
	case "ControlMaster":
		return utils.NonEmpty(o.ControlMaster)
	case "ControlPath":
		return utils.NonEmpty(o.ControlPath)
	case "ForwardAgent":
		return utils.BoolField(o.ForwardAgent)
	case "ForwardX11":
		return utils.BoolField(o.ForwardX11)
	case "ForwardX11Trusted":
		return utils.BoolField(o.ForwardX11Trusted)
	case "HostKeyAlias":
		return utils.NonEmpty(o.HostKeyAlias)
	case "HostName":
		return utils.NonEmpty(o.HostName)
	case "IdentitiesOnly":
		return utils.BoolField(o.IdentitiesOnly)
	case "IdentityAgent":
		return utils.NonEmpty(o.IdentityAgent)
	case "IdentityFile":
		return utils.NonEmpty(o.IdentityFile)
	case "KbdInteractiveAuthentication":
		return utils.BoolField(o.KbdInteractiveAuthentication)
	case "LocalForward":
		return utils.NonEmpty(o.LocalForward)
	case "LogLevel":
		return utils.NonEmpty(o.LogLevel)
	case "PasswordAuthentication":
		return utils.BoolField(o.PasswordAuthentication)
	case "Port":
		return utils.UintField(o.Port)
	case "PreferredAuthentications":
		return utils.NonEmpty(o.PreferredAuthentications)
	case "ProxyCommand":
		return utils.NonEmpty(o.ProxyCommand)
	case "ProxyJump":
		return utils.NonEmpty(o.ProxyJump)
	case "PubkeyAuthentication":
		return utils.BoolField(o.PubkeyAuthentication)
	case "RemoteForward":
		return utils.NonEmpty(o.RemoteForward)
	case "SendEnv":
		return utils.NonEmpty(o.SendEnv)
	case "ServerAliveCountMax":
		return utils.UintField(o.ServerAliveCountMax)
	case "ServerAliveInterval":
		return utils.UintField(o.ServerAliveInterval)
	case "SetEnv":
		return utils.NonEmpty(o.SetEnv)
	case "StrictHostKeyChecking":
		return utils.NonEmpty(o.StrictHostKeyChecking)
	case "TCPKeepAlive":
		return utils.BoolField(o.TCPKeepAlive)
	case "Tunnel":
		return utils.NonEmpty(o.Tunnel)
	case "UpdateHostKeys":
		return utils.BoolField(o.UpdateHostKeys)
	case "User":
		return utils.NonEmpty(o.User)
	case "VerifyHostKeyDNS":
		return utils.NonEmpty(o.VerifyHostKeyDNS)
	case "VisualHostKey":
		return utils.BoolField(o.VisualHostKey)
	default:
		return "", false
	}
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

func (b *Block) ensureOptions() {
	if b.Options == nil {
		b.Options = &ClientOptions{}
	}
}

func (f *File) Serialize() []byte {
	var b strings.Builder
	for _, blk := range f.Blocks {
		for _, it := range blk.Items {
			if it.Field != "" {
				if value, ok := blk.Options.field(it.Field); ok {
					b.WriteString("\t")
					b.WriteString(it.Field)
					b.WriteString(" ")
					b.WriteString(value)
					b.WriteString("\n")
				}
				continue
			}
			b.WriteString(it.Raw)
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
			if it.Field == canon {
				value, ok := blk.Options.field(canon)
				if ok {
					return value, nil
				}
			}
		}
	}
	return "", fmt.Errorf("%w: %s in %s", utils.ErrNotFound, canon, target)
}

func (f *File) setOption(target, canon, value string) error {
	for bi := range f.Blocks {
		blk := &f.Blocks[bi]
		if blk.Match || !containsPattern(blk.Patterns, target) {
			continue
		}
		blk.ensureOptions()
		if err := blk.Options.setField(canon, value); err != nil {
			return err
		}
		for ii := range blk.Items {
			if blk.Items[ii].Field == canon {
				return nil
			}
		}
		for ii := range blk.Items {
			if utils.CommentKey(blk.Items[ii].Raw) == canon {
				blk.Items[ii] = Item{Field: canon}
				return nil
			}
		}
		blk.Items = append(blk.Items, Item{Field: canon})
		return nil
	}
	blockOptions := &ClientOptions{}
	if err := blockOptions.setField(canon, value); err != nil {
		return err
	}
	f.Blocks = append(f.Blocks, Block{
		Patterns: strings.Fields(strings.ReplaceAll(target, ",", " ")),
		Options:  blockOptions,
		Items: []Item{
			{Raw: "Host " + target},
			{Field: canon},
		},
	})
	return nil
}

// Set applies every option configured on opts to the block matching
// target. Missing blocks are created.
func (f *File) Set(target string, opts *ClientOptions) error {
	for _, canon := range clientKeysOrdered {
		value, ok := opts.field(string(canon))
		if !ok {
			continue
		}
		if err := f.setOption(target, string(canon), value); err != nil {
			return err
		}
	}
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
		key, value := utils.SplitLine(trimmed)
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
			if canon, ok := keywords.Canonical(keywords.Client, key); ok {
				blk.ensureOptions()
				if _, set := blk.Options.field(canon); set {
					blk.Items = append(blk.Items, Item{Raw: raw})
				} else {
					_ = blk.Options.setField(canon, value)
					blk.Items = append(blk.Items, Item{Field: canon})
				}
			} else {
				blk.Items = append(blk.Items, Item{Raw: raw})
			}
		}
	}
	return f
}
