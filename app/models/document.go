package models

type Document struct {
	ID			int64		`db:"id" json:"id"`
	Name		string		`db:"name" json:"name"`
	Language	string		`db:"language" json:"language"`
	Content		string		`db:"content" json:"content"`
}

