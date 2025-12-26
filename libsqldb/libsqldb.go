package libsqldb

import (
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

// register the libsql driver to use the same query builder
// implementation as the already existing sqlite3 builder
func init() {
	dbx.BuilderFuncMap["libsql"] = dbx.BuilderFuncMap["sqlite3"]
}

func DBConnect(dbURL string, storeLogsLocally bool) core.DBConnectFunc {
	return func(dbPath string) (*dbx.DB, error) {
		if strings.Contains(dbPath, "data.db") || !storeLogsLocally {
			return dbx.Open("libsql", dbURL)
		}
		// optionally for the logs (aka. pb_data/auxiliary.db) use the default local filesystem driver
		return core.DefaultDBConnect(dbPath)
	}
}
