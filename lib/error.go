package lib

import (
	"encoding/json"
)

type Error struct {
	Err error `json:"error"`
}

func (err *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Message string `json:"error"`
	}{
		Message: err.Err.Error(),
	})
}
