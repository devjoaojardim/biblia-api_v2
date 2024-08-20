package repository

import (
	"database/sql"
	"biblia-api_v2/src/models"
)

func GetAllBooks(db *sql.DB) ([]models.Book, error) {
	rows, err := db.Query("SELECT id, name, abbrev, testament FROM books")
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
