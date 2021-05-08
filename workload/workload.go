package workload

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	ConnStr  string
}

type Query struct {
	Query         string  `db:"query"`
	TotalExecTime float32 `db:"total_exec_time"`
}

// Start reads previous db tables state/requests saves new ones
// resets queryset of queries run through proxy
// starts workload analyzer
func Start(dbconf DatabaseConfig) *QuerySet {
	dbconf.ConnStr = getConnStr(dbconf)

	syncTables(dbconf)
	qs := ResetQuerySet()

	// analyze workload
	ws := &WorkloadAnalyzer{
		querySet: qs,
		dbconf:   dbconf,
	}
	go ws.service()

	return qs
}

type WorkloadAnalyzer struct {
	querySet *QuerySet
	dbconf   DatabaseConfig
}

// WorkloadAnalyzer.service starts workload analyzer to
// analyze queries and possible migrations
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

// WorkloadAnalyzer.getQueriesWithStats requests pg_stat_statements queries from db
// and compares them with queries run through proxy
// by matching queries to stat queries pattern
func (wa *WorkloadAnalyzer) getQueriesWithStats() ([]string, []*Query) {
	statQueries := wa.loadQueriesStats()
	regexpStatQueries := makeRegexpFromStatsQueries(statQueries)

	// filter proxied queries by stat queries pattern
	fQueries, fStatQueries := filterQueriesByRegexp(regexpStatQueries, wa.querySet)
	return fQueries, fStatQueries
}

// WorkloadAnalyzer.loadQueriesStats requests pg_stat_statements queries from db
func (wa *WorkloadAnalyzer) loadQueriesStats() []*Query {
	db, err := sqlx.Open("postgres", wa.dbconf.ConnStr)
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

// WorkloadAnalyzer.findOptimizations clones db and applies
// possible migrations to optimaze it
// After migrating database tests queries total exec time
// Migration set that leads to less total exec time applies to db
func (wa *WorkloadAnalyzer) findOptimizations([]string, []*Query) {
	fmt.Println("Optimizing")
}

func generateMigrationsSets() {}

func testMigrationsSet() {}

func cloneDB() {}
