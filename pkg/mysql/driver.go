package mysql

import (
	"database/sql"
	"fmt"

	"github.com/abrochard/solibase/pkg/solibase"

	_ "github.com/go-sql-driver/mysql"
)

type Driver struct {
	User     string
	Password string
	Host     string
	Port     int
	DB       string

	db *sql.DB
}

func (d *Driver) AddFlags(fs solibase.FlagSet) {
	fs.StringVar(&d.User, "user", d.User, "MySQL username (defaults to root)")
	fs.StringVar(&d.Password, "password", d.Password, "MySQL password (defaults to no password)")
	fs.StringVar(&d.Host, "host", d.Host, "MySQL host (defaults to localhost)")
	fs.IntVar(&d.Port, "port", d.Port, "MySQL port (defaults to 3306)")
	fs.StringVar(&d.DB, "db", d.DB, "MySQL database")
}

func (d *Driver) Connect() error {
	if d.User == "" {
		d.User = "root"
	}
	if d.Password != "" {
		d.Password = ":" + d.Password
	}
	if d.Host == "" {
		d.Host = "localhost"
	}
	if d.Port == 0 {
		d.Port = 3306
	}

	conn := fmt.Sprintf("%s%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.DB)
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Driver) CreateChangelogTableIfNotExists() error {
	rows, err := d.db.Query("SELECT * FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA=DATABASE() AND TABLE_NAME='SOLIBASE_CHANGELOG'")
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		// it exists
		return nil
	}

	query := `CREATE TABLE SOLIBASE_CHANGELOG (
	id BIGINT NOT NULL AUTO_INCREMENT,
	name TEXT NOT NULL,
	metadata TEXT NOT NULL,
	hash TEXT NOT NULL,
	rollback TINYINT NOT NULL DEFAULT '0',
	created_at DATETIME NOT NULL DEFAULT NOW(),
	updated_at DATETIME NOT NULL DEFAULT NOW(),
	PRIMARY KEY (id)
);`
	stmt, err := d.db.Prepare(query)
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
	stmt, err := d.db.Prepare("UPDATE SOLIBASE_CHANGELOG SET rollback=1, updated_at=NOW() WHERE name=? AND hash=?")
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
