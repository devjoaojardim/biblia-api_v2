package main

import (
	"biblia-api_v2/src/database"
	"biblia-api_v2/src/routes"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func init() {
	database.Connect()
	router = routes.SetupRouter()
}

func main() {
	router.Run(":8080")
}
