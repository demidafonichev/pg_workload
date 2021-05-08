package workload

import "fmt"

// getConnStr creates connection string from DatabaseConfig
func getConnStr(dbconf DatabaseConfig) string {
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s application_name=pgproxy sslmode=disable",
		dbconf.Host, dbconf.Port, dbconf.User, getPassword(dbconf.Password), dbconf.DBName)
	return connstr
}

// getPassword returns password for db connection
func getPassword(pass string) string {
	if pass == "" {
		return "''"
	}
	return pass
}
