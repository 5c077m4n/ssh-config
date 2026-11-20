package clientconfig

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const sample = `# personal hosts

Host host1
	HostName host1.example.com
	User git
	Port 22

Host *.example.com
	User deploy
	ForwardAgent yes

Host *
	ServerAliveInterval 60

Match host bastion
	# inside match
`

func TestRoundTrip(t *testing.T) {
	require.Equal(t, sample, string(Parse([]byte(sample)).Serialize()))
}

func TestGet(t *testing.T) {
	f := Parse([]byte(sample))

	value, err := f.Get("host1", "user")
	require.NoError(t, err)
	require.Equal(t, "git", value)

	value, err = f.Get("HOST1", "USER")
	require.NoError(t, err)
	require.Equal(t, "git", value)

	value, err = f.Get("web.example.com", "user")
	require.NoError(t, err)
	require.Equal(t, "deploy", value)

	value, err = f.Get("unknown.example", "serveraliveinterval")
	require.NoError(t, err)
	require.Equal(t, "60", value)

	value, err = f.Get("host1", "SERVERALIVEINTERVAL")
	require.NoError(t, err)
	require.Equal(t, "60", value)

	_, err = f.Get("host1", "bogus")
	require.Error(t, err)

	_, err = f.Get("host1", "identityagent")
	require.Error(t, err)
}

func TestSetReplacesExisting(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host1", "port", "2222"))

	out := string(f.Serialize())
	require.Contains(t, out, "	Port 2222\n")
	require.Contains(t, out, "HostName host1.example.com")
	require.NotContains(t, out, "	Port 22\n")
}

func TestSetAppendsToBlock(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host1", "compression", "yes"))

	out := string(f.Serialize())
	require.Contains(t, out, "Compression yes")
	require.Contains(t, out, "User git")
	require.NotContains(t, out, "Host host1.example.com")
}

func TestSetCreatesBlock(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host2", "user", "admin"))

	out := string(f.Serialize())
	require.Contains(t, out, "Host host2")
	require.Contains(t, out, "User admin")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := Parse([]byte(sample))

	require.Error(t, f.Set("host1", "port", "abc"))
	require.Error(t, f.Set("host1", "bogus", "yes"))
	require.Contains(t, string(f.Serialize()), "Port 22")
}

func TestMatchBlocksUntouched(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("bastion", "user", "admin"))

	out := string(f.Serialize())
	require.Contains(t, out, "Match host bastion")
	require.Contains(t, out, "# inside match")
	require.Contains(t, out, "Host bastion")
}
