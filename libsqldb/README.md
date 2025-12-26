# PocketBase LibSQL (Turso) Database Module

This module allows you to use [Turso](https://turso.tech/) — a distributed, scalable SQLite database powered by libsql — as the backend database for [PocketBase](https://pocketbase.io/).

Turso provides edge-replicated SQLite databases with low-latency global access, making it an excellent choice for production PocketBase deployments that need scalability beyond a single local SQLite file.

## Features

- Seamless integration with PocketBase using the official `libsql-client-go` driver.
- Reuses PocketBase's existing SQLite query builder for full compatibility.
- Optional local storage for the auxiliary logs database (`auxiliary.db`) while using Turso for the main `data.db`.

## Installation

```bash
go get github.com/jim-ww/pocketbase-modules/libsqldb
```

## Usage
1. Create a Turso database and obtain its connection URL in the format:
`libsql://your-database-name.your-org.turso.io?authToken=your-auth-token`
(You can find this URL in the Turso dashboard.)
2. In your PocketBase main.go, use the provided DBConnect function:
```go
package main

import (
    "log"
    "os"

    "github.com/pocketbase/pocketbase"
    "github.com/jim-ww/pocketbase-modules/libsqldb"
)

func main() {
    app := pocketbase.NewWithConfig(pocketbase.Config{
        DBConnect: libsqldb.DBConnect(
            os.Getenv("DB_URL"),      // e.g. libsql://... ?authToken=...
            true,                     // storeLogsLocally – set to true to keep auxiliary.db local
        ),
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}
```

## Parameters for DBConnect

- `dbURL` (string): The full Turso libsql connection URL including the `authToken` query parameter.
- `storeLogsLocally` (bool):
    - `true` (recommended): The main database (data.db) uses Turso, while the auxiliary logs database (auxiliary.db) remains on the local filesystem. This avoids unnecessary remote writes for frequent log operations.
    - `false`: Both databases use Turso (everything remote).
