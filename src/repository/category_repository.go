package repository

import (
	"database/sql"
	"biblia-api_v2/src/models"
)

func GetAllCategories(db *sql.DB) ([]models.Category, error) {
	rows, err := db.Query("SELECT id, name FROM book_categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func GetBooksByCategory(db *sql.DB, idCategory int) ([]models.Book, error) {
	rows, err := db.Query(`
		SELECT b.id, b.name, b.abbrev, b.testament
		FROM books b
		JOIN book_categories_relation bc ON b.id = bc.book_id
		WHERE bc.category_id = ?`, idCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Name, &book.Abbrev, &book.Testament); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
