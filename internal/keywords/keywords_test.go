package keywords

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatePathAuthorizedKeysFile(t *testing.T) {
	dir := t.TempDir()
	existing := filepath.Join(dir, "keys")
	require.NoError(t, os.Mkdir(existing, 0o700))

	get := filepath.Join(dir, "sshd_config")
	require.NoError(t, os.WriteFile(get, []byte("AuthorizedKeysFile\n"), 0o600))

	// Absolute token-free path that exists is accepted.
	require.NoError(t, SSHD.Validate("authorizedkeysfile", existing+" "+get, ""))

	// Missing absolute path is rejected.
	require.Error(t, SSHD.Validate("authorizedkeysfile", filepath.Join(dir, "missing"), ""))

	// Relative paths resolve against the config directory when base is set.
	require.NoError(t, SSHD.Validate("authorizedkeysfile", "keys", dir))
	require.Error(t, SSHD.Validate("authorizedkeysfile", "missing", dir))

	// Without base, relative paths cannot be stat-ed; accepted.
	require.NoError(t, SSHD.Validate("authorizedkeysfile", "keys", ""))

	// Token-expanded paths cannot be stat-ed; accepted.
	require.NoError(t, SSHD.Validate("authorizedkeysfile", "%h/.ssh/authorized_keys", dir))
	require.NoError(t, SSHD.Validate("authorizedkeysfile", "keys/%u keys", dir))
	require.NoError(t, SSHD.Validate("authorizedkeysfile", "%%", dir))
}

func TestValidatePathHostKeyBanner(t *testing.T) {
	key := t.TempDir()
	hostKey := filepath.Join(key, "ssh_host_key")
	require.NoError(t, os.WriteFile(hostKey, []byte(""), 0o600))

	require.NoError(t, SSHD.Validate("hostkey", hostKey, ""))
	require.Error(t, SSHD.Validate("hostkey", filepath.Join(key, "missing_key"), ""))

	// Banner accepts "none" (disables the banner) since it is not an
	// absolute path sshd would stat.
	require.NoError(t, SSHD.Validate("banner", "none", key))
	require.NoError(t, SSHD.Validate("banner", hostKey, key))
	require.Error(t, SSHD.Validate("banner", filepath.Join(key, "missing_banner"), key))
}
