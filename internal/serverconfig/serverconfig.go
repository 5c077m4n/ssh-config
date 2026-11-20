// Package serverconfig for parsing server config
package serverconfig

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/5c077m4n/ssh-config/internal/keywords"
)

type (
	ServerOptions struct {
		AcceptEnv                    string
		AddressFamily                string
		AllowAgentForwarding         *bool
		AllowGroups                  string
		AllowTCPForwarding           string
		AllowUsers                   string
		Banner                       string
		ClientAliveCountMax          *uint64
		ClientAliveInterval          *uint64
		Compression                  string
		DenyGroups                   string
		DenyUsers                    string
		GSSAPIAuthentication         *bool
		HostKey                      string
		KbdInteractiveAuthentication *bool
		ListenAddress                string
		LogLevel                     string
		MaxAuthTries                 *uint64
		MaxSessions                  *uint64
		PasswordAuthentication       *bool
		PermitEmptyPasswords         *bool
		PermitRootLogin              string
		PermitTunnel                 string
		Port                         *uint64
		PubkeyAuthentication         *bool
		StrictModes                  *bool
		Subsystem                    string
		SyslogFacility               string
		TCPKeepAlive                 *bool
		UseDNS                       *bool
		X11DisplayOffset             *uint64
		X11Forwarding                *bool
	}
	Item struct {
		Raw   string
		Field string
	}
	File struct {
		Items   []Item
		Options *ServerOptions
	}
	ServerKey string
)

var serverKeysOrdered = []ServerKey{
	"AcceptEnv",
	"AddressFamily",
	"AllowAgentForwarding",
	"AllowGroups",
	"AllowTCPForwarding",
	"AllowUsers",
	"Banner",
	"ClientAliveCountMax",
	"ClientAliveInterval",
	"Compression",
	"DenyGroups",
	"DenyUsers",
	"GSSAPIAuthentication",
	"HostKey",
	"KbdInteractiveAuthentication",
	"ListenAddress",
	"LogLevel",
	"MaxAuthTries",
	"MaxSessions",
	"PasswordAuthentication",
	"PermitEmptyPasswords",
	"PermitRootLogin",
	"PermitTunnel",
	"Port",
	"PubkeyAuthentication",
	"StrictModes",
	"Subsystem",
	"SyslogFacility",
	"TCPKeepAlive",
	"UseDNS",
	"X11DisplayOffset",
	"X11Forwarding",
}

var ErrNotFound = errors.New("key not found")

func yesno(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func (o *ServerOptions) Set(key, value string) error {
	canon, ok := keywords.Canonical(keywords.SSHD, key)
	if !ok {
		return fmt.Errorf("unknown key %q", key)
	}
	if err := keywords.Validate(keywords.SSHD, key, value); err != nil {
		return err
	}
	return o.setField(canon, value)
}

func (o *ServerOptions) setField(canon, value string) error {
	switch canon {
	case "AcceptEnv":
		o.AcceptEnv = value
	case "AddressFamily":
		o.AddressFamily = value
	case "AllowAgentForwarding":
		b := strings.EqualFold(value, "yes")
		o.AllowAgentForwarding = &b
	case "AllowGroups":
		o.AllowGroups = value
	case "AllowTCPForwarding":
		o.AllowTCPForwarding = value
	case "AllowUsers":
		o.AllowUsers = value
	case "Banner":
		o.Banner = value
	case "ClientAliveCountMax":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.ClientAliveCountMax = &n
	case "ClientAliveInterval":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.ClientAliveInterval = &n
	case "Compression":
		o.Compression = value
	case "DenyGroups":
		o.DenyGroups = value
	case "DenyUsers":
		o.DenyUsers = value
	case "GSSAPIAuthentication":
		b := strings.EqualFold(value, "yes")
		o.GSSAPIAuthentication = &b
	case "HostKey":
		o.HostKey = value
	case "KbdInteractiveAuthentication":
		b := strings.EqualFold(value, "yes")
		o.KbdInteractiveAuthentication = &b
	case "ListenAddress":
		o.ListenAddress = value
	case "LogLevel":
		o.LogLevel = value
	case "MaxAuthTries":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.MaxAuthTries = &n
	case "MaxSessions":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.MaxSessions = &n
	case "PasswordAuthentication":
		b := strings.EqualFold(value, "yes")
		o.PasswordAuthentication = &b
	case "PermitEmptyPasswords":
		b := strings.EqualFold(value, "yes")
		o.PermitEmptyPasswords = &b
	case "PermitRootLogin":
		o.PermitRootLogin = value
	case "PermitTunnel":
		o.PermitTunnel = value
	case "Port":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.Port = &n
	case "PubkeyAuthentication":
		b := strings.EqualFold(value, "yes")
		o.PubkeyAuthentication = &b
	case "StrictModes":
		b := strings.EqualFold(value, "yes")
		o.StrictModes = &b
	case "Subsystem":
		o.Subsystem = value
	case "SyslogFacility":
		o.SyslogFacility = value
	case "TCPKeepAlive":
		b := strings.EqualFold(value, "yes")
		o.TCPKeepAlive = &b
	case "UseDNS":
		b := strings.EqualFold(value, "yes")
		o.UseDNS = &b
	case "X11DisplayOffset":
		n, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("value %q for %s must be a positive integer", value, canon)
		}
		o.X11DisplayOffset = &n
	case "X11Forwarding":
		b := strings.EqualFold(value, "yes")
		o.X11Forwarding = &b
	default:
		return fmt.Errorf("unknown key %q", canon)
	}
	return nil
}

