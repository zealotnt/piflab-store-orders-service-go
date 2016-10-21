package handlers

import (
	"encoding/json"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	"net/http"
)

type Index struct {
	Version string `json:"version"`
}

func IndexHandler(app *App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, c Context) {
		index := Index{"Gateway API 1.0.0"}

		json.NewEncoder(w).Encode(index)
	}
}
