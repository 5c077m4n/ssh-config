package serverconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const sample = `# sshd
Port 22
PermitRootLogin prohibit-password

Match group admins
	PermitRootLogin yes
`

func opts(t *testing.T, key, value string) *ServerOptions {
	t.Helper()
	o := &ServerOptions{}
	require.NoError(t, o.Set(key, value))
	return o
}

func TestRoundTrip(t *testing.T) {
	require.Equal(t, sample, string(Parse([]byte(sample)).Serialize()))
}

func TestGet(t *testing.T) {
	f := Parse([]byte(sample))

	value, err := f.Get("permitrootlogin")
	require.NoError(t, err)
	require.Equal(t, "prohibit-password", value)

	value, err = f.Get("PERMITROOTLOGIN")
	require.NoError(t, err)
	require.Equal(t, "prohibit-password", value)

	_, err = f.Get("x11forwarding")
	require.Error(t, err)

	_, err = f.Get("bogus")
	require.Error(t, err)
}

func TestOptionsStructTyped(t *testing.T) {
	f := Parse([]byte(sample))

	require.NotNil(t, f.Options.Port)
	require.Equal(t, uint64(22), *f.Options.Port)
	require.Equal(t, "prohibit-password", f.Options.PermitRootLogin)
	require.Nil(t, f.Options.X11Forwarding)
}

func TestSetReplacesExisting(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set(opts(t, "permitrootlogin", "no")))

	out := string(f.Serialize())
	require.Contains(t, out, "PermitRootLogin no")
	require.NotContains(t, out, "PermitRootLogin prohibit-password")
	require.Contains(t, out, "Match group admins")
	require.Contains(t, out, "	PermitRootLogin yes")
}

func TestSetInsertsBeforeMatch(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set(opts(t, "x11forwarding", "yes")))

	out := string(f.Serialize())
	require.Contains(t, out, "X11Forwarding yes")
	require.NotContains(t, out, "	X11Forwarding yes")
	require.Contains(t, out, "Match group admins")
}

func TestSetMultipleOptions(t *testing.T) {
	f := Parse([]byte(sample))

	o := &ServerOptions{}
	require.NoError(t, o.Set("x11forwarding", "yes"))
	require.NoError(t, o.Set("maxauthtries", "4"))
	require.NoError(t, f.Set(o))

	out := string(f.Serialize())
	require.Contains(t, out, "X11Forwarding yes")
	require.Contains(t, out, "MaxAuthTries 4")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := Parse([]byte(sample))

	require.Error(t, (&ServerOptions{}).Set("permitrootlogin", "maybe"))
	require.Error(t, (&ServerOptions{}).Set("bogus", "yes"))
	require.Contains(t, string(f.Serialize()), "PermitRootLogin prohibit-password")
}
