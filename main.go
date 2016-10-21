package main

import (
	"github.com/o0khoiclub0o/piflab-store-api-go/handlers"
	. "github.com/o0khoiclub0o/piflab-store-api-go/lib"
)

func main() {
	app := NewApp()
	app.AddRoutes(handlers.GetRoutes())
	app.Run()
	defer app.Close()
}
