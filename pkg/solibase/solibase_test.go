package solibase

import (
	"testing"
)

func TestFindChangeIndex(t *testing.T) {
	changelog := Changelog{
		Names: []string{"changeset2.toml", "changeset1.toml", "changeset.toml"},
		Files: []string{"changeset2.toml", "changeset1.toml", "changeset.toml"},
	}

	name := "changeset1.toml"

	index := findChangeIndex(changelog, name)

	if index != 1 {
		t.Fatalf("Got wrong index back. Expected %d, got %d", 1, index)
	}
}

func TestApplyChangeset(t *testing.T) {
	driver := &DriverMock{
		ExecFunc: func(string) error {
			return nil
		},
		SaveChangeSetFunc: func(string, string, string) error {
			return nil
		},
	}
	changeset := Changeset{}

	err := applyChangeset(driver, changeset)
	if err != nil {
		t.Fatal(err)
	}
}
