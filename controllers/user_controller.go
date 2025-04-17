package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
		Email:       registerRequest.Email,
		Password:    registerRequest.Password,
		AccountType: models.AccountTypeCustomer,
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

	token, err := utils.GenerateToken(user.ID, user.AccountType)
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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uuid.UUID)
	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		utils.Logger.Error("get user failed", zap.String("user_id", userID.String()), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// SpecialEmployeeEndpoint godoc
// @Summary Special command for Sort employees
// @Description This route is only accessible by Sort employees (account_type=employee)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]string "Success message"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /users/special [post]
func (ctrl *UserController) SpecialEmployeeEndpoint(c *gin.Context) {
	userID := c.GetString("user_id")
	accountType := c.GetString("account_type")

	utils.Logger.Info("special employee endpoint triggered",
		zap.String("user_id", userID),
		zap.String("account_type", accountType),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Sort employee command executed successfully",
		"user_id": userID,
		"role":    accountType,
	})
}

// CreateEmployee godoc
// @Summary Create a Sort employee account
// @Description Allows creation of a user with employee privileges
// @Tags admin
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param request body dto.CreateEmployeeRequest true "Employee creation data"
// @Success 201 {object} map[string]any "Created employee in 'user' field"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 500 {object} map[string]string "Server error"
// @Router /users/create-employee [post]
func (ctrl *UserController) CreateEmployee(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:       req.Email,
		Password:    req.Password,
		AccountType: "employee",
	}

	if err := ctrl.service.Register(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uuid.UUID)
	var updateRequest dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&updateRequest); err != nil {
		utils.Logger.Warn("invalid user update input", zap.String("user_id", userID.String()), zap.Error(err))
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
	updatedUser, err := ctrl.service.UpdateUser(userID, &user)
	if err != nil {
		utils.Logger.Error("profile update failed", zap.String("user_id", userID.String()), zap.Error(err))
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
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDRaw.(uuid.UUID)

	if err := ctrl.service.DeleteUser(userID); err != nil {
		utils.Logger.Error("user delete failed", zap.String("user_id", userID.String()), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
