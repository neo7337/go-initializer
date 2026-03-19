# How to Use the goini CLI

`goini` is the command-line companion to the web UI. It supports both fully interactive use and non-interactive flag-driven use for CI pipelines.

## Installation

```bash
git clone https://github.com/neo7337/go-initializer.git
cd go-initializer
make build-cli
sudo mv bin/goini /usr/local/bin/
```

Verify:

```bash
goini version
```

---

## Commands

### `goini new` — scaffold a project

Run with no flags to launch the interactive TUI:

```bash
goini new
```

The TUI walks you through every option with arrow-key navigation and a live preview.

**Non-interactive (flag-driven) usage:**

```bash
goini new \
  --name my-service \
  --module github.com/acme/my-service \
  --description "Order management service" \
  --type microservice \
  --framework gin \
  --go-version 1.25.0 \
  --addon cache=redis \
  --addon database=gorm \
  --docker
```

Any flag you omit will be prompted interactively.

**Available flags:**

| Flag | Description | Example |
| --- | --- | --- |
| `--name` | Project name (root directory) | `my-app` |
| `--module` | Go module path | `github.com/acme/my-app` |
| `--description` | Short description | `"Order service"` |
| `--type` | Project type | `microservice` |
| `--framework` | Framework | `gin` |
| `--go-version` | Go version | `1.25.0` |
| `--addon` | Addon in `category=value` format, repeatable | `--addon cache=redis` |
| `--docker` | Include a Dockerfile | _(flag, no value)_ |
| `--output` | Output directory | `./projects/my-app` |

The project is extracted directly to the output directory (default: `./<name>`).

---

### `goini list` — discover valid values

Use `goini list` sub-commands to see every valid value before scaffolding.

**List project types:**

```bash
goini list types
```

```
TYPE              LABEL
----              -----
api-server        API Server
cli-app           CLI Application
microservice      Microservice
simple-project    Simple Project
```

**List frameworks for a type:**

```bash
goini list frameworks --type microservice
```

```
FRAMEWORK
---------
echo
fiber
gin
golly
gokit
```

**List addons:**

```bash
goini list addons
```

```
CATEGORY    ADDON
--------    -----
cache       memcached
cache       redis
database    ent
database    gorm
other       cobra
other       logrus
other       zap
```

---

### Shell completion

Enable tab-completion for your shell:

```bash
# bash
goini completion bash > /etc/bash_completion.d/goini

# zsh
goini completion zsh > "${fpath[1]}/_goini"

# fish
goini completion fish | source

# PowerShell
goini completion powershell | Out-String | Invoke-Expression
```

---

## CI / scripted use

Because every option has a flag, `goini new` can run in a fully non-interactive mode inside CI:

```yaml
# GitHub Actions example
- name: Scaffold service
  run: |
    goini new \
      --name order-service \
      --module github.com/acme/order-service \
      --description "Order management" \
      --type api-server \
      --framework gin \
      --go-version 1.25.0 \
      --docker
```
