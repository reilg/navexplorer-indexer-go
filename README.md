# NavExplorer Indexer

Indexer for NavExplorer.com

## Requirements

- Go 1.14 or higher
- Elasticsearch 8.2.2 or higher
- [navcoin-core](https://github.com/navcoin/navcoin-core)

## Local Setup

Copy `.env.example` contents to `.env` file.

```sh
# Install Go deps
go mod tidy

# Install dependency injection script
bin/di

# Run indexer
go run cmd/indexerd/main.go
```
