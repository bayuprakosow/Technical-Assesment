package docs

import (
	_ "embed"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed static/index.html
var indexHTML []byte

//go:embed static/guide.html
var guideHTML []byte

//go:embed static/api.html
var apiHTML []byte

//go:embed static/style.css
var styleCSS []byte

//go:embed static/api-overview.xml
var apiOverviewXML []byte

//go:embed static/postman-collection.json
var postmanCollectionJSON []byte

// Register: GET / → /docs; /docs beranda; /docs/guide panduan setup; /docs/api referensi API;
// /docs/_/style.css stylesheet; /docs/api-overview.xml ringkasan XML;
// /docs/postman-collection.json koleksi Postman v2.1 (download / import).
func Register(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs")
	})

	r.GET("/docs/_/style.css", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/css; charset=utf-8", styleCSS)
	})

	r.GET("/docs", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})
	r.GET("/docs/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHTML)
	})

	r.GET("/docs/guide", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", guideHTML)
	})
	r.GET("/docs/guide/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", guideHTML)
	})

	r.GET("/docs/api", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", apiHTML)
	})
	r.GET("/docs/api/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", apiHTML)
	})

	r.GET("/docs/api-overview.xml", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/xml; charset=utf-8", apiOverviewXML)
	})

	r.GET("/docs/postman-collection.json", func(c *gin.Context) {
		c.Header("Content-Disposition", `attachment; filename="findings-api.postman_collection.json"`)
		c.Data(http.StatusOK, "application/json; charset=utf-8", postmanCollectionJSON)
	})
}
