package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"

	"github.com/burkaydurdu/shortly/pkg/log"
)

type Shortly struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	VisitCount  int    `json:"visit_count"`
	Code        string `json:"code"`
}

type DB struct {
	ShortURL []Shortly `json:"short_url"`
}

func (d *DB) FindByCode(code string) *Shortly {
	for _, s := range d.ShortURL {
		if s.Code == code {
			return &s
		}
	}

	return nil
}

type ShortlyBase struct {
	FileName   string
	Log        *log.ShortlyLog
	MemoryPath string
}

func (s *ShortlyBase) InitialDB() (*DB, error) {
	path := fmt.Sprintf("%s/%s.json", s.MemoryPath, s.FileName)

	// File exist control
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		_, err = f.Write([]byte("{}"))

		if err != nil {
			return nil, err
		}

		// Close the new file.
		err = f.Close()

		if err != nil {
			return nil, err
		}

		s.Log.Zap(fmt.Sprintf("%s has created!", path))
	}

	return s.ReadFromFile()
}

func (s *ShortlyBase) SaveToFile(db *DB) error {
	filePermissionCode := 0600

	path := fmt.Sprintf("%s/%s.json", s.MemoryPath, s.FileName)

	fileData, _ := json.Marshal(db)

	return ioutil.WriteFile(path, fileData, fs.FileMode(filePermissionCode))
}

func (s *ShortlyBase) ReadFromFile() (*DB, error) {
	path := fmt.Sprintf("%s/%s.json", s.MemoryPath, s.FileName)

	// read in the local data
	file, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var DB = new(DB)

	err = json.Unmarshal(file, &DB)

	if err != nil {
		return nil, err
	}

	return DB, nil
}
