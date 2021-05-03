package workload

import "pgworkload/schema"

type DBConf struct {
	Addr     string
	User     string
	Password string
	DbName   string
}

func Start(connStr string) {
	schema.SyncTables(connStr)
}
