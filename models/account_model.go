package models

type Account struct {
	ID      int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
