package repository

import (
	"database/sql"
	"biblia-api_v2/src/models"
)

func GetVersesByBook(db *sql.DB, idBook int) ([]models.Verse, error) {
	rows, err := db.Query(`
		SELECT chapter, verse, text
		FROM verses
		WHERE book = ?
		ORDER BY chapter, verse`, idBook)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []models.Verse
	for rows.Next() {
		var verse models.Verse
		if err := rows.Scan(&verse.Chapter, &verse.Verse, &verse.Text); err != nil {
			return nil, err
		}
		verses = append(verses, verse)
	}
	return verses, nil
}

func SearchVerses(db *sql.DB, searchText string) ([]models.Verse, error) {
	rows, err := db.Query(`
		SELECT v.id, b.name, v.chapter, v.verse, v.text
		FROM verses v
		JOIN books b ON v.book = b.id-1
		WHERE v.text LIKE ?`, "%"+searchText+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var verses []models.Verse
	for rows.Next() {
		var verse models.Verse
		if err := rows.Scan(&verse.ID, &verse.Book, &verse.Chapter, &verse.Verse, &verse.Text); err != nil {
			return nil, err
		}
		verses = append(verses, verse)
	}
	return verses, nil
}

func GetVerseOfTheDay(db *sql.DB) (*models.Verse, error) {
	row := db.QueryRow(`
        SELECT v.id, b.name, v.chapter, v.verse, v.text
        FROM verses v
        JOIN books b ON v.book = b.id-1
        ORDER BY RAND()
        LIMIT 1`)

	var verse models.Verse
	if err := row.Scan(&verse.ID, &verse.Book, &verse.Chapter, &verse.Verse, &verse.Text); err != nil {
		return nil, err
	}

	return &verse, nil
}
