package handlers

import (
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/form"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"

	"net/http"
)

func CheckoutCartHandler(app *App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, c Context) {
		form := new(CheckoutCartForm)

		if err := Bind(form, r); err != nil {
			JSON(w, err, 400)
			return
		}

		if err := form.Validate(); err != nil {
			JSON(w, err, 422)
			return
		}

		order, err := form.Order(app)
		if err != nil {
			JSON(w, err, 424)
			return
		}

		if err := (OrderRepository{app.DB}).SaveOrder(order); err != nil {
			JSON(w, err, 500)
			return
		}

		ret := order.ReturnCheckoutRequest()
		JSON(w, ret)
	}
}

func GetCheckoutHandler(app *App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, c Context) {
		host_url := HostURL(r)

		form := new(GetCheckoutForm)

		if err := Bind(form, r); err != nil {
			form.Offset = 0
			form.Limit = 10
		}

		if err := form.Validate(); err != nil {
			JSON(w, err, 422)
			return
		}

		orders, total, err := OrderRepository{app.DB}.GetPage(form.Offset, form.Limit, *form.Status, form.SortField, form.SortOrder, form.Search)
		if err != nil {
			JSON(w, err, 500)
			return
		}

		// Remove items from the checkout GET
		for idx, _ := range *orders {
			(*orders)[idx].Items = nil
		}

		orders_by_pages := orders.GetPaging(host_url, form.Offset, form.Limit, *form.Sort, form.Search, total)
		JSON(w, orders_by_pages)
	}
}

func GetCheckoutDetailHandler(app *App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, c Context) {
		// the purpose of form is to get form.Fields, so don't care about Binding errors
		form := new(GetCheckoutForm)
		Bind(form, r)

		order, err := (OrderRepository{app.DB}).FindByOrderCode(c.Params["id"])
		if err != nil {
			JSON(w, err, 404)
			return
		}

		order.CalculateAmount()
		ret := order.ReturnCheckoutRequest()
		JSON(w, ret)
	}
}

func UpdateCheckoutStatusHandler(app *App) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, c Context) {
		form := new(UpdateCheckoutForm)
		if err := Bind(form, r); err != nil {
			JSON(w, err, 400)
			return
		}

		if err := form.Validate(); err != nil {
			JSON(w, err, 422)
			return
		}

		order, err := form.Order(app, c.Params["id"])
		if err != nil {
			JSON(w, err, 424)
			return
		}

		if err := (OrderRepository{app.DB}).SaveOrder(order); err != nil {
			JSON(w, err, 500)
			return
		}

		order.CalculateAmount()
		ret := order.ReturnCheckoutRequest()
		JSON(w, ret)
	}
}
