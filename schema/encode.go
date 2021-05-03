package schema

import (
	"os"

	"github.com/BurntSushi/toml"
)

func getTalbes() ([]*Table, error) {
	tmap := map[string][]*Table{}
	if _, err := toml.DecodeFile("schema/tables.toml", &tmap); err != nil {
		return nil, err
	}
	return tmap["tables"], nil
}

func saveTables(tmap map[string][]*Table) error {
	f, err := os.Create("schema/tables.toml")
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
