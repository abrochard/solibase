package mysql

import (
	"database/sql"

	"github.com/abrochard/solibase/pkg/solibase"
)

type Driver struct {
	User string
	Host string
	Port int
	DB   string

	db *sql.DB
}

const query = `CREATE TABLE SOLIBASE_CHANGELOG (
	id BIGINT NOT NULL AUTO_INCREMENT,
	name TEXT NOT NULL,
	metadata TEXT NOT NULL,
	hash TEXT NOT NULL,
	rollback TINYINT NOT NULL DEFAULT '0',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (id)
);`

func (d *Driver) AddFlags(fs solibase.FlagSet) {

}

func (d *Driver) Connect() error {
	return nil
}

func (d *Driver) CreateChangelogTableIfNotExists() error {
	return nil
}

func (d *Driver) Exec(query string) error {
	return nil
}

func (d *Driver) ConditionApply(query string) (bool, error) {
	return true, nil
}

func (d *Driver) SaveChangeSet(name, metadata, hash string, rollback bool) error {
	return nil
}

func (d *Driver) SaveRollback(name, hash string) error {
	return nil
}

func (d *Driver) LastChangeSet() (string, string, error) {
	return "", "", nil
}
