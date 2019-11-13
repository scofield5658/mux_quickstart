package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

func download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date := vars["date"]
	file := vars["file"]

	var cwd, err = getStorageDir()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ErrorResponse{
			ErrCode: -1,
			ErrMsg:  err.Error(),
		})
		return
	}

	var storagePath = filepath.Join(cwd, date, file)

	_, err = os.Stat(storagePath)
	if err == nil {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		http.ServeFile(w, r, storagePath)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ErrorResponse{
		ErrCode: -1,
		ErrMsg:  "未找到资源",
	})
}
