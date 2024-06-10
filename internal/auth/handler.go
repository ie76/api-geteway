package auth

import (
	"assignment/internal/errors"
	"assignment/internal/models"
	"assignment/internal/services"
	"assignment/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser godoc
//
// @Summary Login a new user
// @Description Login a user
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   LoginRequest body LoginRequest true "Login request"
// @Success 201 {object} LoginResponse
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the user exists and the password is correct
	user, err := authenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenString, err := utils.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// RegisterUser godoc
//
// @Summary Register a new user
// @Description Register a new user with a username, password and plan id
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   RegisterRequest body RegisterRequest true "Register request"
// @Success 201 {object} RegisterResponse
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.GenerateFromPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		Username: req.Username,
		PlanId:   req.PlanId,
		Password: hashedPassword,
	}

	createdUser, errCreate := services.NewUserService().CreateUser(user)
	if errCreate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"at": errCreate.Time, "message": errCreate.Message, "code": errCreate.Code})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{Message: "User registered successfully", User: createdUser})

}

func authenticateUser(username, password string) (*models.User, *errors.Error) {

	user, err := services.NewUserService().GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	errorHash := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errorHash != nil {
		return nil, errors.New(errors.ErrHash)
	}

	return &user, nil
}
