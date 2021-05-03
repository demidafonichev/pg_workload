package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "host=127.0.0.1 port=5432 user=demidafonichev password='' dbname=master_thesis application_name=pgproxy sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}

	// Set connections num
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(100)

	defer func() {
		db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := db.Query("select column_name from information_schema.columns where table_schema = 'public'")
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			var name string

			if err := rows.Scan(&name); err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(name)
			}
		}
	}
}
