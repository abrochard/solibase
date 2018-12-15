package main

import (
	_ "database/sql"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/abrochard/solibase/pkg/solibase"
	"github.com/abrochard/solibase/pkg/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	os.Exit(runApp())

	// load changelog

	// connect to db
	// db, err := sql.Open("sqlite3", "./foo.db")
	// if err != nil {
	// 	panic(err)
	// }
	// rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name=?;", "users")
	// if err != nil {
	// 	panic(err)
	// }
	// cols, _ := rows.Columns()
	// fmt.Printf("cols %+v\n", cols)
	// if rows.Next() {
	// 	fmt.Printf("rows %+v\n", rows)
	// }
	// rows.Close()

	// var sql interface{}

	// checkChangelogTable(&sql)

	// // get the last changeset applied
	// lastChangeName := "fakechange.toml"

	// // find that changeset in the list of changes
	// index := len(changelog.Files) - 1
	// for {
	// 	if changelog.Files[index] == lastChangeName {
	// 		index++
	// 		break
	// 	}

	// 	index--
	// 	if index < 0 {
	// 		panic("Last change not found in changelog")
	// 	}
	// }

	// if index == len(changelog.Files) {
	// 	fmt.Println("Already up to date")
	// 	return
	// }

	// // apply them from here
	// toApply := changelog.Files[index:]
	// for _, filename := range toApply {
	// 	_, err := solibase.NewChangeset(filename)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// err = c.ExecChange(&sql)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

func runApp() int {
	fs := flag.NewFlagSet("solibase", flag.ExitOnError)

	var driverName string
	fs.StringVar(&driverName, "driver", driverName, "driver to use")
	var changelogFile string
	fs.StringVar(&changelogFile, "changelog", changelogFile, "location of the changelog")

	fsx := solibase.FlagSetGenerator{FS: fs}

	drivers := registerDrivers()

	for name, d := range drivers {
		d.AddFlags(fsx.New(name))
	}

	parseEnv("solibase", fs)

	fs.Parse(os.Args[1:])

	changelog, err := solibase.NewChangelog(changelogFile)
	if err != nil {
		panic(err)
	}

	driver, ok := drivers[driverName]
	if !ok {
		fmt.Println("Must specify a driver name")
		return 1
	}

	return solibase.Run(driver, changelog)
}

func registerDrivers() map[string]solibase.Driver {
	return map[string]solibase.Driver{
		"sqlite": &sqlite.Driver{},
	}
}

func parseEnv(namespace string, fs *flag.FlagSet) {
	prefix := strings.ToUpper(namespace) + "_"
	fs.VisitAll(func(f *flag.Flag) {
		if v, ok := os.LookupEnv(prefix + strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))); ok {
			f.Value.Set(v)
		}
	})
}
