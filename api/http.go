package api

import (
	"io/ioutil"
	"log"
	"net/http"

	js "github.com/go-url-shortener/serializer/json"
	"github.com/go-url-shortener/shortener"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type handler struct {
	redirectService shortener.RedirectService
}

type RedirectHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
}

func NewHandler(redirectService shortener.RedirectService) RedirectHandler {
	return &handler{redirectService}
}

func setupResponse(w http.ResponseWriter, body []byte, statusCode int) {
	w.Header().Set("Content-Type", "json")
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

func (h *handler) serializer() shortener.RedirectSerializer {
	return &js.Redirect{}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

func (h *handler) Post(w http.ResponseWriter, r *http.Request) {
	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirect, err := h.serializer().Decode(requestBody)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	responseBody, err := h.serializer().Encode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	setupResponse(w, responseBody, http.StatusCreated)
}
