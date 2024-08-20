package controllers

import (
	"net/http"
	"biblia-api_v2/src/database"
	"biblia-api_v2/src/repository"

	"github.com/gin-gonic/gin"
)

func GetBooks(c *gin.Context) {
	var requestBody struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	books, err := repository.GetAllBooks(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}
