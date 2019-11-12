package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Path("/upload").
		Methods("POST").
		HandlerFunc(FileUpload)

	storageDir, err := getWorkDir()
	if err != nil {
		log.Fatal("服务启动失败, 文件存储区不存在")
		return
	}

	fmt.Println("Starting")
	log.Fatal(http.ListenAndServe(":8080", router))
}
