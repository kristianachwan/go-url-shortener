package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-url-shortener/server"
)

func main() {
	srv := server.NewServer()
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		server := &http.Server{
			Addr:    ":8080",
			Handler: srv,
		}
		errs <- server.ListenAndServe()
	}()
	err := <-errs
	log.Fatal(err)
}
