package sqlite

import (
	"database/sql"

	"github.com/abrochard/solibase/pkg/solibase"

	_ "github.com/mattn/go-sqlite3"
)

type Driver struct {
	DB string

	db *sql.DB
}

func (d *Driver) AddFlags(fs solibase.FlagSet) {
	fs.StringVar(&d.DB, "db", d.DB, "SQLite database")
}

func (d *Driver) Connect() error {
	db, err := sql.Open("sqlite3", d.DB)
	d.db = db
	return err
}

func (d *Driver) CreateChangelogTableIfNotExists() error {
	rows, err := d.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='SOLIBASE_CHANGELOG'")
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// it exists
		return nil
	}

	// we need to create it
	stmt, err := d.db.Prepare(`CREATE TABLE SOLIBASE_CHANGELOG (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	metadata TEXT NOT NULL,
	hash TEXT NOT NULL,
	rollback TINYINT NOT NULL DEFAULT '0',
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
)`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func (d *Driver) Exec(query string) error {
	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}

func (d *Driver) ConditionApply(query string) (bool, error) {
	rows, err := d.db.Query(query)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return !rows.Next(), nil
}

func (d *Driver) SaveChangeSet(name, metadata, hash string) error {
	stmt, err := d.db.Prepare("INSERT INTO SOLIBASE_CHANGELOG (name, metadata, hash, rollback) VALUES (?,?,?,0)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, metadata, hash)
	return err
}

func (d *Driver) SaveRollback(name, hash string) error {
	stmt, err := d.db.Prepare("UPDATE SOLIBASE_CHANGELOG SET rollback=1, updated_at=CURRENT_TIMESTAMP WHERE name=? AND hash=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, hash)
	return err
}

func (d *Driver) LastChangeSet() (string, string, error) {
	rows, err := d.db.Query("SELECT name, hash FROM SOLIBASE_CHANGELOG WHERE rollback=0 ORDER BY id DESC LIMIT 1")
	if err != nil {
		return "", "", err
	}
	defer rows.Close()

	if !rows.Next() {
		return "", "", nil
	}

	var name, hash string
	err = rows.Scan(&name, &hash)
	return name, hash, err
}
