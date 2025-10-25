package services

import (
	"github.com/gin-gonic/gin"
)

// POST /countries/refresh → Fetch all countries and exchange rates, then cache them in the database
// GET /countries → Get all countries from the DB (support filters and sorting) - ?region=Africa | ?currency=NGN | ?sort=gdp_desc
// GET /countries/:name → Get one country by name
// DELETE /countries/:name → Delete a country record
// GET /status → Show total countries and last refresh timestamp
// GET /countries/image → serve summary image

func RefreshCountries(c *gin.Context) {
	c.JSON(200, "refresh working")
}

func GetAllCountries(c *gin.Context) {
	c.JSON(200, "get all contries working")
}

func GetCountry(c *gin.Context) {
	c.JSON(200, "get country working")
}

func DeleteCountry(c *gin.Context) {
	c.JSON(200, "delete country working")
}

func GetStatus(c *gin.Context) {
	c.JSON(200, "get status working")
}

func GetSummaryImage(c *gin.Context) {
	c.JSON(200, "get summary image working")
}
