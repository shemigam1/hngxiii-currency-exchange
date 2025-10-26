package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shemigam1/hngxiii-currency-exchange/db"
	"github.com/shemigam1/hngxiii-currency-exchange/models"
	"gorm.io/gorm"
)

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
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "External data source unavailable",
			"details": "Could not fetch data from restcountries API",
		})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "External data source unavailable",
			"details": "Could not fetch data from restcountries API",
		})
		return
	}

	var countries []models.CountryDataResponse
	err = json.Unmarshal(body, &countries)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "External data source unavailable",
			"details": "Could not fetch data from restcountries API",
		})
		return
	}

	rates, err := GetExchangeRates()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error":   "External data source unavailable",
			"details": "Could not fetch data from exchange rates API",
		})
		return
	}

	for _, country := range countries[:10] {
		fmt.Println(country.Name)
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

	c.JSON(http.StatusOK, gin.H{"message": "Countries refreshed successfully"})
}

func GetAllCountries(c *gin.Context) {
	var countries []models.CountryInfo
	query := db.DB

	if region := c.Query("region"); region != "" {
		query = query.Where("region = ?", region)
	}

	if currency := c.Query("currency"); currency != "" {
		query = query.Where("currency_code = ?", currency)
	}

	if sort := c.Query("sort"); sort != "" {
		switch sort {
		case "gdp_desc":
			query = query.Order("estimated_gdp DESC")
		case "gdp_asc":
			query = query.Order("estimated_gdp ASC")
		case "population_desc":
			query = query.Order("population DESC")
		case "population_asc":
			query = query.Order("population ASC")
		case "name_asc":
			query = query.Order("name ASC")
		case "name_desc":
			query = query.Order("name DESC")
		}
	}

	result := query.Find(&countries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if len(countries) == 0 {
		c.JSON(http.StatusOK, []models.CountryResponse{})
		return
	}

	var response []models.CountryResponse
	for _, country := range countries {
		response = append(response, models.CountryResponse{
			ID:           fmt.Sprintf("%d", country.ID),
			Name:         country.Name,
			Capital:      country.Capital,
			Region:       country.Region,
			Population:   country.Population,
			CurrencyCode: country.CurrencyCode,
			ExchangeRate: country.ExchangeRate,
			EstimatedGdp: country.EstimatedGdp,
			FlagUrl:      country.FlagUrl,
			UpdatedAt:    country.UpdatedAt.Format(time.RFC3339),
		})
	}

	c.JSON(http.StatusOK, response)
}

func GetCountry(c *gin.Context) {
	name := c.Param("name")
	fmt.Println(name)

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": gin.H{"name": "is required"},
		})
		return
	}

	var country models.CountryInfo
	result := db.DB.Where("name = ?", name).First(&country)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	response := models.CountryResponse{
		ID:           fmt.Sprintf("%d", country.ID),
		Name:         country.Name,
		Capital:      country.Capital,
		Region:       country.Region,
		Population:   country.Population,
		CurrencyCode: country.CurrencyCode,
		ExchangeRate: country.ExchangeRate,
		EstimatedGdp: country.EstimatedGdp,
		FlagUrl:      country.FlagUrl,
		UpdatedAt:    country.UpdatedAt.Format(time.RFC3339),
	}

	c.JSON(http.StatusOK, response)
}

func DeleteCountry(c *gin.Context) {
	name := c.Param("name")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": gin.H{"name": "is required"},
		})
		return
	}

	result := db.DB.Where("name = ?", name).Delete(&models.CountryInfo{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Country not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
}

func GetStatus(c *gin.Context) {
	var total int64
	var lastCountry models.CountryInfo

	result := db.DB.Model(&models.CountryInfo{}).Count(&total)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	result = db.DB.Order("updated_at DESC").First(&lastCountry)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var lastRefreshedAt string
	if result.Error == gorm.ErrRecordNotFound {
		lastRefreshedAt = ""
	} else {
		lastRefreshedAt = lastCountry.UpdatedAt.Format(time.RFC3339)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_countries":   total,
		"last_refreshed_at": lastRefreshedAt,
	})
}

func GetSummaryImage(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{"error": "Summary image not found"})
}
