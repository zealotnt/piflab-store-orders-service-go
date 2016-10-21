package lib

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func Logger(inner HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		params := mux.Vars(r)
		get_params := r.URL.Query()

		inner(w, r, Context{Params: params, GetParams: get_params})

		log.Printf(
			"%s %q\t%s",
			r.Method,
			r.URL.String(),
			time.Since(start),
		)
	}
}
