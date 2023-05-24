package server

import (
	"time"

	"github.com/go-kit/log"
	"github.com/go-url-shortener/shortener"
)

type loggingMiddleware struct {
	logger log.Logger
	next   shortener.RedirectService
}

func NewLoggingMiddleware(logger log.Logger, service shortener.RedirectService) shortener.RedirectService {
	return loggingMiddleware{logger, service}
}

func (mw loggingMiddleware) Find(code string) (output *shortener.Redirect, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "find",
			"input", code,
			"output", output.URL,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.next.Find(code)
	return
}

func (mw loggingMiddleware) Store(redirect *shortener.Redirect) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "store",
			"input", redirect.Code+"|"+redirect.URL,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.Store(redirect)
	return
}
