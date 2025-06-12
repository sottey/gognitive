# gognitive

**Gognitive** is a Go client for the [Limitless AI Pendant](https://limitless.ai) API, providing typed access to lifelog data for analysis, archiving, or intelligent retrieval.

> 🧠 Designed for developers building lifelogging tools, personal knowledge bases, or RAG systems.

---

## ✨ Features

- Authenticate and connect to the Limitless API
- List and filter lifelogs with full pagination support
- Fetch individual lifelogs by ID
- Work with structured data (timestamps, blockquotes, speakers, headings)
- Export-ready integration with CLI and automation tools

---

## 🚀 Installation

```bash
go get github.com/sottey/gognitive
```

If you're contributing or using locally:

```go
replace github.com/sottey/gognitive => ../gognitive
```

---

## 🧱 Client API

```go
client := gognitive.NewClient("sk-your-api-key")
```

### Available Methods

```go
ListLifelogs(limit int, cursor, date, start, end, timezone string) ([]Lifelog, string, error)

GetLifelog(id string) (*Lifelog, error)
```

---

## 🔡 Types

```go
type Lifelog struct {
	ID        string
	Title     string
	Markdown  string
	Contents  []ContentNode
	StartTime string
	EndTime   string
}

type ContentNode struct {
	Type              string
	Content           string
	StartTime         string
	EndTime           string
	SpeakerName       *string
	SpeakerIdentifier *string
	Children          []ContentNode
}
```

---

## 🧪 Example CLI: `relifevc`

A real-world CLI app using `gognitive` is included in the `example/relifevc` directory.

### Usage

```bash
cd example/relifevc
go run main.go list --limit 5
```

Or build it:

```bash
go build -o relifevc
./relifevc list --limit 5
./relifevc export --save --all --format=json --dir=./backups
```

### Config

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

## Troubleshooting
## To test the API calls using curl:
```
curl -X GET 'https://api.limitless.ai/v1/lifelogs?limit=100' -H 'Authorization: Bearer [YOUR API KEY HERE]' -H 'Content-Type: application/json'
```

---

## 📦 Roadmap

- [ ] POST/annotation support (if the API expands)
- [ ] Rate limiting + retry middleware
- [ ] Built-in local indexer for RAG-ready storage

---

## 🛡 License

MIT © 2025 [Sean Ottey](https://github.com/sottey)

---

## 🙋‍♂️ Questions?

Open an issue, fork the repo, or build something awesome with your lifelogs.
