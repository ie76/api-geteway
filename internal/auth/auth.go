package auth

import "assignment/internal/models"

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	Message string      `json:"message"`
	User    models.User `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	PlanId   int    `json:"plan_id" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
