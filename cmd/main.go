package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/log"

	"github.com/go-url-shortener/server"
)

func main() {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}
	srv := server.NewServer(logger)

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
		logger.Log(
			"started at port", "8080",
		)
		errs <- server.ListenAndServe()
	}()
	err := <-errs
	logger.Log(err)
}
