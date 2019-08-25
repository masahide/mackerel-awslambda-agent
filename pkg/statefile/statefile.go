package statefile

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetStatefiles get all data from state files
func GetStatefiles(dir string) ([]byte, error) {
	data := map[string][]byte{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Wrapf(err, "filepath.Walk err path:%s", path)
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return errors.Wrapf(err, "filepath.Rel err path:%s", path)
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return errors.Wrapf(err, "ioutil.ReadFile err path:%s", path)
			}
			data[relPath] = b
			return nil
		})
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
}

// PutStatefiles put all data to state files
func PutStatefiles(baseDir string, jsonBlob []byte) error {
	var data map[string][]byte
	if err := json.Unmarshal(jsonBlob, &data); err != nil {
		return errors.Wrap(err, "json.Unmarshal err")
	}
	for relPath, body := range data {
		dir := filepath.Join(baseDir, filepath.Dir(relPath))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return errors.Wrap(err, "os.MkdirAll")
		}
		fullPath := filepath.Join(baseDir, relPath)
		if err := ioutil.WriteFile(fullPath, body, 0644); err != nil {
			return errors.Wrapf(err, "ioutil.WriteFile path:%s", fullPath)
		}
	}
	return nil
}
