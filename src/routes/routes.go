package routes

import (
	"biblia-api_v2/src/controllers" // Importação da pasta controllers dentro de src

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/books", controllers.GetBooks)
	r.POST("/pesquisar", controllers.SearchVerses)
	r.POST("/capitulos", controllers.GetChapters)
	r.POST("/books_by_category", controllers.GetBooksByCategory)
	r.POST("/categories", controllers.GetCategories)
	r.POST("/verse_of_the_day", controllers.GetVerseOfTheDay)

	return r
}
