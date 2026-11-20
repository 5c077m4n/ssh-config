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

func mustParse(t *testing.T, data string) *File {
	t.Helper()

	f := &File{}
	require.NoError(t, f.UnmarshalSSH([]byte(data)))
	return f
}

func marshal(t *testing.T, f *File) string {
	t.Helper()

	b, err := f.MarshalSSH()
	require.NoError(t, err)
	return string(b)
}

func opts(t *testing.T, key, value string) *ServerOptions {
	t.Helper()
	o := &ServerOptions{}
	require.NoError(t, o.Set(key, value, ""))
	return o
}

func TestRoundTrip(t *testing.T) {
	require.Equal(t, sample, marshal(t, mustParse(t, sample)))
}

func TestGet(t *testing.T) {
	f := mustParse(t, sample)

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
	f := mustParse(t, sample)

	require.NotNil(t, f.Options.Port)
	require.Equal(t, uint64(22), *f.Options.Port)
	require.Equal(t, "prohibit-password", f.Options.PermitRootLogin)
	require.Nil(t, f.Options.X11Forwarding)
}

func TestSetReplacesExisting(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set(opts(t, "permitrootlogin", "no")))

	out := marshal(t, f)
	require.Contains(t, out, "PermitRootLogin no")
	require.NotContains(t, out, "PermitRootLogin prohibit-password")
	require.Contains(t, out, "Match group admins")
	require.Contains(t, out, "	PermitRootLogin yes")
}

func TestSetInsertsBeforeMatch(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set(opts(t, "x11forwarding", "yes")))

	out := marshal(t, f)
	require.Contains(t, out, "X11Forwarding yes")
	require.NotContains(t, out, "	X11Forwarding yes")
	require.Contains(t, out, "Match group admins")
}

func TestSetMultipleOptions(t *testing.T) {
	f := mustParse(t, sample)

	o := &ServerOptions{}
	require.NoError(t, o.Set("x11forwarding", "yes", ""))
	require.NoError(t, o.Set("maxauthtries", "4", ""))
	require.NoError(t, f.Set(o))

	out := marshal(t, f)
	require.Contains(t, out, "X11Forwarding yes")
	require.Contains(t, out, "MaxAuthTries 4")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := mustParse(t, sample)

	require.Error(t, (&ServerOptions{}).Set("permitrootlogin", "maybe", ""))
	require.Error(t, (&ServerOptions{}).Set("bogus", "yes", ""))
	require.Contains(t, marshal(t, f), "PermitRootLogin prohibit-password")
}
