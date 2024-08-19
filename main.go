package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var verseOfTheDay map[string]interface{}

func init() {
	var err error

	db, err = sql.Open("mysql", "root:lobio2541@tcp(127.0.0.1:3306)/biblia")
	if err != nil {
		log.Fatal(err)
	}
}

func getChapters(c *gin.Context) {
	// Verifique o token
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	token, ok := requestBody["token"].(string)
	if !ok || token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	idBook, ok := requestBody["id_book"].(float64) // JSON números são float64
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	// Consultar os capítulos e versos do livro
	rows, err := db.Query(`
		SELECT chapter, verse, text
		FROM verses
		WHERE book = ?
		ORDER BY chapter, verse`, int(idBook))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var chapters []map[string]interface{}
	var currentChapter int
	var chapterVerses []map[string]interface{}

	for rows.Next() {
		var chapter, verse int
		var text string
		if err := rows.Scan(&chapter, &verse, &text); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if chapter != currentChapter {
			if currentChapter != 0 {
				chapters = append(chapters, map[string]interface{}{
					"chapter": currentChapter,
					"verses":  chapterVerses,
				})
			}
			chapterVerses = []map[string]interface{}{}
			currentChapter = chapter
		}

		chapterVerses = append(chapterVerses, map[string]interface{}{
			"verse": verse,
			"text":  text,
		})
	}

	// Adiciona o último capítulo
	if currentChapter != 0 {
		chapters = append(chapters, map[string]interface{}{
			"chapter": currentChapter,
			"verses":  chapterVerses,
		})
	}

	c.JSON(http.StatusOK, chapters)
}

func getBooks(c *gin.Context) {
	// Define uma estrutura para o corpo da requisição
	type RequestBody struct {
		Token string `json:"token"`
	}

	var requestBody RequestBody
	// Faz o binding do corpo da requisição para a estrutura RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Verifica o token no corpo da requisição
	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var id int
		var name, abbrev, testament string
		if err := rows.Scan(&id, &name, &abbrev, &testament); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, map[string]interface{}{
			"id":        id,
			"name":      name,
			"abbrev":    abbrev,
			"testament": testament,
		})
	}

	c.JSON(http.StatusOK, books)
}

func searchVerses(c *gin.Context) {
	var requestBody struct {
		Pesquisa string `json:"pesquisa"`
		Token    string `json:"token"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	searchText := strings.TrimSpace(requestBody.Pesquisa)
	if searchText == "" {
		// Se o texto de pesquisa estiver vazio, retorna uma lista vazia
		c.JSON(http.StatusOK, []map[string]interface{}{})
		return
	}

	// Se o texto de pesquisa não estiver vazio, realiza a pesquisa no banco de dados
	rows, err := db.Query(`
		SELECT v.id, b.name, v.chapter, v.verse, v.text
		FROM verses v
		JOIN books b ON v.book = b.id-1
		WHERE v.text LIKE ?`, "%"+searchText+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var id, chapter, verse int
		var bookName, text string
		if err := rows.Scan(&id, &bookName, &chapter, &verse, &text); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, map[string]interface{}{
			"id":      id,
			"book":    bookName,
			"chapter": chapter,
			"verse":   verse,
			"text":    text,
		})
	}

	c.JSON(http.StatusOK, results)
}

func getBooksByCategory(c *gin.Context) {
	var requestBody struct {
		IDCategory int    `json:"id_category"`
		Token      string `json:"token"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	idCategory := requestBody.IDCategory
	if idCategory <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	rows, err := db.Query(`
		SELECT b.id, b.name, b.abbrev, b.testament
		FROM books b
		JOIN book_categories_relation bc ON b.id = bc.book_id
		WHERE bc.category_id = ?`, idCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []map[string]interface{}
	for rows.Next() {
		var id int
		var name, abbrev, testament string
		if err := rows.Scan(&id, &name, &abbrev, &testament); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, map[string]interface{}{
			"id":        id,
			"name":      name,
			"abbrev":    abbrev,
			"testament": testament,
		})
	}

	c.JSON(http.StatusOK, books)
}

func getCategories(c *gin.Context) {
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

	rows, err := db.Query("SELECT id, name FROM book_categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var categories []map[string]interface{}
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		categories = append(categories, map[string]interface{}{
			"id":   id,
			"name": name,
		})
	}

	c.JSON(http.StatusOK, categories)
}

func updateVerseOfTheDay() {
	row := db.QueryRow(`
        SELECT v.id, b.name, v.chapter, v.verse, v.text
        FROM verses v
        JOIN books b ON v.book = b.id-1
        ORDER BY RAND()
        LIMIT 1`)

	var id, chapter, verse int
	var bookName, text string
	if err := row.Scan(&id, &bookName, &chapter, &verse, &text); err != nil {
		log.Printf("Erro ao escanear o resultado: %v", err)
		return
	}

	verseOfTheDay = map[string]interface{}{
		"id":      id,
		"book":    bookName,
		"chapter": chapter,
		"verse":   verse,
		"text":    text,
	}

	log.Printf("Versículo do dia atualizado: %v", verseOfTheDay)
}

func getVerseOfTheDay(c *gin.Context) {
	var requestBody struct {
		Token string `json:"token"`
	}

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if requestBody.Token != "biblia" {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}

	updateVerseOfTheDay()
	c.JSON(200, verseOfTheDay)
}

func main() {
	r := gin.Default()
	r.POST("/books", getBooks)
	r.POST("/pesquisar", searchVerses)
	r.POST("/capitulos", getChapters)
	r.POST("/books_by_category", getBooksByCategory)
	r.POST("/categories", getCategories)
	r.POST("/verse_of_the_day", getVerseOfTheDay)
	r.Run(":8080")
}
