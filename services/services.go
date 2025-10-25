package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shemigam1/hngxiii-currency-exchange/db"
	"github.com/shemigam1/hngxiii-currency-exchange/models"
)

// POST /countries/refresh → Fetch all countries and exchange rates, then cache them in the database
// GET /countries → Get all countries from the DB (support filters and sorting) - ?region=Africa | ?currency=NGN | ?sort=gdp_desc
// GET /countries/:name → Get one country by name
// DELETE /countries/:name → Delete a country record
// GET /status → Show total countries and last refresh timestamp
// GET /countries/image → serve summary image

type ExchangeRateResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func GetExchangeRates() (map[string]float64, error) {
	exchange_rates_url := "https://open.er-api.com/v6/latest/USD"
	resp, err := http.Get(exchange_rates_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rates ExchangeRateResponse
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, err
	}

	return rates.Rates, nil
}

func RefreshCountries(c *gin.Context) {
	countries_url := "https://restcountries.com/v2/all?fields=name,capital,region,population,flag,currencies"
	resp, err := http.Get(countries_url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch countries"})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Unmarshal into a slice of structs
	var countries []models.CountryDataResponse
	// fmt.Println(body[:4])
	err = json.Unmarshal(body, &countries)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	// Now you can loop through and process
	rates, err := GetExchangeRates()
	fmt.Println(rates)
	for _, country := range countries[:10] {
		fmt.Println(country.Name)
		// Save to database, transform data, etc.
		var countryInfo models.CountryInfo
		countryInfo.Name = country.Name
		countryInfo.Capital = country.Capital
		countryInfo.Region = country.Region
		countryInfo.Population = country.Population
		countryInfo.FlagUrl = country.Flag
		if country.Currencies == nil || len(country.Currencies) <= 0 {
			countryInfo.CurrencyCode = ""
			countryInfo.ExchangeRate = 0
			countryInfo.EstimatedGdp = 0

		} else {
			countryInfo.CurrencyCode = country.Currencies[0].Code
		}
		if _, exists := rates[countryInfo.CurrencyCode]; exists {
			countryInfo.ExchangeRate = rates[countryInfo.CurrencyCode]
			countryInfo.EstimatedGdp = float64(countryInfo.Population) * float64(1000+rand.Intn(1001)) / countryInfo.ExchangeRate
		} else {
			countryInfo.ExchangeRate = 0
			countryInfo.EstimatedGdp = 0
		}
		result := db.DB.Where("name = ?", countryInfo.Name).
			Assign(countryInfo).
			FirstOrCreate(&countryInfo)

		if result.Error != nil {
			fmt.Println("Error saving country:", result.Error)
			continue
		}
	}

	c.JSON(http.StatusOK, countries)
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
