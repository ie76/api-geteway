package services

import (
	"assignment/database"
	"assignment/internal/errors"
	"assignment/internal/models"
	"database/sql"
	"sync"
)

type UserService interface {
	CreateUser(user models.User) (models.User, *errors.Error)
	GetUserById(id int) (models.User, *errors.Error)
	GetUserByUsername(username string) (models.User, *errors.Error)
	DeductCredits(userID int) *errors.Error
	GetCredits(userID int) (int, *errors.Error)
}

type IUserService struct {
	user  *models.User
	db    *sql.DB
	mutex sync.Mutex
}

func NewUserService() UserService {
	return &IUserService{
		user: &models.User{},
		db:   database.GetDB(),
	}
}

func (s *IUserService) CreateUser(user models.User) (models.User, *errors.Error) {
	db := s.db

	var planCredits int
	err := db.QueryRow("SELECT credits from plans where id = $1", user.PlanId).Scan(&planCredits)
	if err != nil {
		return models.User{}, errors.New(errors.ErrPlanNotFound)
	}

	query := "INSERT INTO users (username, password, plan_id, credits) VALUES ($1, $2, $3, $4)"
	_, err = db.Exec(query, user.Username, user.Password, user.PlanId, planCredits)
	if err != nil {
		return models.User{}, errors.New(errors.ErrCreateUser)
	}
	return user, nil
}

func (s *IUserService) GetUserById(id int) (models.User, *errors.Error) {
	db := s.db

	query := "SELECT id, username, password, credits, plan_id FROM users WHERE id = $1"
	row := db.QueryRow(query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Credits, &user.PlanId)
	if err != nil {
		return models.User{}, errors.New(errors.ErrUserNotFound)
	}

	return user, nil
}

func (s *IUserService) GetUserByUsername(username string) (models.User, *errors.Error) {
	db := s.db

	query := "SELECT id, username, password, credits, plan_id FROM users WHERE username = $1"
	row := db.QueryRow(query, username)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Credits, &user.PlanId)
	if err != nil {
		return models.User{}, errors.New(errors.ErrUserNotFound)
	}

	return user, nil
}

func (s *IUserService) DeductCredits(userID int) *errors.Error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Start a transaction
	tx, err := s.db.Begin()
	if err != nil {
		return errors.New(errors.ErrDatabaseError)
	}
	defer tx.Rollback()

	var currentCredit int
	err = tx.QueryRow("SELECT credits FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&currentCredit)
	if err != nil {
		return errors.New(errors.ErrUserNotFound)
	}

	if currentCredit <= 0 {
		return errors.New(errors.ErrInsufficientCredits)
	}

	_, err = tx.Exec("UPDATE users SET credits = credits - $1 WHERE id = $2", 1, userID)
	if err != nil {
		return errors.New(errors.ErrDatabaseError)
	}

	err = tx.Commit()
	if err != nil {
		return errors.New(errors.ErrDatabaseError)
	}

	return nil
}

func (s *IUserService) GetCredits(userID int) (int, *errors.Error) {

	var currentCredit int
	err := s.db.QueryRow("SELECT credits FROM users WHERE id = $1", userID).Scan(&currentCredit)
	if err != nil {
		return 0, errors.New(errors.ErrUserNotFound)
	}

	return currentCredit, nil
}
