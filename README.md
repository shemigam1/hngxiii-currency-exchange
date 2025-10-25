# Currency Exchange API

## Overview
This project is a Go Gin backend API designed to manage and provide comprehensive country information, including dynamic currency exchange rates, population data, and estimated GDP. It leverages PostgreSQL for data persistence and integrates with external APIs to fetch and cache up-to-date global country details and financial data.

## Features
-   **Go Gin**: High-performance and efficient web framework for building robust API endpoints.
-   **PostgreSQL**: Reliable relational database for structured storage and retrieval of country data.
-   **GORM**: An elegant ORM for Go, simplifying database interactions and object-relational mapping.
-   **External API Integration**: Seamlessly fetches country data from `restcountries.com` and real-time exchange rates from `open.er-api.com`.
-   **Data Caching**: Caches fetched country and exchange rate data in the database, reducing external API calls and improving response times.
-   **Dynamic GDP Calculation**: Computes estimated GDP for each country based on its population, a random economic factor, and the latest exchange rate.
-   **Comprehensive CRUD Operations**: Provides API endpoints for creating, retrieving, updating (via refresh), and deleting country records.
-   **Advanced Filtering and Sorting**: Supports querying countries by specific regions, currency codes, and sorting results by estimated GDP, population, or name in ascending or descending order.
-   **API Status Monitoring**: An endpoint to quickly check the total number of cached countries and the timestamp of the last data refresh.
-   **Environment Configuration**: Securely manages sensitive configuration parameters using `.env` files.

## Getting Started
### Installation
Installation steps are skipped as per user request.

### Environment Variables
The following environment variables are required to run the application:

-   `PORT`: The port number on which the Gin server will listen.
    Example: `PORT=8080`
-   `DB_HOST`: The hostname or IP address of the PostgreSQL database server.
    Example: `DB_HOST=localhost`
-   `DB_PORT`: The port number of the PostgreSQL database server.
    Example: `DB_PORT=5432`
-   `DB_USER`: The username for connecting to the PostgreSQL database.
    Example: `DB_USER=postgres`
-   `DB_PASSWORD`: The password for the PostgreSQL database user.
    Example: `DB_PASSWORD=your_password`
-   `DB_NAME`: The name of the database to connect to.
    Example: `DB_NAME=currency_exchange_db`
-   `SSL_MODE`: The SSL mode for the database connection. Typically `disable` for local development or `require` for production.
    Example: `SSL_MODE=disable`

## Usage
The API provides various endpoints to interact with country and currency exchange data. Below are practical examples for interacting with the API using `curl`.

### Refresh Country Data
Refreshes the country data and exchange rates from external sources, caching them in the database.
```bash
curl -X POST http://localhost:8080/countries/refresh
```

### Get All Countries
Retrieves a list of all cached countries. Supports filtering by `region` and `currency`, and sorting by `gdp_desc`, `gdp_asc`, `population_desc`, `population_asc`, `name_asc`, or `name_desc`.

**Example: Get all countries**
```bash
curl http://localhost:8080/countries
```

**Example: Get countries in Africa, sorted by GDP (descending)**
```bash
curl "http://localhost:8080/countries?region=Africa&sort=gdp_desc"
```

**Example: Get countries using NGN currency**
```bash
curl "http://localhost:8080/countries?currency=NGN"
```

### Get a Specific Country by Name
Retrieves details for a single country by its name. The country name must be provided as a query parameter.
```bash
curl "http://localhost:8080/countries?name=Nigeria"
```

### Delete a Country by Name
Deletes a country record from the database by its name. The country name must be provided as a query parameter.
```bash
curl -X DELETE "http://localhost:8080/countries?name=Nigeria"
```

### Get API Status
Retrieves the total number of countries currently cached in the database and the timestamp of the last data refresh.
```bash
curl http://localhost:8080/status
```

### Get Summary Image
This endpoint is intended to serve a summary image, but currently returns a "Not Found" error.
```bash
curl http://localhost:8080/countries/image
```

## API Documentation
### Base URL
`http://localhost:[PORT]`

### Endpoints

#### POST /countries/refresh
Refreshes all country data and exchange rates from external APIs and updates the database.

**Request**:
No payload required.

**Response**:
```json
{
  "message": "Countries refreshed successfully"
}
```

**Errors**:
-   `503 Service Unavailable`: External data source unavailable (e.g., restcountries.com or open.er-api.com could not be reached or returned an error).

#### GET /countries
Retrieves a list of all country records from the database. Supports optional filtering and sorting via query parameters.

