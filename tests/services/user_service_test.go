package services_test

import (
	"assignment/database"
	"assignment/internal/errors"
	"assignment/internal/models"
	"assignment/internal/services"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()
	user := models.User{
		Username: "testuser",
		Password: "password123",
		PlanId:   1,
	}

	mock.ExpectQuery("SELECT credits from plans where id = ?").
		WithArgs(user.PlanId).
		WillReturnRows(sqlmock.NewRows([]string{"credits"}).AddRow(100))

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.Password, user.PlanId, 100).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdUser, _ := userService.CreateUser(user)
	assert.Equal(t, user, createdUser)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserPlanNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	user := models.User{
		Username: "testuser",
		Password: "password123",
		PlanId:   999,
	}

	mock.ExpectQuery("SELECT credits from plans where id = ?").
		WithArgs(user.PlanId).
		WillReturnError(sql.ErrNoRows)

	_, errCreate := userService.CreateUser(user)
	assert.Equal(t, errors.ErrPlanNotFound, errCreate.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	user := models.User{
		ID:       1,
		Username: "testuser",
		Password: "password123",
		PlanId:   1,
		Credits:  100,
	}

	mock.ExpectQuery("SELECT id, username, password, credits, plan_id FROM users WHERE id = ?").
		WithArgs(user.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "credits", "plan_id"}).
			AddRow(user.ID, user.Username, user.Password, user.Credits, user.PlanId))

	result, _ := userService.GetUserById(user.ID)
	assert.Equal(t, user, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	mock.ExpectQuery("SELECT id, username, password, credits, plan_id FROM users WHERE id = ?").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	_, errGet := userService.GetUserById(999)
	assert.Equal(t, errors.ErrUserNotFound, errGet.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	user := models.User{
		ID:       1,
		Username: "testuser",
		Password: "password123",
		PlanId:   1,
		Credits:  100,
	}

	mock.ExpectQuery("SELECT id, username, password, credits, plan_id FROM users WHERE username = ?").
		WithArgs(user.Username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "credits", "plan_id"}).
			AddRow(user.ID, user.Username, user.Password, user.Credits, user.PlanId))

	result, _ := userService.GetUserByUsername(user.Username)
	assert.Equal(t, user, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsernameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	mock.ExpectQuery("SELECT id, username, password, credits, plan_id FROM users WHERE username = ?").
		WithArgs("nonexistentuser").
		WillReturnError(sql.ErrNoRows)

	_, errorGet := userService.GetUserByUsername("nonexistentuser")
	assert.Equal(t, errors.ErrUserNotFound, errorGet.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeductCredits(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	userID := 1
	currentCredits := 10

	mock.ExpectBegin()

	mock.ExpectQuery("SELECT credits FROM users WHERE id = \\$1 FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"credits"}).AddRow(currentCredits))

	mock.ExpectExec("UPDATE users SET credits = credits - \\$1 WHERE id = \\$2").
		WithArgs(1, userID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	errCredit := userService.DeductCredits(userID)
	assert.Nil(t, errCredit)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCredits(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	userID := 1
	currentCredits := 10

	mock.ExpectQuery("SELECT credits FROM users WHERE id = \\$1").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"credits"}).AddRow(currentCredits))

	credits, err := userService.GetCredits(userID)
	assert.Equal(t, currentCredits, credits)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeductCreditsInsufficientStartingBalance(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db
	userService := services.NewUserService()

	userID := 1
	currentCredits := 0

	mock.ExpectBegin()

	mock.ExpectQuery("SELECT credits FROM users WHERE id = \\$1 FOR UPDATE").
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"credits"}).AddRow(currentCredits))

	mock.ExpectRollback()

	errDeducate := userService.DeductCredits(userID)
	assert.Equal(t, errors.ErrInsufficientCredits, errDeducate.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
