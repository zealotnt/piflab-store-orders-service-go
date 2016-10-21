package models

import (
	"time"
)

type Product struct {
	Id       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    int     `json:"price"`
	Provider string  `json:"provider"`
	Rating   float32 `json:"rating"`
	Status   string  `json:"status"`
	Detail   string  `json:"detail"`

	ImageData          []byte    `json:"-" sql:"-"`
	ImageThumbnailData []byte    `json:"-" sql:"-"`
	ImageDetailData    []byte    `json:"-" sql:"-"`
	Image              string    `json:"-"`
	NewImage           string    `json:"-" sql:"-"`
	ImageUpdatedAt     time.Time `json:"-"`
	ImageUrl           *string   `json:"image_url" sql:"-"`
	ImageThumbnailUrl  *string   `json:"image_thumbnail_url" sql:"-"`
	ImageDetailUrl     *string   `json:"image_detail_url" sql:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductSlice []Product

type PageUrl struct {
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type ProductPage struct {
	Data   *ProductSlice `json:"data"`
	Paging PageUrl       `json:"paging"`
}

type ImageField int

const (
	IMAGE ImageField = iota
)

type ImageSize int

const (
	ORIGIN ImageSize = iota
	THUMBNAIL
	DETAIL
)
