# ssh-config

Programmatic SSH configuration via a CLI: `sshc` for the client
(`~/.ssh/config`) and `sshdc` for the server (`/etc/ssh/sshd_config`).

## Install

```
go install github.com/5c077m4n/ssh-config/cmd/sshc@latest
go install github.com/5c077m4n/ssh-config/cmd/sshdc@latest
```

## Usage examples

### Client (`sshc`)

Given this `~/.ssh/config`:

```
# personal hosts

Host bastion
	HostName bastion.example.com
	User admin
	Port 22

Host *.lab
	User deploy
	ForwardAgent yes
```

Set `User git` in the `Host bastion` block, then read it back:

```
$ sshc 'bastion.user' git
$ sshc 'bastion.user'
git
```

Keys and values are case-insensitive — these do the same thing:

```
$ sshc 'BASTION.Port' 2222
$ sshc 'bastion.port'
2222
```

Wildcards resolve like ssh itself (first match wins):

```
$ sshc 'db1.lab.forwardagent'
yes
```

A setter creates a new `Host` block when no block matches:

```
$ sshc 'web01.hostname' web01.internal
$ sshc 'web01.hostname'
web01.internal
```

Invalid input is rejected before the file is touched:

```
$ sshc 'bastion.bogus' yes
sshc: unknown key "bogus"
$ sshc 'bastion.port' abc
sshc: value "abc" for Port must be a positive integer
```

### Server (`sshdc`)

```
$ sshdc permitrootlogin no
$ sshdc permitrootlogin
no
$ sshdc maxauthtries 4
$ sshdc clientaliveinterval 120
```

Enum and integer values are validated:

```
$ sshdc permitrootlogin maybe
sshdc: value "maybe" for PermitRootLogin must be one of yes, no, prohibit-password, forced-commands-only
```

## Validation

Every getter and setter is validated before the config file is read or written:

- Keys must be known SSH keywords (case-insensitive; unknown keys are rejected).
- Values are type-checked against the key's kind: `yes`/`no` booleans, enum
  membership, positive integers, non-empty strings, or paths: path-valued
  keys (`AuthorizedKeysFile`, `Banner`, `HostKey`) accept non-empty values, but
  every token-free path in the value must exist: absolute paths are stat-ed
  directly and relative paths are joined with the server config's directory
  (via `sshdc`). `%`-tokens and the special value `none` are skipped — they are
  resolved by `sshd` at runtime.

## Behavior notes

- Client config: `host.key` addresses the first `Host` block whose patterns
  match `host` (first match wins, like ssh). A setter targets a block by exact
  pattern and creates a new `Host` block when none matches. `Match` blocks are
  parsed and preserved but never addressable.
- Server config: getters and setters apply to the main section only; `Match`
  sections are preserved verbatim and never modified.
- Client config path: `--config-path` flag or `~/.ssh/config`. Server config path:
  `--config-path` flag or `/etc/ssh/sshd_config`.
- `Include` directives are not followed.