func (o *ServerOptions) field(canon string) (string, bool) {
	if o == nil {
		return "", false
	}
	switch canon {
	case "AcceptEnv":
		return nonEmpty(o.AcceptEnv)
	case "AddressFamily":
		return nonEmpty(o.AddressFamily)
	case "AllowAgentForwarding":
		return boolField(o.AllowAgentForwarding)
	case "AllowGroups":
		return nonEmpty(o.AllowGroups)
	case "AllowTCPForwarding":
		return nonEmpty(o.AllowTCPForwarding)
	case "AllowUsers":
		return nonEmpty(o.AllowUsers)
	case "Banner":
		return nonEmpty(o.Banner)
	case "ClientAliveCountMax":
		return uintField(o.ClientAliveCountMax)
	case "ClientAliveInterval":
		return uintField(o.ClientAliveInterval)
	case "Compression":
		return nonEmpty(o.Compression)
	case "DenyGroups":
		return nonEmpty(o.DenyGroups)
	case "DenyUsers":
		return nonEmpty(o.DenyUsers)
	case "GSSAPIAuthentication":
		return boolField(o.GSSAPIAuthentication)
	case "HostKey":
		return nonEmpty(o.HostKey)
	case "KbdInteractiveAuthentication":
		return boolField(o.KbdInteractiveAuthentication)
	case "ListenAddress":
		return nonEmpty(o.ListenAddress)
	case "LogLevel":
		return nonEmpty(o.LogLevel)
	case "MaxAuthTries":
		return uintField(o.MaxAuthTries)
	case "MaxSessions":
		return uintField(o.MaxSessions)
	case "PasswordAuthentication":
		return boolField(o.PasswordAuthentication)
	case "PermitEmptyPasswords":
		return boolField(o.PermitEmptyPasswords)
	case "PermitRootLogin":
		return nonEmpty(o.PermitRootLogin)
	case "PermitTunnel":
		return nonEmpty(o.PermitTunnel)
	case "Port":
		return uintField(o.Port)
	case "PubkeyAuthentication":
		return boolField(o.PubkeyAuthentication)
	case "StrictModes":
		return boolField(o.StrictModes)
	case "Subsystem":
		return nonEmpty(o.Subsystem)
	case "SyslogFacility":
		return nonEmpty(o.SyslogFacility)
	case "TCPKeepAlive":
		return boolField(o.TCPKeepAlive)
	case "UseDNS":
		return boolField(o.UseDNS)
	case "X11DisplayOffset":
		return uintField(o.X11DisplayOffset)
	case "X11Forwarding":
		return boolField(o.X11Forwarding)
	default:
		return "", false
	}
}

func nonEmpty(s string) (string, bool) {
	return s, s != ""
}

func boolField(b *bool) (string, bool) {
	if b == nil {
		return "", false
	}
	return yesno(*b), true
}

func uintField(n *uint64) (string, bool) {
	if n == nil {
		return "", false
	}
	return strconv.FormatUint(*n, 10), true
}

func splitLine(line string) (string, string) {
	i := strings.IndexAny(line, " \t=")
	if i < 0 {
		return line, ""
	}
	value := strings.Trim(strings.TrimSpace(line[i+1:]), `"'`)
	return line[:i], value
}

func commentKey(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "#") {
		return ""
	}
	key, _ := splitLine(strings.TrimSpace(trimmed[1:]))
	return key
}

func (f *File) mainEnd() int {
	for i, it := range f.Items {
		if it.Field == "" &&
			strings.HasPrefix(strings.TrimSpace(strings.ToLower(it.Raw)), "match") {
			return i
		}
	}
	return len(f.Items)
}

func (f *File) Serialize() []byte {
	var b strings.Builder
	for _, it := range f.Items {
		if it.Field != "" {
			if value, ok := f.Options.field(it.Field); ok {
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
	return []byte(b.String())
}

func (f *File) Get(key string) (string, error) {
	canon, ok := keywords.Canonical(keywords.SSHD, key)
	if !ok {
		return "", fmt.Errorf("unknown key %q", key)
	}
	for _, it := range f.Items[:f.mainEnd()] {
		if it.Field == canon {
			value, ok := f.Options.field(canon)
			if ok {
				return value, nil
			}
		}
	}
	return "", fmt.Errorf("%w: %s", ErrNotFound, canon)
}

// Set applies every option configured on opts to the main (pre-Match)
// section of the config.
func (f *File) Set(opts *ServerOptions) error {
	for _, canon := range serverKeysOrdered {
		value, ok := opts.field(string(canon))
		if !ok {
			continue
		}
		if err := f.setOption(string(canon), value); err != nil {
			return err
		}
	}
	return nil
}

func (f *File) setOption(canon, value string) error {
	if f.Options == nil {
		f.Options = &ServerOptions{}
	}
	if err := f.Options.setField(canon, value); err != nil {
		return err
	}
	end := f.mainEnd()
	for i := range f.Items[:end] {
		if f.Items[i].Field == canon {
			return nil
		}
	}
	for i := range f.Items[:end] {
		if commentKey(f.Items[i].Raw) == canon {
			f.Items[i] = Item{Field: canon}
			return nil
		}
	}
	entry := Item{Field: canon}
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
		if canon, ok := keywords.Canonical(keywords.SSHD, key); ok {
			if f.Options == nil {
				f.Options = &ServerOptions{}
			}
			if _, set := f.Options.field(canon); set {
				f.Items = append(f.Items, Item{Raw: raw})
			} else {
				_ = f.Options.setField(canon, value)
				f.Items = append(f.Items, Item{Field: canon})
			}
		} else {
			f.Items = append(f.Items, Item{Raw: raw})
		}
	}
	return f
}
