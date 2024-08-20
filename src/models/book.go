package models

type Book struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Abbrev    string `json:"abbrev"`
	Testament string `json:"testament"`
}
