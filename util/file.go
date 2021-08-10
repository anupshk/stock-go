package util

import (
	"io/fs"
	"os"
	"path/filepath"
)

const CSV_DIR = "csv"

func GetAllCsvWalk() ([]string, error) {
	var files []string
	filepath.WalkDir(CSV_DIR, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".csv" {
			files = append(files, s)
		}
		return nil
	})
	return files, nil
}

func GetAllCsv() ([]string, error) {
	var files []string
	f, err := os.Open(CSV_DIR)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		fn := file.Name()
		if filepath.Ext(fn) == ".csv" {
			files = append(files, CSV_DIR+"/"+fn)
		}
	}
	return files, nil
}
