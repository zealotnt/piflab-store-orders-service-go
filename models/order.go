package models

import (
	"strconv"
	"time"
)

type Amount struct {
	Subtotal uint `json:"subtotal"`
	Shipping uint `json:"shipping"`
	Total    uint `json:"total"`
}

type OrderInfo struct {
	OrderCode       string `json:"-" sql:"column:code"`
	CustomerName    string `json:"name" sql:"customer_name"`
	CustomerAddress string `json:"address" sql:"customer_address"`
	CustomerPhone   string `json:"phone" sql:"customer_phone"`
	CustomerEmail   string `json:"email" sql:"customer_email"`
	CustomerNote    string `json:"note" sql:"column:note"`
}

type CheckoutReturn struct {
	Id        string      `json:"id,omitempty"`
	Items     []OrderItem `json:"items,omitempty"`
	Amounts   Amount      `json:"amounts" sql:"-"`
	OrderInfo *OrderInfo  `json:"customer,omitempty" sql:"-"`
	Status    string      `json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CheckoutReturnSlice []CheckoutReturn
type OrderSlice []Order

type Order struct {
	Id          uint   `json:"-"`
	AccessToken string `json:"access_token,omitempty"`
	Status      string `json:"status"`

	Items []OrderItem `json:"items" sql:"order_items"`

	OrderInfo `json:"-"`

	Amounts    Amount `json:"amounts" sql:"-"`
	TotalPrice uint   `json:"-" sql:"total_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderStatusLog struct {
	Id        uint   `sql:"id"`
	OrderId   uint   `sql:"order_id, REFERENCES Orders(id)"`
	Status    string `sql:"status"`
	CreatedAt time.Time
}

type OrderItem struct {
	Id                       uint    `json:"id,string" sql:"id"`
	OrderId                  uint    `json:"-" sql:"REFERENCES Orders(id)"`
	ProductId                uint    `json:"product_id,string" sql:"column:product_id"`
	ProductName              string  `json:"name" sql:"column:name"`
	ProductImageThumbnailUrl *string `json:"image_thumbnail_url" sql:"-"`
	ProductPrice             int     `json:"price" sql:"column:price"`
	Quantity                 int     `json:"quantity"`
}

type OrderPage struct {
	Data   *CheckoutReturnSlice `json:"data"`
	Paging PageUrl              `json:"paging"`
}

func getOrderPage(host_url string, offset uint, limit uint, sort string, search string, total uint) PageUrl {
	prevNum := uint64(offset - limit)
	nextNum := uint64(offset + limit)

	if total == 0 {
		// return null next/previous field
		return PageUrl{}
	}
	if offset < limit {
		prevNum = 0
	}
	if total <= offset {
		if total > limit {
			prevNum = uint64(total - limit)
		} else {
			prevNum = 0
		}
	}

	next := host_url + "/orders?offset=" + strconv.FormatUint(nextNum, 10) + "&limit=" + strconv.FormatUint(uint64(limit), 10)
	previous := host_url + "/orders?offset=" + strconv.FormatUint(prevNum, 10) + "&limit=" + strconv.FormatUint(uint64(limit), 10)
	if sort != "" {
		next += "&sort=" + sort
		previous += "&sort=" + sort
	}
	if search != "" {
		next += "&q=" + search
		previous += "&q=" + search
	}

	// Nothing to show on next_url
	if uint64(total) <= nextNum {
		// If offset already zero, not thing to show on previous_url also
		if offset == 0 {
			return PageUrl{}
		}

		// At least, we have something to show on previous_url
		return PageUrl{
			Previous: &previous,
		}
	}
	if offset == 0 {
		return PageUrl{
			Next: &next,
		}
	}
	return PageUrl{
		Next:     &next,
		Previous: &previous,
	}
}

func (orders OrderSlice) GetPaging(host_url string, offset uint, limit uint, sort string, search string, total uint) *OrderPage {
	orders_return := CheckoutReturnSlice{}
	for idx, _ := range orders {
		orders_return = append(orders_return, orders[idx].ReturnCheckoutRequest())
		orders_return[idx].Id = orders[idx].OrderCode
	}
	return &OrderPage{
		Data:   &orders_return,
		Paging: getOrderPage(host_url, offset, limit, sort, search, total),
	}
}

func (order *Order) CalculateAmount() {
	for _, item := range order.Items {
		order.Amounts.Subtotal += uint(item.ProductPrice) * uint(item.Quantity)
	}
	order.Amounts.Shipping = 0
	order.Amounts.Total = order.Amounts.Shipping + order.Amounts.Subtotal

	order.TotalPrice = order.Amounts.Total
}

func (order *Order) EraseAccessToken() {
	order.AccessToken = ""
}

func (order *Order) RemoveZeroQuantityItems() {
	for idx, _ := range order.Items {
		if order.Items[idx].Quantity <= 0 {
			order.Items = append(order.Items[:idx], order.Items[idx+1:]...)
			return
		}
	}
}

func (order *Order) ReturnCheckoutRequest() CheckoutReturn {
	ret := new(CheckoutReturn)
	ret.Id = order.OrderCode
	ret.Items = order.Items
	ret.Amounts = order.Amounts
	if order.OrderInfo.CustomerName != "" {
		ret.OrderInfo = &order.OrderInfo
	}
	ret.Status = order.Status
	ret.UpdatedAt = order.UpdatedAt
	ret.CreatedAt = order.CreatedAt
	return *ret
}
