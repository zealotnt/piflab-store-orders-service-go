package repository_test

import (
	"github.com/o0khoiclub0o/piflab-store-api-go/lib"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}

var app *lib.App

var _ = BeforeSuite(func() {
	app = lib.NewApp()
})

var _ = AfterSuite(func() {
	app.Close()
})
