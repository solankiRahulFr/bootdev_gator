# Gator

Gator is a command-line RSS feed aggregator written in Go. It allows users to register, follow RSS feeds, aggregate posts from those feeds, and browse them directly from the terminal.

---

## Requirements

Before running Gator, make sure you have the following installed:

* **Go (1.21 or later)**
  https://go.dev/doc/install

* **PostgreSQL**
  https://www.postgresql.org/download/

You also need a running PostgreSQL database.

---

## Installation

Install the `gator` CLI using `go install`:

```bash
go install github.com/YOUR_GITHUB_USERNAME/gator@latest
```

After installation, the `gator` command will be available in your terminal.

Make sure your `$GOPATH/bin` or `$HOME/go/bin` directory is in your `PATH`.

---

## Configuration

Gator requires a config file stored at:

```
~/.gatorconfig.json
```

Example configuration:

```json
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

### Config fields

* **db_url** – connection string for your PostgreSQL database
* **current_user_name** – the currently logged in user (managed automatically)

---

## Database Setup

Create the database:

```bash
createdb gator
```

Run the migrations to set up the tables:

```bash
goose postgres "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable" up
```

---

## Running the Program

Gator is a CLI application where you run commands like:

```bash
gator <command> [arguments]
```

---

## Example Commands

### Register a user

```bash
gator register alice
```

### Login

```bash
gator login alice
```

### Add a feed

```bash
gator addfeed "TechCrunch" https://techcrunch.com/feed/
```

### Follow a feed

```bash
gator follow https://techcrunch.com/feed/
```

### List feeds

```bash
gator feeds
```

### See feeds you follow

```bash
gator following
```

### Aggregate feeds (background worker)

```bash
gator agg 1m
```

This will collect posts from feeds every **1 minute** and store them in the database.

### Browse posts

```bash
gator browse
```

### Unfollow a feed

```bash
gator unfollow https://techcrunch.com/feed/
```

---

## Example RSS Feeds

You can try adding these feeds:

* https://techcrunch.com/feed/
* https://news.ycombinator.com/rss
* https://www.boot.dev/blog/index.xml

---

## Development

Clone the repository:

```bash
git clone https://github.com/YOUR_GITHUB_USERNAME/gator.git
cd gator
```

Build locally:

```bash
go build
```

Run commands:

```bash
./gator register testuser
```

---

## How It Works

Gator:

1. Stores RSS feeds in PostgreSQL
2. Periodically fetches feeds
3. Parses posts from RSS XML
4. Saves them to the database
5. Allows users to browse posts from the CLI

---

## License

MIT License
