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

	ws := &WorkloadAnalyzer{
		querySet: qs,
		conn:     connStr,
	}
	go ws.service()

	return qs
}

type WorkloadAnalyzer struct {
	querySet *query.QuerySet
	conn     string
}

func (wa *WorkloadAnalyzer) service() {
	for {
		time.Sleep(3 * time.Second)

		filteredQueries, filteredStatQueries := wa.getQueriesWithStats()

		fmt.Printf("Filtered queries:\n")
		for _, q := range filteredQueries {
			fmt.Printf("%s\n", q)
		}
		fmt.Println("-------------------------")
		fmt.Printf("Filtered stat queries:\n")
		for _, q := range filteredStatQueries {
			fmt.Printf("%s\n", q.Query)
		}
		fmt.Println("-------------------------")

		wa.findOptimizations(filteredQueries, filteredStatQueries)
	}
}

func (wa *WorkloadAnalyzer) getQueriesWithStats() ([]string, []*Query) {
	statQueries := wa.loadQueriesStats()
	regexpStatQueries := makeRegexpFromStatsQueries(statQueries)

	fQueries, fStatQueries := filterQueriesByRegexp(regexpStatQueries, wa.querySet)
	return fQueries, fStatQueries
}

func (wa *WorkloadAnalyzer) loadQueriesStats() []*Query {
	db, err := sqlx.Open("postgres", wa.conn)
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

func (wa *WorkloadAnalyzer) findOptimizations([]string, []*Query) {
	fmt.Println("Optimizing")
}
