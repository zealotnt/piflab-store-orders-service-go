package models

import (
	"github.com/mholt/binding"

	"errors"
	"net/http"
	"strings"
)

type GetCheckoutForm struct {
	Offset    uint
	Limit     uint
	Status    *string
	Sort      *string
	SortField string
	SortOrder string
	Fields    string
	Search    string
}

func (form *GetCheckoutForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Offset: binding.Field{
			Form: "offset",
		},
		&form.Limit: binding.Field{
			Form: "limit",
		},
		&form.Status: binding.Field{
			Form: "status",
		},
		&form.Sort: binding.Field{
			Form: "sort",
		},
		&form.Fields: binding.Field{
			Form: "fields",
		},
		&form.Search: binding.Field{
			Form: "q",
		},
	}
}

func (form *GetCheckoutForm) Validate() error {
	if form.Status != nil {
		if *form.Status != "processing" &&
			*form.Status != "shipping" &&
			*form.Status != "completed" &&
			*form.Status != "cancelled" {
			return errors.New("Only support processing/shipping/completed/cancelled in status field")
		}
	} else {
		var default_status string
		default_status = ""
		form.Status = &default_status
	}

	if form.Sort != nil {
		params := strings.Split(*form.Sort, "|")
		switch len(params) {
		case 1:
			if params[0] != "created_at" &&
				params[0] != "updated_at" &&
				params[0] != "id" {
				return errors.New("Only support created_at/updated_at/id in sort field")
			}
			form.SortField = params[0]
			form.SortOrder = "desc"

		case 2:
			if params[0] != "created_at" &&
				params[0] != "updated_at" &&
				params[0] != "id" {
				return errors.New("Only support created_at/updated_at/id in sort field")
			}

			if params[1] != "desc" &&
				params[1] != "asc" {
				return errors.New("Only support desc/asc in sort field")
			}
			form.SortField = params[0]
			form.SortOrder = params[1]

		default:
			return errors.New("Invalid sort field format, should be: create_at/updated_at/id|desc/asc")
		}
	} else {
		var default_sort string
		default_sort = "id|desc"
		form.Sort = &default_sort
		form.SortField = "id"
		form.SortOrder = "desc"
	}

	return nil
}
