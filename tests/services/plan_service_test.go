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

func TestCreatePlan(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db

	plan := models.Plan{
		Name:    "Basic Plan",
		Credits: 10,
	}

	mock.ExpectQuery("INSERT INTO plans \\(name, credits\\) VALUES \\(\\$1, \\$2\\) RETURNING id").
		WithArgs(plan.Name, plan.Credits).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdPlan, err := services.NewPlanService().CreatePlan(plan)
	assert.Equal(t, 1, createdPlan.ID)
	assert.Equal(t, plan.Name, createdPlan.Name)
	assert.Equal(t, plan.Credits, createdPlan.Credits)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db

	planID := 1
	plan := models.Plan{
		ID:      planID,
		Name:    "Basic Plan",
		Credits: 10,
	}

	mock.ExpectQuery("SELECT id, name, credits FROM plans WHERE id = \\$1").
		WithArgs(planID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "credits"}).
			AddRow(plan.ID, plan.Name, plan.Credits))

	retrievedPlan, err := services.NewPlanService().GetPlanById(planID)
	assert.Equal(t, plan, retrievedPlan)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPlanByIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	database.Db = db

	planID := 999

	mock.ExpectQuery("SELECT id, name, credits FROM plans WHERE id = \\$1").
		WithArgs(planID).
		WillReturnError(sql.ErrNoRows)

	_, errPlan := services.NewPlanService().GetPlanById(planID)
	assert.Equal(t, errors.ErrPlanNotFound, errPlan.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}
