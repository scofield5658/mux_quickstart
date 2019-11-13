package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func getStorageDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	var storageDir = getEnv("MUX_STORE_PATH")
	if storageDir == "" {
		storageDir = filepath.Join(cwd, "tmp")
	}

	err = os.MkdirAll(storageDir, os.ModePerm)
	if err != nil {
		return "", err
	}
	return storageDir, nil
}

func getEnv(key string) string {
	return os.Getenv(key)
}

func getMD5ForFile(file *multipart.File) (string, error) {
	hash := md5.New()
	defer (*file).Close()

	if _, err := io.Copy(hash, *file); err != nil {
		return "", err
	}

	hashInBytes := hash.Sum(nil)[:16]
	fileMd5 := hex.EncodeToString(hashInBytes)
	return fileMd5, nil
}
