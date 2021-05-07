package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "host=127.0.0.1 port=9090 user=demidafonichev password='' dbname=master_thesis application_name=pgproxy sslmode=disable")
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

	// rows, err := db.Query("select email from userprofile")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	for rows.Next() {
	// 		var a string

	// 		if err := rows.Scan(&a); err != nil {
	// 			fmt.Println(err)
	// 		} else {
	// 			fmt.Println(a)
	// 		}
	// 	}
	// }

	// rows, err := db.Query("select first_name, last_name, email from userprofile")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	for rows.Next() {
	// 		var fn, ln, e string

	// 		if err := rows.Scan(&fn, &ln, &e); err != nil {
	// 			fmt.Println(err)
	// 		} else {
	// 			fmt.Println(fn, ln, e)
	// 		}
	// 	}
	// }

	rows, err := db.Queryx(`
	insert into userprofile(first_name, last_name, email)
	values
		('Demid', 'Afonichev', 'demidafonichev@gmail.com'),
		('Example', 'User', 'example@gmail.com'),
		('qwe', 'asd', 'zxc')`)
	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			fmt.Println("Inserted")
		}
	}
}
