package schema

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
)

type Table struct {
	Name string
	Cols []*Column
}

type Column struct {
	TableName string `db:"table_name"`
	Name      string `db:"column_name"`
	Type      string `db:"data_type"`
}

var Tables []*Table

func SyncTables(connStr string) {
	tables, err := getTalbes()
	if err == nil {
		glog.Info("Read tables from file...")
	} else {
		glog.Info("No serialized tables found, requesting from db...")

		tables, err = readTablesFromDB(connStr)
		fmt.Println(tables)
		if err != nil {
			glog.Infof("Error reading tables form db: %s\n", err)
		}

		tmap := map[string][]*Table{"tables": tables}

		err = saveTables(tmap)
		if err != nil {
			glog.Fatalln("Error saving tables to file")
		}
	}

	Tables = tables
}

func readTablesFromDB(connStr string) ([]*Table, error) {
	fmt.Println(connStr)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		glog.Fatalln(err)
	}
	defer db.Close()

	rows, err := db.Queryx("select table_name, column_name, data_type from information_schema.columns where table_schema='public' and table_name <> 'pg_stat_statements'")
	if err != nil {
		glog.Fatalln(err)
	}

	var cols []*Column
	for rows.Next() {
		col := &Column{}
		if err := rows.StructScan(&col); err != nil {
			glog.Fatalln(err)
		}
		cols = append(cols, col)
	}

	tables := combineColumnsToTables(cols)

	return tables, nil
}

func combineColumnsToTables(cols []*Column) []*Table {
	tmap := map[string]*Table{}
	for _, col := range cols {
		tn := col.TableName
		_, inmap := tmap[tn]

		if !inmap {
			tmap[tn] = &Table{tn, []*Column{col}}
		} else {
			tmap[tn].Cols = append(tmap[tn].Cols, col)
		}
	}

	tables := []*Table{}
	for _, t := range tmap {
		tables = append(tables, t)
	}

	return tables
}
