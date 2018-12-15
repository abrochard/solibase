package solibase

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Changelog struct {
	Files []string
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

	for i, relativeFileName := range changelog.Files {
		if filepath.Ext(relativeFileName) != ".toml" {
			return changelog, errors.New("invalid file: " + relativeFileName)
		}
		changelog.Files[i] = filepath.Join(filepath.Dir(filename), relativeFileName)
	}

	return changelog, nil
}
