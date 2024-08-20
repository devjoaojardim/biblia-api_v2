package controllers

import (
	"net/http"
	"biblia-api_v2/src/database"
	"biblia-api_v2/src/repository"

	"github.com/gin-gonic/gin"
)

func GetChapters(c *gin.Context) {
	var requestBody struct {
		Token  string  `json:"token"`
		IDBook float64 `json:"id_book"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	if requestBody.IDBook <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	verses, err := repository.GetVersesByBook(database.DB, int(requestBody.IDBook))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verses)
}

func SearchVerses(c *gin.Context) {
	var requestBody struct {
		Token    string `json:"token"`
		Pesquisa string `json:"pesquisa"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	searchText := requestBody.Pesquisa
	if searchText == "" {
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	verses, err := repository.SearchVerses(database.DB, searchText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verses)
}

func GetVerseOfTheDay(c *gin.Context) {
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

	verse, err := repository.GetVerseOfTheDay(database.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, verse)
}
