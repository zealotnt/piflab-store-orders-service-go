package handlers_test

import (
	. "github.com/o0khoiclub0o/piflab-store-api-go/handlers"
	"github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var app *lib.App
var extraParams = map[string]string{
	"name":     "xbox",
	"price":    "70000",
	"provider": "Microsoft",
	"rating":   "3.5",
	"status":   "available",
	"detail":   "some text",
}
var initialize_product = &Product{}

var _ = BeforeSuite(func() {
	app = lib.NewApp()
	app.AddRoutes(GetRoutes())

	// Create a product with image
	path := os.Getenv("FULL_IMPORT_PATH") + "/db/seeds/factory/golang.png"
	response := MultipartRequest("POST", "/products", extraParams, "image", path)
	Expect(response.Code).To(Equal(201))
	Expect(response.Body).To(ContainSubstring(`"image_url":"https://`))
	Expect(response.Body).To(ContainSubstring(`"image_thumbnail_url":"https://`))
	Expect(response.Body).To(ContainSubstring(`"image_detail_url":"https://`))

	// Create a product without image
	response = MultipartRequest("POST", "/products", extraParams, "", "")
	Expect(response.Code).To(Equal(201))
	Expect(response.Body).To(ContainSubstring(`"image_url":null`))
	Expect(response.Body).To(ContainSubstring(`"image_thumbnail_url":null`))
	Expect(response.Body).To(ContainSubstring(`"image_detail_url":null`))

	body, _ := ioutil.ReadAll(response.Body)
	err := json.Unmarshal(body, initialize_product)
	Expect(err).To(BeNil())
})

var _ = AfterSuite(func() {
	response := Request("DELETE", "/products/"+strconv.FormatUint(uint64(initialize_product.Id), 10), "")
	Expect(response.Code).To(Equal(200))
	app.Close()
})

func Request(method string, route string, body interface{}) *httptest.ResponseRecorder {
	return app.Request(method, route, body)
}

func getProducts(body []byte) (*ProductSlice, error) {
	products_pages := ProductPage{}
	if err := json.Unmarshal(body, &products_pages); err != nil {
		return nil, err
	}

	return products_pages.Data, nil
}

func getFirstAvailableId(response *httptest.ResponseRecorder) uint {
	body, _ := ioutil.ReadAll(response.Body)
	if products, _ := getProducts(body); products != nil {
		return (*products)[0].Id
	}

	return 0
}

func getFirstAvailableUrl() string {
	response := Request("GET", "/products", "")
	return fmt.Sprintf("/products/%d", getFirstAvailableId(response))
}

func getFirstImagelessProductUrl() string {
	response := Request("GET", "/products?offset=0&limit=100", "")
	body, _ := ioutil.ReadAll(response.Body)

	if products, _ := getProducts(body); products != nil {
		for idx, product := range *products {
			if product.ImageUrl == nil {
				return fmt.Sprintf("/products/%d", (*products)[idx].Id)
			}
		}
	}

	return "/products/0"
}

func MultipartRequest(method string, route string, params map[string]string, paramName, path string) *httptest.ResponseRecorder {
	body := lib.BodyMultipart{}

	writer := multipart.NewWriter(&body.Buff)

	file, err := os.Open(path)
	if err == nil {
		part, _ := writer.CreateFormFile(paramName, filepath.Base(path))
		io.Copy(part, file)
	}
	defer file.Close()

	for key, val := range params {
		writer.WriteField(key, val)
	}
	writer.Close()

	body.ContentType = writer.FormDataContentType()

	return Request(method, route, body)
}
