# relifevc

**Relifevc** is a CLI app built using [gognitive](https://github.com/sottey/gognitive), a Go client for the Limitless AI Pendant API.

This tool provides convenient ways to list, view, and export lifelog data collected by your Limitless pendant — perfect for personal analysis, journaling, or archiving.

---

## 🧪 Features

- List recent lifelog entries
- Export one or all entries as JSON or Markdown
- Generate metadata files for use in RAG/AI systems
- Daily automation-friendly (via cron, etc.)

---

## 🏁 Quick Start

```bash
cd example/relifevc
go run main.go list --limit 5
```

### Build It

```bash
go build -o relifevc
./relifevc list --limit 5
./relifevc export --save --all --format=json --dir=./backups
```

---

## ⚙️ Config

Create `~/.relifevc/config.json`:

```json
{
  "api_key": "sk-xxx",
  "timezone": "America/Los_Angeles",
  "export": {
    "format": "json",
    "dir": "./exports"
  }
}
```

Or pass a different config via:

```bash
./relifevc --config ./myconfig.json list
```

---

## 🗓 Example Cron Job

To export your lifelogs daily:

```cron
0 1 * * * /path/to/relifevc export --save --all >> ~/lifelog_export.log 2>&1
```

---

## 🛠 Powered By

- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [Gognitive](https://github.com/sottey/gognitive)

---

## 📄 License

MIT © 2025 [Sean Ottey](https://github.com/sottey)
