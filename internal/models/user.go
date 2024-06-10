package models

type User struct {
	ID       int    `json:"id" gorm:"primaryKey"`
	PlanId   int    `json:"plan_id"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
	Credits  int    `json:"credits"`
}
