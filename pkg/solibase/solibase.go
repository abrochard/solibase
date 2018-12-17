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
	SaveChangeSet(name, metadata, hash string) error
	SaveRollback(name, hash string) error
	LastChangeSet() (string, string, error)
}

func Run(driver Driver, changelog Changelog, rollback string) int {
	err := driver.Connect()
	if err != nil {
		panic(err)
	}

	err = driver.CreateChangelogTableIfNotExists()
	if err != nil {
		panic(err)
	}

	if rollback != "" {
		fmt.Println("Starting a rollback to before " + rollback)
		bottom := findChangeIndex(changelog, rollback)
		if bottom < 0 {
			panic("rollback target not found in changelog")
		}
		lastChange, _, err := driver.LastChangeSet()
		top := findChangeIndex(changelog, lastChange)
		if top < 0 {
			panic("Last applied change not found")
		}
		if bottom > top {
			panic("You are trying to rollback something newer than the last applied change")
		}

		err = rollbackChangelog(driver, Changelog{Files: changelog.Files[bottom:top], Names: changelog.Names[bottom:top]})
		if err != nil {
			panic(err)
		}
		fmt.Println("Done")
		return 0
	}

	name, hash, err := driver.LastChangeSet()
	if err != nil {
		panic(err)
	}

	if name == changelog.Names[len(changelog.Files)-1] {
		fmt.Println("Already up to date")
		return 0
	}

	if name == "" {
		fmt.Println("No changes applied found, running all changelog")
		err := applyChangelog(driver, changelog)
		if err != nil {
			panic(err)
		}
		fmt.Println("Done")
		return 0
	}

	index := findChangeIndex(changelog, name)
	if index < 0 {
		panic("Last applied change not found")
	}

	c, err := NewChangeset(changelog.Names[index], changelog.Files[index])
	if c.Hash != hash {
		fmt.Println("Wrong hash of last change applied")
		return 1
	}

	err = applyChangelog(driver, Changelog{Files: changelog.Files[index+1:], Names: changelog.Names[index+1:]})
	if err != nil {
		panic(err)
	}

	fmt.Println("Done")

	return 0
}

func findChangeIndex(changelog Changelog, name string) int {
	index := len(changelog.Files) - 1
	for {
		if changelog.Names[index] == name {
			break
		}

		index--
		if index < 0 {
			return -1
		}
	}
	return index
}

func applyChangelog(driver Driver, changelog Changelog) error {
	for i, filename := range changelog.Files {
		c, err := NewChangeset(changelog.Names[i], filename)
		if err != nil {
			return err
		}

		err = applyChangeset(driver, c)
		if err != nil {
			return err
		}
	}

	return nil
}

func applyChangeset(driver Driver, changeset Changeset) error {
	if changeset.Change.Condition != "" {
		valid, err := driver.ConditionApply(changeset.Change.Condition)
		if err != nil {
			return err
		}
		if !valid {
			// doesn't meet condition
			fmt.Printf("Changeset %s doesn't meet its condition\n", changeset.Name)
			return nil
		}
	}

	err := driver.Exec(changeset.Change.SQL)
	if err != nil {
		return err
	}

	fmt.Printf("Applied change %s\n", changeset.Name)
	return driver.SaveChangeSet(changeset.Name, changeset.MetadataJSON, changeset.Hash)
}

func rollbackChangelog(driver Driver, changelog Changelog) error {
	i := len(changelog.Names) - 1
	for {
		c, err := NewChangeset(changelog.Names[i], changelog.Files[i])
		if err != nil {
			return err
		}

		err = rollbackChangeset(driver, c)
		if err != nil {
			return err
		}

		i--
		if i < 0 {
			return nil
		}
	}
}

func rollbackChangeset(driver Driver, changeset Changeset) error {
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

	fmt.Printf("Rolled back change %s\n", changeset.Name)
	return driver.SaveRollback(changeset.Name, changeset.Hash)
}
