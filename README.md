# revue

A TUI dashboard for GitHub pull request reviews and CI status.

## Features

- View PRs where you're requested as reviewer
- View your own open PRs
- CI status at a glance (✓ / ✗ / ⏳)
- Open PRs in browser with Enter
- Tab navigation between views

## Requirements

- [gh CLI](https://cli.github.com/) installed and authenticated (`gh auth login`)

## Install

```bash
go install github.com/ryoh827/revue@latest
```

## Usage

```bash
revue
```

### Key Bindings

| Key | Action |
|-----|--------|
| `j` / `k` | Navigate up/down |
| `Tab` / `h` / `l` | Switch tab |
| `Enter` | Open PR in browser |
| `r` | Refresh |
| `q` / `Ctrl+C` | Quit |
