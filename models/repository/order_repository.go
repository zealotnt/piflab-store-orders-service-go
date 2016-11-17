package repository

import (
	"github.com/icrowley/fake"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"

	"errors"
	"math/rand"
	"strings"
	"time"
)

type OrderRepository struct {
	*DB
}

func (repo OrderRepository) getOrderItemsImageUrl(order_items []OrderItem) {
	for idx, item := range order_items {
		product := &Product{}

		product.Id = item.ProductId
		// don't care if the product still present
		// if it is remove, the image url will be blank
		product, _ = (ProductRepository{}).FindById(product.Id)
		order_items[idx].ProductImageThumbnailUrl = product.ImageThumbnailUrl
	}
}

func (repo OrderRepository) generateOrderCode(order *Order) error {
	rand.Seed(time.Now().UTC().UnixNano())

try_gen_other_value:
	order.OrderInfo.OrderCode = fake.CharactersN(32)

	temp_order := &Order{}
	if err := repo.DB.Where("code = ?", order.OrderInfo.OrderCode).Find(temp_order).Error; err != nil {
		// Check if err is not found -> code is unique
		if err.Error() == "record not found" {
			return nil
		}

		// Otherwise, this is database operation error
		return errors.New("Database error")
	}

	// duplicate, try again
	goto try_gen_other_value
}

func (repo OrderRepository) createOrder(order *Order) error {
	if err := repo.generateOrderCode(order); err != nil {
		return err
	}

	if err := repo.DB.Create(order).Error; err != nil {
		return err
	}

	// Create the order_status_log item
	order_status_log := OrderStatusLog{
		Code:   order.OrderInfo.OrderCode,
		Status: order.Status,
	}
	if err := repo.DB.Create(&order_status_log).Error; err != nil {
		return err
	}

	repo.getOrderItemsImageUrl(order.Items)

	return nil
}

func (repo OrderRepository) updateOrder(order *Order) error {
	tx := repo.DB.Begin()

	// Update the order
	if err := tx.Save(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// TODO: Each time there is a change in checkout, save it to db
	order_status_log := OrderStatusLog{
		Code:   order.OrderInfo.OrderCode,
		Status: order.Status,
	}

	if err := tx.Create(&order_status_log).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	repo.getOrderItemsImageUrl(order.Items)

	// Don't return access_token when updating
	order.EraseAccessToken()

	return nil
}

func (repo OrderRepository) FindByOrderCode(order_code string) (*Order, error) {
	order := &Order{}
	items := &[]OrderItem{}

	// find a order by its access_token
	if err := repo.DB.Where("code = ?", order_code).Find(order).Error; err != nil {
		return nil, err
	}

	// use order.Id to find its OrderItem data (order.Id is its forein key)
	if err := repo.DB.Where("order_id = ?", order.Id).Find(items).Error; err != nil {
		return nil, err
	}

	// use the order.Items to update products information
	order.Items = *items

	repo.getOrderItemsImageUrl(order.Items)

	return order, nil
}

func (repo OrderRepository) GetOrder(access_token string) (*Order, error) {
	order := &Order{}
	items := &[]OrderItem{}

	// find a order by its access_token
	if err := repo.DB.Where("access_token = ?", access_token).Find(order).Error; err != nil {
		return nil, err
	}

	// use order.Id to find its OrderItem data (order.Id is its forein key)
	if err := repo.DB.Where("order_id = ?", order.Id).Find(items).Error; err != nil {
		return nil, err
	}

	// use the order.Items to update products information
	order.Items = *items

	repo.getOrderItemsImageUrl(order.Items)

	return order, nil
}

func (repo OrderRepository) SaveOrder(order *Order) error {
	if order.OrderInfo.OrderCode == "" {
		return repo.createOrder(order)
	}
	return repo.updateOrder(order)
}

func (repo OrderRepository) GetPage(offset uint, limit uint, status string, sort_field string, sort_order string, search string) (*OrderSlice, uint, error) {
	orders := &OrderSlice{}
	items := &[]OrderItem{}
	var count uint
	var err error
	var where_param string

	if status != "" {
		where_param = "status='" + status + "'"
	}

	if search != "" {
		if where_param != "" {
			where_param += " AND"
		}
		where_param += " LOWER(customer_name) LIKE  '%" + strings.ToLower(search) + "%'"
	}

	err = repo.DB.Order(sort_field + " " + sort_order).Offset(int(offset)).Where(where_param).Limit(int(limit)).Find(orders).Error

	for idx, order := range *orders {
		// use order.Id to find its OrderItem data (order.Id is its forein key)
		if err := repo.DB.Where("order_id = ?", order.Id).Find(items).Error; err != nil {
			return nil, 0, err
		}
		// use the order.Items to update products information
		(*orders)[idx].Items = *items
		(*orders)[idx].CalculateAmount()
	}

	// TODO: count number of orders and return
	repo.DB.Table("orders").Count(&count)

	return orders, count, err
}
