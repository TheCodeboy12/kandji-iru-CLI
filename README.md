# Kandji Iru CLI

A CLI for the [Kandji](https://www.kandji.io/) (Iru) Endpoint Management API. Manage devices, blueprints, users, tags, audit logs, and more from the command line.

## Requirements

- Go 1.21+
- A Kandji API token and your tenant base URL (or subdomain)

## Installation

```bash
go build -o kandji-iru-cli .
# or
go install .
```

## Configuration

The CLI resolves the **API token** in this order (first found wins):

1. **`--token` flag** — explicit override when you pass it
2. **System keyring** — macOS Keychain, Linux Secret Service, or Windows Credential Manager (recommended; token not on disk or in env)
3. **Config file** — `token:` in `~/.config/kandji-iru-cli/config.yaml`
4. **Environment** — `KANDJI_TOKEN`

Base URL is still from flags, env (`KANDJI_BASE_URL` / `KANDJI_SUBDOMAIN`), or config file.

If no token is found, the CLI suggests using the keyring first, then config or env.

### Create the config file

**Option A — Init with keyring (recommended):** store the token in the system keyring and keep the config file without secrets:

```bash
kandji-iru-cli init --keyring
# Creates config file (no token in file), prompts for API token and stores it in the keyring
# Then edit the config file to add base-url (or subdomain) only
```

You can also pipe the token or use env: `echo "YOUR_TOKEN" | kandji-iru-cli init --keyring` or `KANDJI_TOKEN=xxx kandji-iru-cli init --keyring`.

**Option B — Init with config file only:** create the config file and add the token in plain text:

```bash
kandji-iru-cli init
# Creates ~/.config/kandji-iru-cli/config.yaml; edit it and add your token and base-url (or subdomain)
```

Use `--config path/to/file` to use a different config path, or `--force` to overwrite an existing file.

### Config file example (`~/.config/kandji-iru-cli/config.yaml`)

```yaml
token: your-api-token   # omit if using keyring
base-url: https://your-tenant.api.kandji.io
# or (US): subdomain: your-tenant
```

For EU tenants, set `base-url` explicitly (e.g. `https://your-tenant.api.eu.kandji.io`).

### System keyring

Storing the token in the keyring is the most secure option (no plain-text token in config or env). The CLI checks the keyring **before** the config file and environment.

```bash
# Interactive setup (creates config + keyring entry)
kandji-iru-cli init --keyring

# Store token from flag, env, or stdin into keyring
kandji-iru-cli token store
echo "YOUR_API_TOKEN" | kandji-iru-cli token store

# Remove token from keyring
kandji-iru-cli token delete
```

Uses [go-keyring](https://github.com/zalando/go-keyring) (macOS Keychain, Linux Secret Service, Windows Credential Manager).

## Arbitrary query params (`--params`)

Any list or GET endpoint that supports query parameters accepts `--params` with a JSON object. Param names must match the API (snake_case). These are merged with (and can override) the normal flags.

**Endpoints that support `--params`:** `devices list`, `audit events list`, `users list`, `blueprints list`, `tags list`, `devices activity <id>`, `blueprint-routing activity`.

```bash
# Devices
kandji-iru-cli devices list --params '{"serial_number":"ABC123"}' -o json
kandji-iru-cli devices list --params '{"platform":"Mac","limit":5}' -o json

# Audit events
kandji-iru-cli audit events list --params '{"limit":100,"sort_by":"-occurred_at"}' -o json

# Users, blueprints, tags
kandji-iru-cli users list --params '{"email":"@company.com"}' -o json
kandji-iru-cli blueprints list --params '{"name":"Engineering"}' -o json
kandji-iru-cli tags list --params '{"search":"eng"}' -o json

# Device activity, blueprint-routing activity
kandji-iru-cli devices activity <device_id> --params '{"limit":50}' -o json
kandji-iru-cli blueprint-routing activity --params '{"limit":100}' -o json
```

## Pagination

List commands that return paged results will show a **next (and previous) page hint** on stderr when more data is available, so table output stays clean and scripts can ignore stderr.

- **Audit events** — cursor-based: use `--cursor=<value>` from the hint to get the next page.
- **Users** — cursor-based: use `--cursor=<value>`.
- **Blueprints** — offset-based: use `--offset=N --limit=M` from the hint.
- **Devices** — offset-based (300 per page max): use `--offset=N --limit=M`; the hint is shown when you received a full page.

With `-o json` or `--raw`, the full API response is printed (including `next` and `previous` for audit, blueprints, and users), so you can parse the cursor or URL yourself.

## Output formats

- **Table** (default) — Human-readable tables for list/get commands.
- **JSON** — Use `-o json` for pretty-printed JSON (re-marshaled from the CLI structs).
- **Raw** — Use `-o raw` or `--raw` to get the exact API response body (no parsing or re-encoding).


Examples:

```bash
kandji-iru-cli devices list -o json    # pretty JSON from CLI structs
kandji-iru-cli devices list -o raw      # exact API response bytes
kandji-iru-cli --raw audit events list
```

## Commands overview

| Area | Commands |
|------|----------|
| **Audit** | `audit events list` — list audit events (limit, sort, date range, cursor) |
| **Blueprints** | `blueprints list`, `blueprints get <id>`, `blueprints library-items <id>`, `blueprints templates` |
| **Blueprint routing** | `blueprint-routing get`, `blueprint-routing activity` |
| **Devices** | `devices list`, `devices get <id>`, `devices details <id>`, `devices action <id> <action>` |
| **Device info** | `devices activity`, `devices apps`, `devices library-items`, `devices parameters`, `devices status` |
| **Device notes** | `devices notes list/get/create/update/delete` |
| **Device update** | `devices update <id>` — PATCH user, asset tag, blueprint, tags (`--user`, `--asset-tag`, `--blueprint-id`, `--tags`, `--clear-asset-tag`, `--clear-tags`, or `--body`) |
| **Lost mode** | `devices lostmode get <id>`, `devices lostmode cancel <id>` |
| **Secrets** | `devices secrets <id> <type>` — types: `bypasscode`, `filevaultkey`, `unlockpin`, `recoverypassword` |
| **Tags** | `tags list`, `tags get <id>`, `tags create --name ...`, `tags update <id> --name ...`, `tags delete <id>` |
| **Settings** | `settings licensing` |
| **Token (keyring)** | `token store` — store API token in system keyring; `token delete` — remove from keyring |
| **Users** | `users list`, `users get <id>`, `users delete <id>` |
| **Shell completion** | `completion [bash\|zsh\|fish\|powershell]` |

## Examples

```bash
# List devices (table)
kandji-iru-cli devices list

# List devices as JSON
kandji-iru-cli devices list -o json

# List with arbitrary query params (JSON); param names must match API (snake_case, e.g. serial_number)
kandji-iru-cli devices list --params '{"serial_number":"ABC123"}' -o json
kandji-iru-cli devices list --params '{"platform":"Mac","limit":10}' -o json

# Get a device and its full details
kandji-iru-cli devices get <device_id>
kandji-iru-cli devices details <device_id>

# Run a device action (e.g. lock, restart, blank push)
kandji-iru-cli devices action <device_id> lock
kandji-iru-cli devices action <device_id> setname --body '{"DeviceName":"My Mac"}'

# List and create device notes
kandji-iru-cli devices notes list <device_id>
kandji-iru-cli devices notes create <device_id> --content "Deployed 2024-01-15"

# Update device (assign user, blueprint, tags)
kandji-iru-cli devices update <device_id> --blueprint-id <blueprint_uuid> --user <user_uuid>

# Raw JSON for any command
kandji-iru-cli --raw blueprints list
kandji-iru-cli audit events list -o json --limit 10
```

## Chaining commands

The CLI does not read arguments from stdin. To use the output of one command as input to another, use **command substitution** `$(...)` or **xargs**, and use **`-o json`** so you can parse output (e.g. with [jq](https://jqlang.github.io/jq/)).

**Get one device by user email:**
```bash
kandji-iru-cli devices get $(kandji-iru-cli devices list --user-email user@company.com -o json | jq -r '.[0].device_id')
```

**Run a command for every device matching a filter (e.g. send blank push to all Macs):**
```bash
kandji-iru-cli devices list --platform Mac -o json | jq -r '.[].device_id' | xargs -I {} kandji-iru-cli devices action {} blankpush
```

**Get full details for every device for a user (one JSON object per line):**
```bash
kandji-iru-cli devices list --user-email user@company.com -o json | jq -r '.[].device_id' | while read -r id; do kandji-iru-cli devices get "$id" -o json; done
```

**Update every device for a user (e.g. assign a blueprint):**
```bash
kandji-iru-cli devices list --user-email user@company.com -o json | jq -r '.[].device_id' | while read -r id; do kandji-iru-cli devices update "$id" --blueprint-id <blueprint_uuid>; done
```

**Add a note to every device in a blueprint:**
```bash
kandji-iru-cli devices list --blueprint-id <blueprint_uuid> -o json | jq -r '.[].device_id' | xargs -I {} kandji-iru-cli devices notes create {} --content "Deployed 2024-01-15"
```

**Extract fields from a list (e.g. device IDs only, or count):**
```bash
kandji-iru-cli devices list -o json | jq -r '.[].device_id'
kandji-iru-cli devices list -o json | jq 'length'
```

**Chain other resources (e.g. list blueprints, then get first one):**
```bash
kandji-iru-cli blueprints list -o json | jq -r '.[0].id' | xargs kandji-iru-cli blueprints get
```

## API reference

The CLI follows the [Iru Endpoint Management API](https://api-docs.kandji.io/) (and the included Postman collection). Rate limit: 10,000 requests/hour per tenant.

## License

See repository license.
