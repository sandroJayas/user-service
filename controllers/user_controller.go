package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/sandroJayas/user-service/dto"
	"github.com/sandroJayas/user-service/models"
	"github.com/sandroJayas/user-service/usecase"
	"github.com/sandroJayas/user-service/utils"
	"go.uber.org/zap"
	"net/http"
)

type UserController struct {
	service *usecase.UserService
}

func NewUserController(service *usecase.UserService) *UserController {
	return &UserController{service: service}
}

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param registerRequest body dto.RegisterRequest true "Registration data"
// @Success 201 {object} map[string]any "Created user object in 'user' field"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/register [post]
func (ctrl *UserController) Register(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.Logger.Warn("invalid user input", zap.String("email", registerRequest.Email), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:    registerRequest.Email,
		Password: registerRequest.Password,
	}
	if err := ctrl.service.Register(&user); err != nil {
		utils.Logger.Error("user registration failed", zap.String("email", registerRequest.Email), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates user with email and password, returns JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param loginRequest body dto.LoginRequest true "Login data"
// @Success 200 {object} map[string]any "JWT token in 'token' field"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid credentials"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.Logger.Warn("invalid login", zap.String("email", loginRequest.Email), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.service.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.Logger.Warn("login failed", zap.String("email", loginRequest.Email), zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		utils.Logger.Error("token generation failed", zap.String("user_id", user.ID.String()), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})

		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Me godoc
// @Summary Get current user
// @Description Returns the user data for the authenticated user
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} map[string]any "User object in 'user' field"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/me [get]
func (ctrl *UserController) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := ctrl.service.GetUserByID(userID.(string))
	if err != nil {
		utils.Logger.Error("get user failed", zap.String("user_id", userID.(string)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfile godoc
// @Summary Update user's profile
// @Description Updates the logged-in user's profile fields
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param updateRequest body dto.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} map[string]any "Updated user in 'user' field"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/profile [put]
func (ctrl *UserController) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var updateRequest dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.Logger.Warn("invalid user update input", zap.String("user_id", userID.(string)), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{
		FirstName:       updateRequest.FirstName,
		LastName:        updateRequest.LastName,
		AddressLine1:    updateRequest.AddressLine1,
		AddressLine2:    updateRequest.AddressLine2,
		City:            updateRequest.City,
		PostalCode:      updateRequest.PostalCode,
		Country:         updateRequest.Country,
		PhoneNumber:     updateRequest.PhoneNumber,
		PaymentMethodID: updateRequest.PaymentMethodID,
	}
	updatedUser, err := ctrl.service.UpdateUser(userID.(string), &user)
	if err != nil {
		utils.Logger.Error("profile update failed", zap.String("user_id", userID.(string)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser})
}

// DeleteUser godoc
// @Summary Soft-delete the current user
// @Description Marks the user as deleted (is_deleted = true)
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} map[string]any "Deletion confirmation"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/delete [delete]
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID, _ := c.Get("user_id")

	if err := ctrl.service.DeleteUser(userID.(string)); err != nil {
		utils.Logger.Error("user delete failed", zap.String("user_id", userID.(string)), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
