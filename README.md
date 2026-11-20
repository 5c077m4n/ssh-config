# ssh-config

Programmatic SSH configuration via a CLI: `sshc` for the client (`~/.ssh/config`) and `sshdc` for the server (`/etc/ssh/sshd_config`).

## Install

```
go install github.com/5c077m4n/ssh-config/cmd/sshc@latest
go install github.com/5c077m4n/ssh-config/cmd/sshdc@latest
```

## Usage

Client settings use dot notation (`host.key`); the key part is case-insensitive:

```
sshc 'host1.user' git       # set User git inside the Host host1 block
sshc 'host1.user'           # print the value for that key
sshc '*.example.com.user'   # block patterns may contain wildcards
```

Server settings are flat (`key value`):

```
sshdc permitrootlogin no
sshdc permitrootlogin
```

## Validation

Every getter and setter is validated before the config file is read or written:

- Keys must be known SSH keywords (case-insensitive; unknown keys are rejected).
- Values are type-checked against the key's kind: `yes`/`no` booleans, enum membership, positive integers, or non-empty strings.

## Behavior notes

- Client config: `host.key` addresses the first `Host` block whose patterns match `host` (first match wins, like ssh). A setter targets a block by exact pattern and creates a new `Host` block when none matches. `Match` blocks are parsed and preserved but never addressable.
- Server config: getters and setters apply to the main section only; `Match` sections are preserved verbatim and never modified.
- Client config path: `$SSHC_CONFIG` or `~/.ssh/config`. Server config path: `$SSHDC_CONFIG` or `/etc/ssh/sshd_config`.
- `Include` directives are not followed.
