package models_test

import (
	. "github.com/o0khoiclub0o/piflab-store-api-go/models"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"io/ioutil"
	"strconv"
)

type ProductSliceTest struct {
	description string
	url         string
	expect      string
}

var _ = Describe("Product Test", func() {
	var first bool
	var product *Product

	BeforeEach(func() {
		if first == false {
			first = true

			response := Request("GET", "/products", "")
			id := getFirstAvailableId(response)
			product, _ = (ProductRepository{app.DB}).FindById(id)
		}
	})

	var _ = Describe("GetImageUrl Test", func() {
		It("returns valid image url", func() {
			product_that_has_image := getProductThatHasImage()
			product_that_has_image.GetImageUrl()
			Expect(*product_that_has_image.ImageUrl).NotTo(BeNil())
			Expect(*product_that_has_image.ImageThumbnailUrl).NotTo(BeNil())
			Expect(*product_that_has_image.ImageDetailUrl).NotTo(BeNil())
		})
	})

	var _ = Describe("GetImageUrlType Test", func() {
		It("returns valid image url, base on ORIGIN ImageType", func() {
			_, err := product.GetImageUrlType(IMAGE, ORIGIN)
			Expect(err).To(BeNil())
		})

		It("returns error when trying to get image url, because the param is out of scope of enum ImageType", func() {
			url, err := product.GetImageUrlType(IMAGE, 99)
			Expect(url).To(Equal(""))
			Expect(err.Error()).To(ContainSubstring("field too short, minimum length 1: Key"))
		})
	})

	var _ = Describe("GetImagePath and GetImageContentType Test", func() {
		It("returns image url with image extension", func() {
			// rename Image file name to contain png extension
			product.Image = "dummyFile.png"
			type TestStruct_ImageSize_Expected_Extension struct {
				size        ImageSize
				expectedExt string
				extension   string
				contentType string
			}

			test_cases := []TestStruct_ImageSize_Expected_Extension{{
				// Origin will follow origin file's extension
				ORIGIN,
				"/image_origin_",
				".png",
				"image/png"}, {

				// thumbnail always have png extension
				THUMBNAIL,
				"/image_thumbnail_",
				".png",
				"image/png"}, {

				// detail always have png extension
				DETAIL,
				"/image_detail_",
				".png",
				"image/png"},
			}
			for _, test := range test_cases {
				path := product.GetImagePath(IMAGE, test.size)
				expected := "products/" +
					strconv.FormatUint(uint64(product.Id), 10) +
					test.expectedExt
				Expect(path).To(ContainSubstring(expected))
				Expect(path).To(ContainSubstring(test.extension))

				retContentType := product.GetImageContentType(IMAGE, test.size)
				Expect(retContentType).To(Equal(test.contentType))
			}

			// try again with jpg extension
			product.Image = "dummyFile.jpg"
			path := product.GetImagePath(IMAGE, ORIGIN)
			expected := "products/" +
				strconv.FormatUint(uint64(product.Id), 10) +
				"/image_origin_"
			Expect(path).To(ContainSubstring(expected))
			Expect(path).To(ContainSubstring(".jpg"))

			retContentType := product.GetImageContentType(IMAGE, ORIGIN)
			Expect(retContentType).To(Equal("image/jpg"))
		})

		It("returns image url without image extension", func() {
			var fields = []ImageField{IMAGE}
			for _, field := range fields {
				// rename Image file name, so we don't use regex's result to give to file name
				product.Image = "dummyFile"
				path := product.GetImagePath(field, ORIGIN)
				expected_image := "products/" +
					strconv.FormatUint(uint64(product.Id), 10) +
					"/image_origin_"
				if field == IMAGE {
					Expect(path).To(ContainSubstring(expected_image))
				}
				Expect(path).NotTo(ContainSubstring(".png"))

				retContentType := product.GetImageContentType(field, ORIGIN)
				Expect(retContentType).To(Equal("image"))
			}
		})

		It("returns empty string, because the param is out of scope of enum ImageType", func() {
			var fields = []ImageField{IMAGE}
			for _, field := range fields {
				url := product.GetImagePath(field, 99)
				Expect(url).To(Equal(""))

				retContentType := product.GetImageContentType(field, 99)
				Expect(retContentType).To(Equal(""))
			}

			// finally, it will return empty string, because the param is out of scope of enum ImageField
			url := product.GetImagePath(99, ORIGIN)
			Expect(url).To(Equal(""))

			retContentType := product.GetImageContentType(99, ORIGIN)
			Expect(retContentType).To(Equal(""))
		})
	})
})

var _ = Describe("ProductSlice Test", func() {
	It("test all cases of GetPaging-getPage", func() {
		product_counts, err := (ProductRepository{app.DB}).CountProduct()
		Expect(err).To(BeNil())

		test_cases := []ProductSliceTest{{
			/*1*/
			description: `Get only 1 product per page, with offset=0.` +
				`Return val should contain valid "next" field, and null "previous" field`,
			url:    `/products?offset=0&limit=1`,
			expect: `"paging":{"next":"/products?offset=1\u0026limit=1","previous":null}`}, {

			/*2*/
			description: `Get page with very big offset (bigger than maximum current product).
				The result should return null "next" field, and valid "previous" field, that can:
				+ Return all of the products, if "limit" > "product_counts"
				+ Return "limit" number of products, if "limit" < "products_counts"
				--> Return "limit" number of products, if "limit" < "products_counts"`,
			url: `/products?offset=` + strconv.FormatUint(uint64(product_counts), 10) +
				`&limit=1`,
			expect: `"paging":{"next":null,"previous":"/products?offset=` +
				strconv.FormatUint(uint64(product_counts-1), 10) +
				`\u0026limit=1"}`}, {

			/*3*/
			description: `Get page with very big offset (bigger than maximum current product).
				The result should return null "next" field, and valid "previous" field, that can:
				+ Return all of the products, if "limit" > "product_counts"
				+ Return "limit" number of products, if "limit" < "products_counts"
				--> Return all of the products, if "limit" > "product_counts"`,
			url: `/products?offset=` +
				strconv.FormatUint(uint64(product_counts), 10) +
				`&limit=` +
				strconv.FormatUint(uint64(product_counts), 10),
			expect: `"paging":{"next":null,"previous":"/products?offset=0\u0026limit=` +
				strconv.FormatUint(uint64(product_counts), 10)}, {

			/*4*/
			description: `Get page with offset in the "middle" position of products
				with the "limit" value that doesn't exceed the maximum "product_counts".
				So the result will contain both "next" field, and "previous" field`,
			url: `/products?offset=` + strconv.FormatUint(uint64(product_counts/2), 10) +
				`&limit=1`,
			expect: `"paging":{"next":"/products?offset=` +
				strconv.FormatUint(uint64(product_counts/2+1), 10) + `\u0026limit=1",` +
				`"previous":"/products?offset=` +
				strconv.FormatUint(uint64(product_counts/2-1), 10) + `\u0026limit=1`},
		}

		for _, test := range test_cases {
			By(test.description)
			response := Request("GET", test.url, "")
			body, _ := ioutil.ReadAll(response.Body)
			Expect(response.Code).To(Equal(200))
			Expect(body).To(ContainSubstring(test.expect))
		}
	})
})
