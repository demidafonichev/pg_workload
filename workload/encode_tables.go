package workload

import (
	"os"

	"github.com/BurntSushi/toml"
)

// getTables loads tables from config file
func getTalbes() ([]*Table, error) {
	tmap := map[string][]*Table{}
	if _, err := toml.DecodeFile("workload/tables.toml", &tmap); err != nil {
		return nil, err
	}
	return tmap["tables"], nil
}

// saveTables encode db tables config to file
func saveTables(tmap map[string][]*Table) error {
	f, err := os.Create("workload/tables.toml")
	if err != nil {
		return err
	}
	if err = toml.NewEncoder(f).Encode(tmap); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return nil
}
