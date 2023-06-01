package echo_server_test

import (
	"api-gw/pkg/server/echo_server"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Swagger handler tests", func() {
	It("should ensure default for swagger ui opts", func() {
		opts := echo_server.SwaggerUIOpts{}
		opts.EnsureDefaults()

		Expect(opts.BasePath).To(Equal("/"))
		Expect(opts.Path).To(Equal("docs"))
		Expect(opts.SpecURL).To(Equal("/swagger.json"))
		Expect(opts.Title).To(Equal("API documentation"))
	})

	It("should test swagger handler", func() {
		e := echo.New()
		e.GET("/:datamodel/docs", echo_server.SwaggerUI)

		req, _ := http.NewRequest("GET", "/docs", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/test/docs")
		c.SetParamNames("datamodel")
		c.SetParamValues("test")

		err := echo_server.SwaggerUI(c)
		Expect(err).To(BeNil())
		Expect(rec.Code).To(Equal(200))
		Expect(rec.Header().Get("Content-Type")).To(Equal("text/html; charset=UTF-8"))
	})
})
