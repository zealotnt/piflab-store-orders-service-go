package lib

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request, Context)

func (h HandlerFunc) ServeHTTPC(w http.ResponseWriter, r *http.Request, c Context) {
	h(w, r, c)
}

type MiddlewareFunc func(HandlerFunc) HandlerFunc

type Handler func(*App) HandlerFunc
