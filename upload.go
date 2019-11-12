package main

import (
	"crypto/md5"
	"encoding/hex"
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

	hash := md5.New()

	if _, err = io.Copy(hash, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	hashInBytes := hash.Sum(nil)[:16]
	fileMd5 = hex.EncodeToString(hashInBytes)

	storageDir, err := getWorkDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	storagePath := filepath.Join(storageDir, fileMd5)
	fileInfo, err := os.Stat(storagePath)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(FileUploadResponse{
			Filename:     handler.Filename,
			Size:         fileInfo.Size(),
			MD5:          fileMd5,
			ModifiedTime: TimeStamp(fileInfo.ModTime().Unix()),
		})
		return
	}

	f, err := os.OpenFile(storagePath, os.O_WRONLY|os.O_CREATE, 0666)

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

	var now = time.Now()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(FileUploadResponse{
		Filename:     handler.Filename,
		Size:         handler.Size,
		MD5:          fileMd5,
		ModifiedTime: TimeStamp(now.Unix()),
	})
}
