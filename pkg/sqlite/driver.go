package sqlite

import (
	"fmt"
	"github.com/abrochard/solibase/pkg/solibase"
)

type Driver struct {
	DB string
}

func (d *Driver) AddFlags(fs solibase.FlagSet) {
	fs.StringVar(&d.DB, "db", d.DB, "SQLite database")
}

func (d *Driver) Connect() error {
	return nil
}

func (d *Driver) CreateChangelogTableIfNotExists() error {
	return nil
}

func (d *Driver) Exec(query string) error {
	fmt.Printf("Exec: %s\n", query)
	return nil
}

func (d *Driver) ConditionApply(query string) (bool, error) {
	return true, nil
}

func (d *Driver) SaveChangeSet(name, metadata, hash string, rollback bool) error {
	return nil
}

func (d *Driver) LastChangeSet() (string, string, error) {
	return "", "", nil
}
