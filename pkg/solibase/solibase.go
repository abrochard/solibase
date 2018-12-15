package solibase

import (
	"fmt"
)

type FlagSet interface {
	StringVar(p *string, name, value, usage string)
	IntVar(i *int, name string, value int, usage string)
}

type Driver interface {
	AddFlags(fs FlagSet)
	Connect() error
	CreateChangelogTableIfNotExists() error
	Exec(query string) error
	ConditionApply(query string) (bool, error)
	SaveChangeSet(name, metadata, hash string, rollback bool) error
	LastChangeSet() (string, string, error)
}

func Run(driver Driver, changelog Changelog) int {
	err := driver.Connect()
	if err != nil {
		panic(err)
	}

	err = driver.CreateChangelogTableIfNotExists()
	if err != nil {
		panic(err)
	}

	name, hash, err := driver.LastChangeSet()
	if err != nil {
		panic(err)
	}

	if name == changelog.Files[len(changelog.Files)-1] {
		fmt.Println("Already up to date")
		return 0
	}

	if name == "" {
		fmt.Println("No changes applied found, running all changelog")
		err := ApplyChangelog(driver, changelog)
		if err != nil {
			panic(err)
		}
		fmt.Println("Done")
		return 0
	}

	index := len(changelog.Files) - 1
	for {
		if changelog.Files[index] == name {
			break
		}

		index--
		if index < 0 {
			panic("Last applied change not found")
		}
	}

	c, err := NewChangeset(changelog.Files[index])
	if c.Hash != hash {
		fmt.Println("Wrong hash of last change applied")
		return 1
	}

	err = ApplyChangelog(driver, Changelog{Files: changelog.Files[index+1:]})
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")

	return 0
}

func ApplyChangelog(driver Driver, changelog Changelog) error {
	for _, filename := range changelog.Files {
		c, err := NewChangeset(filename)
		if err != nil {
			return err
		}

		err = ApplyChangeset(driver, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyChangeset(driver Driver, changeset Changeset) error {
	if changeset.Change.Condition != "" {
		valid, err := driver.ConditionApply(changeset.Change.Condition)
		if err != nil {
			return err
		}
		if !valid {
			// doesn't meet condition
			return nil
		}
	}

	err := driver.Exec(changeset.Change.SQL)
	if err != nil {
		return err
	}

	return driver.SaveChangeSet(changeset.Name, changeset.MetadataJSON, changeset.Hash, false)
}

func RollbackChangeset(driver Driver, changeset Changeset) error {
	if changeset.Rollback.Condition != "" {
		valid, err := driver.ConditionApply(changeset.Rollback.Condition)
		if err != nil {
			return err
		}
		if !valid {
			// doesn't meet condition
			return nil
		}
	}

	err := driver.Exec(changeset.Rollback.SQL)
	if err != nil {
		return err
	}

	return driver.SaveChangeSet(changeset.Name, changeset.MetadataJSON, changeset.Hash, true)
}
