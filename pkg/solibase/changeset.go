package solibase

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Changeset struct {
	Name         string
	Hash         string
	Metadata     map[string]string
	MetadataJSON string
	Change       struct {
		SQL       string
		Condition string
	}
	Rollback struct {
		SQL       string
		Condition string
	}
}

func NewChangeset(name, filename string) (Changeset, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Changeset{}, err
	}

	var changeset Changeset
	_, err = toml.Decode(string(data), &changeset)
	if err != nil {
		return Changeset{}, err
	}

	changeset.Name = name

	// set hash
	h := sha256.New()
	r := bytes.NewReader(data)
	_, err = io.Copy(h, r)
	if err != nil {
		return Changeset{}, err
	}
	changeset.Hash = hex.EncodeToString(h.Sum(nil))

	// set metadata as JSON string
	metadataJSON, err := json.Marshal(changeset.Metadata)
	if err != nil {
		return Changeset{}, err
	}
	changeset.MetadataJSON = string(metadataJSON)

	return changeset, nil
}
