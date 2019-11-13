package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	var fileMd5 string
	file, handler, err := r.FormFile("file")
	defer file.Close()

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  "未找到文件表单字段",
		})
		return
	}

	storageDir, err := getStorageDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	var now = time.Now()
	var subStorageDir = now.Format("20060102")
	err = os.MkdirAll(filepath.Join(storageDir, subStorageDir), os.ModePerm)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	fileForHash, _, err := r.FormFile("file")
	fileMd5, err = getMD5ForFile(&fileForHash)

	storagePath := filepath.Join(storageDir, subStorageDir, fileMd5)
	fileInfo, err := os.Stat(storagePath)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FileUploadResponse{
			Filename:     handler.Filename,
			Size:         fileInfo.Size(),
			Path:         subStorageDir + "/" + fileMd5,
			ModifiedTime: TimeStamp(fileInfo.ModTime().Unix()),
		})
		return
	}

	f, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer f.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	if _, err = io.Copy(f, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FileUploadResponse{
		Filename:     handler.Filename,
		Size:         handler.Size,
		Path:         subStorageDir + "/" + fileMd5,
		ModifiedTime: TimeStamp(now.Unix()),
	})
}
