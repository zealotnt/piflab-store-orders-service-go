package handlers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IndexHandler", func() {
	It("returns 200", func() {
		response := Request("GET", "/", "")

		Expect(response.Code).To(Equal(200))
		Expect(response.Body).To(ContainSubstring(`"version":"1.0.0"`))
	})
})
