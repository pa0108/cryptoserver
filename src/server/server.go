package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"./crypto"
)

const (
	API_KEY    = ""
	API_SECRET = ""
)

func main() {
	l := log.New(os.Stdout, "Crypto API", log.LstdFlags)
	client := crypto.NewClient(API_KEY, API_SECRET)

	//get data from Crypto server on regular intervals
	go func() {
		client.GetData()
	}()

	//Create handler for http request
	h := crypto.NewCurrencyHandler()
	sm := http.NewServeMux()
	sm.Handle("/currency/", h)

	s := http.Server{
		Addr:         ":8080",
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		log.Fatal(s.ListenAndServe().Error())
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	l.Println("Recieved notification signal", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
