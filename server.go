package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*20, "the duration for which the server gracefully wait for existing connections to finish")
	flag.Parse()

	router := mux.NewRouter()

	router.Path("/upload").
		Methods("POST").
		HandlerFunc(FileUpload)

	router.HandleFunc("/download/{date}/{file}", download).Methods("GET")

	var port = getEnv("MUX_SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         "0.0.0.0:" + port,
		WriteTimeout: time.Second * 20,
		ReadTimeout:  time.Second * 20,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	server.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
