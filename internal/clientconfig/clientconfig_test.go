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

func opts(t *testing.T, key, value string) *ClientOptions {
	t.Helper()
	o := &ClientOptions{}
	require.NoError(t, o.Set(key, value))
	return o
}

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

func TestOptionsStructTyped(t *testing.T) {
	f := Parse([]byte(sample))

	block := f.Blocks[1]
	require.Equal(t, "host1.example.com", block.Options.HostName)
	require.NotNil(t, block.Options.Port)
	require.Equal(t, uint64(22), *block.Options.Port)
	require.Equal(t, "git", block.Options.User)
	require.Nil(t, block.Options.Compression)

	block = f.Blocks[2]
	require.NotNil(t, block.Options.ForwardAgent)
	require.True(t, *block.Options.ForwardAgent)
}

func TestSetReplacesExisting(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host1", opts(t, "port", "2222")))

	out := string(f.Serialize())
	require.Contains(t, out, "	Port 2222\n")
	require.Contains(t, out, "HostName host1.example.com")
	require.NotContains(t, out, "	Port 22\n")
}

func TestSetAppendsToBlock(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host1", opts(t, "compression", "yes")))

	out := string(f.Serialize())
	require.Contains(t, out, "Compression yes")
	require.Contains(t, out, "User git")
	require.NotContains(t, out, "Host host1.example.com")
}

func TestSetMultipleOptions(t *testing.T) {
	f := Parse([]byte(sample))

	o := &ClientOptions{}
	require.NoError(t, o.Set("port", "2222"))
	require.NoError(t, o.Set("user", "admin"))
	require.NoError(t, f.Set("host1", o))

	out := string(f.Serialize())
	require.Contains(t, out, "	Port 2222\n")
	require.Contains(t, out, "User admin")
	require.NotContains(t, out, "	Port 22\n")
}

func TestSetCreatesBlock(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("host2", opts(t, "user", "admin")))

	out := string(f.Serialize())
	require.Contains(t, out, "Host host2")
	require.Contains(t, out, "User admin")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := Parse([]byte(sample))

	require.Error(t, (&ClientOptions{}).Set("port", "abc"))
	require.Error(t, (&ClientOptions{}).Set("bogus", "yes"))
	require.Contains(t, string(f.Serialize()), "Port 22")
}

func TestMatchBlocksUntouched(t *testing.T) {
	f := Parse([]byte(sample))

	require.NoError(t, f.Set("bastion", opts(t, "user", "admin")))

	out := string(f.Serialize())
	require.Contains(t, out, "Match host bastion")
	require.Contains(t, out, "# inside match")
	require.Contains(t, out, "Host bastion")
}
