package main

import (
	"fmt"
	"github.com/o0khoiclub0o/piflab-store-api-go/db/seeds/factory"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/o0khoiclub0o/piflab-store-api-go/models/repository"
)

func main() {
	app := NewApp()
	defer app.Close()

	product_counts, err := ProductRepository{app.DB}.CountProduct()

	if err != nil {
		panic(err)
	}

	sJson := `{"no-image": "yes"}`
	extraParams := make(map[string]interface{})
	factory.Json2Map(sJson, extraParams)

	if product_counts == 0 {
		for i := 0; i < 10; i++ {
			if _, err := factory.CreateProduct(app.DB, extraParams); err != nil {
				panic(err)
			}
		}
	}

	fmt.Println("Seed successfully")
}
