package main

import (
	"os"
	"path/filepath"
)

func getWorkDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, "tmp"), nil
}
