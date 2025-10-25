package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/shemigam1/hngxiii-currency-exchange/db"
	"github.com/shemigam1/hngxiii-currency-exchange/routes"
)

func init() {
	db.LoadEnvVariables()
	db.ConnectDB()
}

func main() {
	r := gin.Default()
	routes.Routes(r)
	r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
}
