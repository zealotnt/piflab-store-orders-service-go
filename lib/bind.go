package lib

import (
	"github.com/mholt/binding"

	"errors"
	"net/http"
)

func Bind(form binding.FieldMapper, r *http.Request) error {
	errs := binding.Bind(r, form)

	if errs.Len() > 0 {
		if errs[0].Classification == binding.RequiredError {
			return errors.New(errs[0].FieldNames[0] + " is required")
		}
		if errs[0].Classification == binding.TypeError {
			return errors.New(errs[0].FieldNames[0] + " is invalid")
		}
		return errs
	}

	return nil
}
