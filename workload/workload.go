package workload

import (
	"fmt"
	"pgworkload/query"
	"pgworkload/schema"
	"time"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

type DBConf struct {
	Addr     string
	User     string
	Password string
	DbName   string
}

type Query struct {
	Query         string  `db:"query"`
	TotalExecTime float32 `db:"total_exec_time"`
}

func Start(connStr string) *query.QuerySet {
	schema.SyncTables(connStr)
	qs := query.ResetQuerySet()

	go workload(connStr, qs)

	return qs
}

func workload(connStr string, qs *query.QuerySet) {
	for {
		time.Sleep(3 * time.Second)

		statQueries := loadQueriesStats(connStr)
		regexpStatQueries := makeRegexpFromStatsQueries(statQueries)

		filteredQueries, filteredStatQueries := filterQueriesByRegexp(regexpStatQueries, qs)

		fmt.Printf("Filtered queries:\n")
		for _, q := range filteredQueries {
			fmt.Printf("%s\n", q[:15])
		}
		fmt.Println("-------------------------")

		fmt.Printf("Filtered stat queries:\n")
		for _, q := range filteredStatQueries {
			fmt.Printf("%s\n", q.Query[:15])
		}
		fmt.Println("-------------------------")
	}
}

func loadQueriesStats(connStr string) []*Query {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		glog.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Queryx("select query, total_exec_time from pg_stat_statements")
	if err != nil {
		glog.Fatalln(err)
	}

	var queries []*Query
	for rows.Next() {
		q := &Query{}
		if err := rows.StructScan(&q); err != nil {
			glog.Fatalln(err)
		}
		queries = append(queries, q)
	}
	return queries
}
