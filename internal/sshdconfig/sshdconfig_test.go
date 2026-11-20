package sshdconfig

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

func TestSetReplacesExisting(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("permitrootlogin", "no"))

	out := string(f.Serialize())
	require.Contains(t, out, "PermitRootLogin no")
	require.NotContains(t, out, "PermitRootLogin prohibit-password")
	require.Contains(t, out, "Match group admins")
	require.Contains(t, out, "\tPermitRootLogin yes")
}

func TestSetInsertsBeforeMatch(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("x11forwarding", "yes"))

	out := string(f.Serialize())
	require.Contains(t, out, "X11Forwarding yes")
	require.NotContains(t, out, "\tX11Forwarding yes")
	require.Contains(t, out, "Match group admins")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := Parse([]byte(sample))

	require.Error(t, f.Set("permitrootlogin", "maybe"))
	require.Error(t, f.Set("bogus", "yes"))
	require.Contains(t, string(f.Serialize()), "PermitRootLogin prohibit-password")
}
