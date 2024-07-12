package services

import (
	"hr-system-go/app/plugins/env"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/internal/auth/constants"
	user_models "hr-system-go/internal/user/models"
	"hr-system-go/utils"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	logger *logger.Logger
	env    *env.Env
	db     *mysql.MySqlStore
	jwtKey []byte
}

func NewAuthService(logger *logger.Logger, env *env.Env, db *mysql.MySqlStore) *AuthService {
	return &AuthService{
		logger: logger,
		env:    env,
		db:     db,
		jwtKey: []byte(env.GetEnv("JWT_TOKEN_KEY")),
	}
}

func (s AuthService) AuthTokenWrapper(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.authTokenAndSetCurrentUser()(ctx)
		if ctx.IsAborted() {
			return
		}
		handler(ctx)
	}
}

func (s AuthService) AuthUserAbilityWrapper(handler gin.HandlerFunc, requireAbility string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s.authTokenAndSetCurrentUser()(ctx)
		if ctx.IsAborted() {
			return
		}

		s.authUserAbility(requireAbility)(ctx)
		if ctx.IsAborted() {
			return
		}
		handler(ctx)
	}
}

// if current user has full permissions or is an admin, then they can view anyone's records
func (s AuthService) AbleToAccessOtherUserData(ctx *gin.Context, targetUserId int, allGrantAbility string) bool {
	currentUser := getCurrentUser(ctx)
	for _, abilityName := range currentUser.Role.GetAbilityNames() {
		if abilityName == constants.ABILITY_ADMIN || abilityName == allGrantAbility {
			return true
		}
	}
	return currentUser.ID == uint(targetUserId)
}

func (s AuthService) GenerateToken(userID uint, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":   int(userID),
		"userName": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(s.jwtKey)
}

func (s AuthService) GetCurrentUser(ctx *gin.Context) *user_models.User {
	return getCurrentUser(ctx)
}

func (s AuthService) authUserAbility(requiredAbility string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		currentUser := getCurrentUser(ctx)
		if currentUser.Role == nil || len(currentUser.Role.Abilities) == 0 {
			s.logger.Error("User does not have any ability")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			ctx.Abort()
			return
		}
		hasAbility := false
		for _, ability := range currentUser.Role.Abilities {
			if ability.Name == requiredAbility || ability.Name == constants.ABILITY_ADMIN {
				hasAbility = true
				break
			}
		}

		if !hasAbility {
			s.logger.Error("User does not have enough ability")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (s AuthService) authTokenAndSetCurrentUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		claims, err := ValidateToken(tokenString, s.jwtKey)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		userId, exist := claims["userId"]
		if !exist {
			s.logger.Error("UserId is not exist from token")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		userID, err := utils.ParseInterfaceToInt(userId)
		if err != nil {
			s.logger.Error("Cannot convert UserId")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		var user *user_models.User
		if err := s.db.DB().Preload("Role.Abilities").First(&user, userID).Error; err != nil {
			s.logger.Error("Cannot find user's Role and Ability")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Set("currentUser", user)
		ctx.Set("userName", claims["userName"])
		ctx.Next()
	}
}

func ValidateToken(tokenString string, jwtKey []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}

func getCurrentUser(ctx *gin.Context) *user_models.User {
	var user *user_models.User
	currentUser, ok := ctx.Get("currentUser")
	if !ok {
		return nil
	}
	user, ok = currentUser.(*user_models.User)
	if !ok {
		return nil
	}

	return user
}
