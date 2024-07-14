package controllers

import (
	"hr-system-go/app/plugins/logger"
	auth_service "hr-system-go/internal/auth/services"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/internal/user/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type SessionsController struct {
	logger      *logger.Logger
	service     services.UserServiceInterface
	authService auth_service.AuthServiceInterface
}

func NewSessionsController(logger *logger.Logger, service services.UserServiceInterface, authService auth_service.AuthServiceInterface) *SessionsController {
	return &SessionsController{
		logger:      logger,
		service:     service,
		authService: authService,
	}
}

func (c *SessionsController) RegisterRoutes(r *gin.Engine) {
	r.POST("api/register", c.SignUp)
	r.POST("api/login", c.SignIn)
	r.POST("api/passwordResetRequest", c.PasswordResetRequest)
	r.POST("api/resetPassword", c.authService.AuthTokenWrapper(c.ResetPassword))
}

type sessionBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type passwordResetRequestBody struct {
	Email string `json:"email"`
}

type resetPasswordBody struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
}

func (c *SessionsController) SignUp(ctx *gin.Context) {
	var payload sessionBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot Parse Body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := user_models.User{
		Name:  payload.Name,
		Email: payload.Email,
	}
	if err := c.service.RegisterUser(&user, payload.Password); err != nil {
		c.logger.Error("Cannot Register User", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to register user"})
		return
	}

	token, _ := c.authService.GenerateToken(user.ID, user.Name)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *SessionsController) SignIn(ctx *gin.Context) {
	var payload sessionBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot Parse Body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *user_models.User
	user, err := c.service.FindUserByEmail(payload.Email)
	if err != nil {
		c.logger.Error("Cannot Find User", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordEncrypt), []byte(payload.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := c.authService.GenerateToken(user.ID, user.Name)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *SessionsController) PasswordResetRequest(ctx *gin.Context) {
	var payload passwordResetRequestBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot Parse Body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.FindUserByEmail(payload.Email)
	if err != nil {
		c.logger.Error("Cannot Find User by Email")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to Send Reset Request"})
	}

	// token is for post resetPassword action
	// todo: send token by mailer
	token, _ := c.authService.GenerateToken(user.ID, user.Name)
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (c *SessionsController) ResetPassword(ctx *gin.Context) {
	var payload resetPasswordBody
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		c.logger.Error("Cannot Parse Body", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := c.authService.GetCurrentUser(ctx)
	if user == nil {
		c.logger.Error("Cannot Get Current User")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
	}

	if err := c.service.UpdatePassword(user, payload.NewPassword); err != nil {
		c.logger.Error("Cannot Find User", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
