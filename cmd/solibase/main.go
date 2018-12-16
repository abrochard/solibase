package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/abrochard/solibase/pkg/mysql"
	"github.com/abrochard/solibase/pkg/solibase"
	"github.com/abrochard/solibase/pkg/sqlite"
)

func main() {
	os.Exit(runApp())
}

func runApp() int {
	fs := flag.NewFlagSet("solibase", flag.ExitOnError)

	var driverName string
	fs.StringVar(&driverName, "driver", driverName, "driver to use")
	var changelogFile string
	fs.StringVar(&changelogFile, "changelog", changelogFile, "location of the changelog")
	var rollback string
	fs.StringVar(&rollback, "rollback", rollback, "rollback to before the specified change")

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

	return solibase.Run(driver, changelog, rollback)
}

func registerDrivers() map[string]solibase.Driver {
	return map[string]solibase.Driver{
		"sqlite": &sqlite.Driver{},
		"mysql":  &mysql.Driver{},
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
