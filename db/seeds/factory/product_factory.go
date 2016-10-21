package factory

import (
	"github.com/icrowley/fake"
	. "github.com/mitchellh/mapstructure"
	"github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/form"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"

	"bytes"
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

func NewProduct(params ...map[string]interface{}) (*Product, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	fh, err := os.Open(os.Getenv("FULL_IMPORT_PATH") + "/db/seeds/factory/golang.png")
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	dataBytes := bytes.Buffer{}

	dataBytes.ReadFrom(fh)

	var product *Product

	if params != nil {
		if _, ok := params[0]["no-image"]; ok {
			product = &Product{
				Name:     fake.ProductName(),
				Price:    rand.Intn(100000),
				Provider: fake.Company(),
				Rating:   rand.Float32() * float32(rand.Intn(5)),
				Status:   STATUS_OPTIONS[rand.Intn(len(STATUS_OPTIONS))],
				Detail:   fake.ParagraphsN(1),
			}
			goto ignore_image
		}
	}

	product = &Product{
		Name:               fake.ProductName(),
		Price:              rand.Intn(100000),
		Provider:           fake.Company(),
		Rating:             rand.Float32() * float32(rand.Intn(5)),
		Status:             STATUS_OPTIONS[rand.Intn(len(STATUS_OPTIONS))],
		Detail:             fake.ParagraphsN(1),
		ImageData:          dataBytes.Bytes(),
		ImageThumbnailData: dataBytes.Bytes(),
		ImageDetailData:    dataBytes.Bytes(),
		Image:              "golang.png",
	}

ignore_image:
	if params != nil {
		err := Decode(params[0], product)
		return product, err
	}

	return product, nil
}

func CreateProduct(DB *lib.DB, params ...map[string]interface{}) (*Product, error) {
	product, err := NewProduct(params...)

	if err != nil {
		return product, err
	}

	return product, (ProductRepository{DB}).SaveProduct(product)
}

func Json2Map(sJson string, myMap map[string]interface{}) error {
	return json.Unmarshal([]byte(sJson), &myMap)
}
