package solibase

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Changelog struct {
	Files []string
	Names []string `toml:"files"`
}

func NewChangelog(filename string) (Changelog, error) {
	var changelog Changelog
	if filename == "" {
		return changelog, errors.New("blank filename for changelog")
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return changelog, err
	}
	_, err = toml.Decode(string(data), &changelog)
	if err != nil {
		return changelog, err
	}

	for _, relativeFileName := range changelog.Names {
		if filepath.Ext(relativeFileName) != ".toml" {
			return changelog, errors.New("invalid file: " + relativeFileName)
		}
		changelog.Files = append(changelog.Files, filepath.Join(filepath.Dir(filename), relativeFileName))
	}

	return changelog, nil
}
