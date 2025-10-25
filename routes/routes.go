package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/shemigam1/hngxiii-currency-exchange/services"
)

// POST /countries/refresh → Fetch all countries and exchange rates, then cache them in the database
// GET /countries → Get all countries from the DB (support filters and sorting) - ?region=Africa | ?currency=NGN | ?sort=gdp_desc
// GET /countries/:name → Get one country by name
// DELETE /countries/:name → Delete a country record
// GET /status → Show total countries and last refresh timestamp
// GET /countries/image → serve summary image
func Routes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to HNGXIII Currency Exchange API",
		})
	})

	countriesGroup := r.Group("/countries")
	{
		countriesGroup.POST("/refresh", services.RefreshCountries)
		countriesGroup.GET("/", services.GetAllCountries)
		countriesGroup.GET("/:name", services.GetCountry)
		countriesGroup.DELETE("/:name", services.DeleteCountry)
	}

	r.GET("/status", services.GetStatus)
	r.GET("/countries/image", services.GetSummaryImage)
}
