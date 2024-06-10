package models

type Plan struct {
	ID      int    `json:"id"`
	Credits int    `json:"credits"`
	Name    string `json:"name" gorm:"unique"`
}
