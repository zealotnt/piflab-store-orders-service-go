package models

import (
	"github.com/mholt/binding"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"

	"errors"
	"net/http"
)

type CheckoutCartForm struct {
	AccessToken     *string     `json:"access_token"`
	Items           []OrderItem `json:"items,omitempty"`
	CustomerName    *string     `json:"name"`
	CustomerAddress *string     `json:"address"`
	CustomerPhone   *string     `json:"phone"`
	CustomerEmail   *string     `json:"email"`
	CustomerNote    *string     `json:"note"`
	Fields          string
}

func (form *CheckoutCartForm) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&form.AccessToken: binding.Field{
			Form: "access_token",
		},
		&form.CustomerName: binding.Field{
			Form: "name",
		},
		&form.CustomerAddress: binding.Field{
			Form: "address",
		},
		&form.CustomerPhone: binding.Field{
			Form: "phone",
		},
		&form.CustomerEmail: binding.Field{
			Form: "email",
		},
		&form.CustomerNote: binding.Field{
			Form: "note",
		},
		&form.Fields: binding.Field{
			Form: "fields",
		},
		&form.Fields: binding.Field{
			Form: "items",
		},
	}
}

func (form *CheckoutCartForm) Validate() error {
	if form.AccessToken == nil {
		return errors.New("Access Token is required")
	}

	if form.CustomerName == nil {
		return errors.New("Customer's Name is required")
	}

	if form.CustomerAddress == nil {
		return errors.New("Customer's Address is required")
	}

	if form.CustomerPhone == nil {
		return errors.New("Customer's Phone number is required")
	}

	if form.CustomerEmail == nil {
		return errors.New("Customer's Email is required")
	}
	if !ValidateEmail(*form.CustomerEmail) {
		return errors.New("Customer's Email address is invalid")
	}

	return nil
}

func (form *CheckoutCartForm) Order(app *App) (*Order, error) {
	var order = new(Order)
	var err error

	// TODO: Check if the cart_id exists in the db
	if order, err = (OrderRepository{app.DB}).GetOrder(*form.AccessToken); err == nil {
		// unknown err, return anyway
		return nil, errors.New("Order is already checked out, please use another order")
	}

	order.OrderInfo.CustomerName = *form.CustomerName
	order.OrderInfo.CustomerAddress = *form.CustomerAddress
	order.OrderInfo.CustomerPhone = *form.CustomerPhone
	order.OrderInfo.CustomerEmail = *form.CustomerEmail

	if form.CustomerNote != nil {
		order.OrderInfo.CustomerNote = *form.CustomerNote
	}

	// Change status to processing, any other change to oder_items is rejected from now on
	order.Status = "processing"

	return order, err
}
