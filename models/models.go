package models

import "gorm.io/gorm"

//  id — auto-generated
//  name — required
//  capital — optional
//  region — optional
//  population — required
//  currency_code — required
//  exchange_rate — required
//  estimated_gdp — computed from population × random(1000–2000) ÷ exchange_rate
//  flag_url — optional
//  last_refreshed_at — auto timestamp

type CountryInfo struct {
	gorm.Model
	Name         string  `gorm:"column:name;not null;uniqueIndex"`
	Capital      string  `gorm:"column:capital"`
	Region       string  `gorm:"column:region"`
	Population   int64   `gorm:"column:population;not null"`
	CurrencyCode string  `gorm:"column:currency_code"`
	ExchangeRate float64 `gorm:"column:exchange_rate"`
	EstimatedGdp float64 `gorm:"column:estimated_gdp"`
	FlagUrl      string  `gorm:"column:flag_url"`
}

type StringAnalysisResult struct {
	Value                 string
	Length                int64
	IsPalindrome          bool
	UniqueCharacters      int64
	WordCount             int64
	SHA256Hash            string
	CharacterFrequencyMap map[string]int64
}

type CountryResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Capital      string  `json:"capital"`
	Region       string  `json:"region"`
	Population   int64   `json:"population"`
	CurrencyCode string  `json:"currency_code"`
	ExchangeRate float64 `json:"exchange_rate"`
	EstimatedGdp float64 `json:"estimated_gdp"`
	FlagUrl      string  `json:"flag_url"`
	UpdatedAt    string  `json:"last_refreshed_at"`
}

type FlagsResponse struct {
	FlagUrl string `json:"flag_url"`
}

type Status struct {
	TotalCountries int64  `json:"total_countries"`
	UpdatedAt      string `json:"last_refreshed_at"`
}

type CacheImage struct {
	TotalCountries                 int64             `json:"total_countries"`
	TopFiveCountriesByEstimatedGdp []CountryResponse `json:"top_five_countries_by_estimated_gdp"`
	UpdatedAt                      string            `json:"last_refreshed_at"`
}

type CountryDataResponse struct {
	Name       string     `json:"name"`
	Capital    string     `json:"capital"`
	Region     string     `json:"region"`
	Population int64      `json:"population"`
	Currencies []Currency `json:"currencies"`
	Flag       string     `json:"flag"`
	// Independent bool       `json:"independent"`
}

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