**Request**:
Query Parameters:
-   `region` (string, optional): Filter by country region (e.g., `Africa`, `Europe`).
-   `currency` (string, optional): Filter by currency code (e.g., `NGN`, `USD`).
-   `sort` (string, optional): Sort order for results. Accepted values:
    -   `gdp_desc`: Sort by estimated GDP in descending order.
    -   `gdp_asc`: Sort by estimated GDP in ascending order.
    -   `population_desc`: Sort by population in descending order.
    -   `population_asc`: Sort by population in ascending order.
    -   `name_asc`: Sort by name in ascending order.
    -   `name_desc`: Sort by name in descending order.

**Response**:
```json
[
  {
    "id": "1",
    "name": "Nigeria",
    "capital": "Abuja",
    "region": "Africa",
    "population": 206139589,
    "currency_code": "NGN",
    "exchange_rate": 780.0,
    "estimated_gdp": 600000000000,
    "flag_url": "https://restcountries.com/data/nga.svg",
    "last_refreshed_at": "2024-07-30T10:00:00Z"
  },
  {
    "id": "2",
    "name": "Ghana",
    "capital": "Accra",
    "region": "Africa",
    "population": 31072940,
    "currency_code": "GHS",
    "exchange_rate": 14.5,
    "estimated_gdp": 70000000000,
    "flag_url": "https://restcountries.com/data/gha.svg",
    "last_refreshed_at": "2024-07-30T10:00:00Z"
  }
]
```

**Errors**:
-   `500 Internal Server Error`: An unexpected error occurred on the server while retrieving countries.

#### GET /countries?name={name}
Retrieves a single country record by its exact name. The name must be provided as a query parameter.

**Request**:
Query Parameters:
-   `name` (string, required): The exact name of the country to retrieve (e.g., `Nigeria`).

**Response**:
```json
{
  "id": "1",
  "name": "Nigeria",
  "capital": "Abuja",
  "region": "Africa",
  "population": 206139589,
  "currency_code": "NGN",
  "exchange_rate": 780.0,
  "estimated_gdp": 600000000000,
  "flag_url": "https://restcountries.com/data/nga.svg",
  "last_refreshed_at": "2024-07-30T10:00:00Z"
}
```

**Errors**:
-   `400 Bad Request`: The `name` query parameter is missing from the request.
-   `404 Not Found`: No country found with the provided name.
-   `500 Internal Server Error`: An unexpected error occurred on the server while retrieving the country.

#### DELETE /countries?name={name}
Deletes a single country record from the database by its exact name. The name must be provided as a query parameter.

**Request**:
Query Parameters:
-   `name` (string, required): The exact name of the country to delete (e.g., `Nigeria`).

**Response**:
```json
{
  "message": "Country deleted successfully"
}
```

**Errors**:
-   `400 Bad Request`: The `name` query parameter is missing from the request.
-   `404 Not Found`: No country found with the provided name to delete.
-   `500 Internal Server Error`: An unexpected error occurred on the server while deleting the country.

#### GET /status
Retrieves the current status of the cached country data, including the total number of countries and the timestamp of the last refresh.

**Request**:
No payload required.

**Response**:
```json
{
  "total_countries": 10,
  "last_refreshed_at": "2024-07-30T10:00:00Z"
}
```
If no countries are present in the database, `last_refreshed_at` will be an empty string.

**Errors**:
-   `500 Internal Server Error`: An unexpected error occurred on the server while retrieving status information.

#### GET /countries/image
This endpoint is intended to generate or serve a summary image of country data.

**Request**:
No payload required.

**Response**:
Currently, this endpoint returns a 404 error and does not serve an image.

**Errors**:
-   `404 Not Found`: The summary image resource is not found or not yet implemented.

## Technologies Used

| Technology    | Description                                   |
| :------------ | :-------------------------------------------- |
| Go            | Primary programming language                  |
| Gin           | High-performance HTTP web framework           |
| GORM          | Object-Relational Mapping (ORM) library       |
| PostgreSQL    | Robust relational database system             |
| godotenv      | Loads environment variables from `.env` file  |
| air           | Live-reloading for Go applications during development |

## 👤 Author

**Your Name**

-   **LinkedIn**: [Your LinkedIn Profile]
-   **Twitter**: [@YourTwitterHandle]
-   **Portfolio**: [Your Portfolio Website]

## Badges
![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Gin Framework](https://img.shields.io/badge/Gin-v1.11.0-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-v1.31.0-blue?style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg?style=for-the-badge)

[![Readme was generated by Dokugen](https://img.shields.io/badge/Readme%20was%20generated%20by-Dokugen-brightgreen)](https://www.npmjs.com/package/dokugen)