package database

import (
	"assignment/config"
	"assignment/internal/errors"
	"assignment/internal/models"
	"assignment/internal/utils"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Migrate() {

	cfg := config.GetConfig()
	connString := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		panic(errors.New(errors.ErrDatabaseError))
	}

	plans := []models.Plan{
		{Credits: 5, Name: "Basic"},
		{Credits: 15, Name: "Medium"},
		{Credits: 25, Name: "Premium"},
	}

	encryptPassword, err := utils.GenerateFromPassword("password")
	users := []models.User{
		{PlanId: 1, Username: "username_basic", Password: encryptPassword, Credits: 5},
		{PlanId: 2, Username: "username_medium", Password: encryptPassword, Credits: 15},
		{PlanId: 2, Username: "username_premium", Password: encryptPassword, Credits: 25},
	}

	db.AutoMigrate(&users, &plans)
	db.Create(plans)
	db.Create(users)
}
