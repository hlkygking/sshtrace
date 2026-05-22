# sshtrace

> Utility to log and replay SSH session commands for audit trails

---

## Installation

```bash
go install github.com/yourusername/sshtrace@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/sshtrace.git && cd sshtrace && go build ./...
```

---

## Usage

Start logging an SSH session by wrapping your connection with `sshtrace`:

```bash
sshtrace record --output session.log user@remote-host
```

Replay a previously recorded session:

```bash
sshtrace replay --input session.log
```

List recorded sessions:

```bash
sshtrace list
```

**Example output:**

```
[2024-01-15 10:32:01] user@remote-host $ ls -la /etc
[2024-01-15 10:32:03] user@remote-host $ cat /etc/passwd
[2024-01-15 10:32:07] user@remote-host $ sudo systemctl status nginx
```

Sessions are stored in `~/.sshtrace/sessions/` by default. Use `--dir` to specify a custom path.

---

## Configuration

| Flag | Default | Description |
|------|---------|-------------|
| `--output` | stdout | Path to write session log |
| `--dir` | `~/.sshtrace` | Base directory for session storage |
| `--format` | `text` | Log format (`text` or `json`) |

---

## License

MIT © [yourusername](https://github.com/yourusername)