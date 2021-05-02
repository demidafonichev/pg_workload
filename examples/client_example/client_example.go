package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "host=127.0.0.1 port=9090 user=demidafonichev password=postgres dbname=master_thesis sslmode=disable")
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

	rows, err := db.Query("select name from customer_account")
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
