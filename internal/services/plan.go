package services

import (
	"assignment/database"
	"assignment/internal/errors"
	"assignment/internal/models"
	"database/sql"
)

type PlanService interface {
	CreatePlan(plan models.Plan) (models.Plan, *errors.Error)
	GetPlanById(id int) (models.Plan, *errors.Error)
}

type IPlanService struct {
	db *sql.DB
}

func NewPlanService() PlanService {
	return &IPlanService{
		db: database.GetDB(),
	}
}

func (s *IPlanService) CreatePlan(plan models.Plan) (models.Plan, *errors.Error) {
	query := "INSERT INTO plans (name, credits) VALUES ($1, $2) RETURNING id"
	err := s.db.QueryRow(query, plan.Name, plan.Credits).Scan(&plan.ID)
	if err != nil {
		return models.Plan{}, errors.New(errors.ErrCreatePlan)
	}
	return plan, nil
}

func (s *IPlanService) GetPlanById(id int) (models.Plan, *errors.Error) {
	query := "SELECT id, name, credits FROM plans WHERE id = $1"
	row := s.db.QueryRow(query, id)

	var plan models.Plan
	err := row.Scan(&plan.ID, &plan.Name, &plan.Credits)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Plan{}, errors.New(errors.ErrPlanNotFound)
		}
		return models.Plan{}, errors.New(errors.ErrDatabaseError)
	}

	return plan, nil
}
