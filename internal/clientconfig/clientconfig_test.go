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

func opts(t *testing.T, key, value string) *ClientOptions {
	t.Helper()
	o := &ClientOptions{}
	require.NoError(t, o.Set(key, value))
	return o
}

func TestRoundTrip(t *testing.T) {
	require.Equal(t, sample, marshal(t, mustParse(t, sample)))
}

func TestGet(t *testing.T) {
	f := mustParse(t, sample)

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
	f := mustParse(t, sample)

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
	f := mustParse(t, sample)

	require.NoError(t, f.Set("host1", opts(t, "port", "2222")))

	out := marshal(t, f)
	require.Contains(t, out, "	Port 2222\n")
	require.Contains(t, out, "HostName host1.example.com")
	require.NotContains(t, out, "	Port 22\n")
}

func TestSetAppendsToBlock(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set("host1", opts(t, "compression", "yes")))

	out := marshal(t, f)
	require.Contains(t, out, "Compression yes")
	require.Contains(t, out, "User git")
	require.NotContains(t, out, "Host host1.example.com")
}

func TestSetMultipleOptions(t *testing.T) {
	f := mustParse(t, sample)

	o := &ClientOptions{}
	require.NoError(t, o.Set("port", "2222"))
	require.NoError(t, o.Set("user", "admin"))
	require.NoError(t, f.Set("host1", o))

	out := marshal(t, f)
	require.Contains(t, out, "	Port 2222\n")
	require.Contains(t, out, "User admin")
	require.NotContains(t, out, "	Port 22\n")
}

func TestSetCreatesBlock(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set("host2", opts(t, "user", "admin")))

	out := marshal(t, f)
	require.Contains(t, out, "Host host2")
	require.Contains(t, out, "User admin")
}

func TestSetValidatesBeforeAction(t *testing.T) {
	f := mustParse(t, sample)

	require.Error(t, (&ClientOptions{}).Set("port", "abc"))
	require.Error(t, (&ClientOptions{}).Set("bogus", "yes"))
	require.Contains(t, marshal(t, f), "Port 22")
}

func TestMatchBlocksUntouched(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set("bastion", opts(t, "user", "admin")))

	out := marshal(t, f)
	require.Contains(t, out, "Match host bastion")
	require.Contains(t, out, "# inside match")
	require.Contains(t, out, "Host bastion")
}

func TestUseKeyChainAddsIgnoreUnknown(t *testing.T) {
	f := mustParse(t, `Host bastion
	UseKeyChain yes
`)

	out := marshal(t, f)
	require.Equal(t, `IgnoreUnknown UseKeychain

Host bastion
	UseKeyChain yes
`, out)
}

func TestUseKeyChainExistingIgnoreUnknown(t *testing.T) {
	f := mustParse(t, `Host bastion
	IgnoreUnknown UseKeychain
	UseKeyChain yes
`)

	require.Equal(t, `IgnoreUnknown UseKeychain

Host bastion
	UseKeyChain yes
`, marshal(t, f))
}

func TestUseKeyChainSetAddsIgnoreUnknown(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set("host1", opts(t, "usekeychain", "yes")))

	out := marshal(t, f)
	require.Contains(t, out, "IgnoreUnknown UseKeychain\n\n")
	require.Contains(t, out, "	UseKeyChain yes\n")
	require.Equal(t, map[string]string{"usekeychain": "UseKeychain"}, f.IgnoreUnknown)
}

func TestIgnoreUnknownMergesLines(t *testing.T) {
	f := mustParse(t, `IgnoreUnknown FooBar

Host bastion
	IgnoreUnknown UseKeychain,FooBar
	UseKeyChain yes
`)

	require.Equal(
		t,
		map[string]string{"foobar": "FooBar", "usekeychain": "UseKeychain"},
		f.IgnoreUnknown,
	)
	require.Equal(t, `IgnoreUnknown FooBar,UseKeychain

Host bastion
	UseKeyChain yes
`, marshal(t, f))
}

func TestIgnoreUnknownGetSet(t *testing.T) {
	f := mustParse(t, `Host bastion
	IgnoreUnknown UseKeychain,FooBar
`)

	value, err := f.Get("bastion", "ignoreunknown")
	require.NoError(t, err)
	require.Equal(t, "FooBar,UseKeychain", value)
}

func TestIgnoreUnknownSetGlobal(t *testing.T) {
	f := mustParse(t, sample)

	require.NoError(t, f.Set("host1", opts(t, "ignoreunknown", "FooBar")))

	require.Equal(t, map[string]string{"foobar": "FooBar"}, f.IgnoreUnknown)
	require.Contains(t, marshal(t, f), "IgnoreUnknown FooBar\n\n# personal hosts\n")
}

func TestIgnoreUnknownRoundTripStable(t *testing.T) {
	f := mustParse(
		t,
		"IgnoreUnknown UseKeychain\n\n# personal hosts\n\nHost host1\n\tHostName host1.example.com\n",
	)

	var out []byte
	for i := 0; i < 5; i++ {
		data, err := f.MarshalSSH()
		require.NoError(t, err)
		out = data
		require.NoError(t, f.UnmarshalSSH(out))
	}

	require.Equal(
		t,
		"IgnoreUnknown UseKeychain\n\n# personal hosts\n\nHost host1\n\tHostName host1.example.com\n",
		string(out),
	)
}
