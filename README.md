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

The CLI reads configuration from (in order of precedence):

1. **Flags** — `--token`, `--base-url` or `--subdomain`
2. **Environment** — `KANDJI_TOKEN`, `KANDJI_BASE_URL` or `KANDJI_SUBDOMAIN`
3. **Config file** — `~/.kandji.yaml` (or `--config path/to/file`)

### Config file example (`~/.kandji.yaml`)

```yaml
token: your-api-token
# One of:
base-url: https://your-tenant.api.kandji.io
# or (US):
subdomain: your-tenant
```

For EU tenants, set `base-url` explicitly (e.g. `https://your-tenant.api.eu.kandji.io`).

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
| **Users** | `users list`, `users get <id>`, `users delete <id>` |
| **Shell completion** | `completion [bash\|zsh\|fish\|powershell]` |

## Examples

```bash
# List devices (table)
kandji-iru-cli devices list

# List devices as JSON
kandji-iru-cli devices list -o json

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

## API reference

The CLI follows the [Iru Endpoint Management API](https://api-docs.kandji.io/) (and the included Postman collection). Rate limit: 10,000 requests/hour per tenant.

## License

See repository license.
