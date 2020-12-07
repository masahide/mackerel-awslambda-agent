package statefile

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"
)

// GetStatefiles get all data from state files.
func GetStatefiles(dir string) ([]byte, error) {
	data := map[string][]byte{}
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return xerrors.Errorf("filepath.Walk err path:%s err: %w", path, err)
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return xerrors.Errorf("filepath.Rel err path:%s err: %w", path, err)
			}
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return xerrors.Errorf("ioutil.ReadFile err path:%s err: %w", path, err)
			}
			data[relPath] = b

			return nil
		})
	if err != nil {
		return nil, xerrors.Errorf("filepath.Walk err: %w", err)
	}

	return json.Marshal(data)
}

// PutStatefiles put all data to state files.
func PutStatefiles(baseDir string, jsonBlob []byte) error {
	var data map[string][]byte
	if err := json.Unmarshal(jsonBlob, &data); err != nil {
		return xerrors.Errorf("json.Unmarshal err: %w", err)
	}
	for relPath, body := range data {
		dir := filepath.Join(baseDir, filepath.Dir(relPath))
		if err := os.MkdirAll(dir, 0755); err != nil {
			return xerrors.Errorf("os.MkdirAll statefile err: %w", err)
		}
		fullPath := filepath.Join(baseDir, relPath)
		if err := ioutil.WriteFile(fullPath, body, 0600); err != nil {
			return xerrors.Errorf("ioutil.WriteFile statefile path:%s err: %w", fullPath, err)
		}
		log.Printf("put state file: %s", fullPath)
	}

	return nil
}
