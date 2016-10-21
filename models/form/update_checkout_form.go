package models

import (
	"github.com/mholt/binding"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"

	"errors"
	"net/http"
)

type UpdateCheckoutForm struct {
	Status *string
	Fields string
}

func (form *UpdateCheckoutForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.Status: binding.Field{
			Form: "status",
		},
		&form.Fields: binding.Field{
			Form: "fields",
		},
	}
}

func (form *UpdateCheckoutForm) Validate() error {
	if form.Status == nil {
		return errors.New("Status is required")
	}

	if *form.Status != "processing" &&
		*form.Status != "shipping" &&
		*form.Status != "completed" &&
		*form.Status != "cancelled" {
		return errors.New("Status is invalid, only accept processing/shipping/completed/cancelled")
	}

	return nil
}

func (form *UpdateCheckoutForm) Order(app *App, order_code string) (*Order, error) {
	var order = new(Order)
	var err error

	if order, err = (OrderRepository{app.DB}).GetOrderByOrdercode(order_code); err != nil {
		if err.Error() == "record not found" {
			return order, errors.New("Order code is invalid")
		}

		// unknown err, return anyway
		return order, err
	}

	if order.Status == "cart" {
		if *form.Status != "processing" {
			return order, errors.New("Current status is cart, only accept processing for next status")
		}
	}

	if order.Status == "processing" {
		if *form.Status != "shipping" &&
			*form.Status != "cancelled" {
			return order, errors.New("Current status is processing, only accept shipping/cancelled for next status")
		}
	}

	if order.Status == "shipping" {
		if *form.Status != "completed" &&
			*form.Status != "cancelled" {
			return order, errors.New("Current status is shipping, only accept completed/cancelled for next status")
		}
	}

	if order.Status == "completed" || order.Status == "cancelled" {
		return order, errors.New("Current status is " + order.Status + " can't change status anymore")
	}

	order.Status = *form.Status

	return order, nil
}
