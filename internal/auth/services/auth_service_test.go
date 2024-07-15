package services

import (
	"encoding/json"
	"fmt"
	"hr-system-go/app/plugins/env"
	http_server "hr-system-go/app/plugins/http"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/app/plugins/mysql"
	"hr-system-go/app/plugins/redis"
	"hr-system-go/internal/auth/constants"
	auth_models "hr-system-go/internal/auth/models"
	user_models "hr-system-go/internal/user/models"
	http "net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AuthService Suite")
}

var (
	authService AuthServiceInterface
	mockLogger  *logger.Logger
	mockEnv     *env.Env
	mockDB      *mysql.MySqlStore
	mockRDS     *redis.RedisStore
)

var _ = BeforeSuite(func() {
	mockEnv = env.NewEnv()
	mockLogger = logger.NewLogger(mockEnv)
	http_server.NewRouter(mockEnv, mockLogger)
	mockDB = mysql.NewMySqlStore(mockEnv, mockLogger)
	mockRDS = redis.NewRedisStore(mockEnv, mockLogger)
	authService = NewAuthService(mockLogger, mockEnv, mockDB, mockRDS)

	mockDB.Connect(
		mockEnv.GetEnv("DB_USER"),
		mockEnv.GetEnv("DB_PASSWORD"),
		mockEnv.GetEnv("DB_DATABASE"),
		mockEnv.GetEnv("DB_HOST"),
		mockEnv.GetEnv("DB_PORT"),
		mockEnv.GetEnv("DB_PARAMS"),
	)

	redisDB, _ := strconv.Atoi(mockEnv.GetEnv("REDIS_DB"))
	mockRDS.Connect(
		mockEnv.GetEnv("REDIS_HOST"),
		mockEnv.GetEnv("REDIS_PORT"),
		redisDB,
	)

	mockDB.DB().AutoMigrate(&user_models.User{}, &auth_models.Role{}, &auth_models.Ability{})
})

var _ = AfterSuite(func() {
	mockRDS.ClearAll()
	mockDB.DB().Migrator().DropTable(&user_models.User{}, &auth_models.Role{}, &auth_models.Ability{})
	mockDB.Close()
})

var _ = Describe("AuthService", func() {
	Describe("GenerateToken", func() {
		It("should generate a valid token", func() {
			token, err := authService.GenerateToken(1, "testuser")
			Expect(err).To(BeNil())
			Expect(token).NotTo(BeEmpty())
		})
	})

	Describe("AbleToAccessOtherUserData", func() {
		var ctx *gin.Context
		BeforeEach(func() {
			ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
		})

		Context("when user has admin ability", func() {
			It("should return true", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: constants.ABILITY_ADMIN}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)

				result := authService.AbleToAccessOtherUserData(ctx, int(user.ID), "some_ability")
				Expect(result).To(BeTrue())
			})
		})

		Context("when user has the required ability", func() {
			It("should return true", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: "required_ability"}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)

				result := authService.AbleToAccessOtherUserData(ctx, int(user.ID), "required_ability")
				Expect(result).To(BeTrue())
			})
		})

		Context("when user is accessing their own data", func() {
			It("should return true", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: "some_ability"}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)

				result := authService.AbleToAccessOtherUserData(ctx, int(user.ID), "some_ability")
				Expect(result).To(BeTrue())
			})
		})

		Context("when user doesn't have permission", func() {
			It("should return false", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: "some_ability"}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)

				result := authService.AbleToAccessOtherUserData(ctx, 999, "required_ability")
				Expect(result).To(BeFalse())
			})
		})

		Context("When abilities are cached in Redis", func() {
			It("should return true with cache abilities", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: constants.ABILITY_ADMIN}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)
				redisKey := fmt.Sprintf("cache:users/%v/abilitiesNames", int(user.ID))
				mockRDS.Set(redisKey, []string{constants.ABILITY_ADMIN}, 24*time.Hour)

				var abilities []string
				Expect(mockRDS.Get(redisKey, &abilities)).To(BeNil())
				result := authService.AbleToAccessOtherUserData(ctx, int(user.ID), "required_ability")
				Expect(result).To(BeTrue())
			})

			It("should return true without cache abilities", func() {
				user := &user_models.User{
					Email: faker.Email(),
					Role: &auth_models.Role{
						Abilities: []auth_models.Ability{{Name: constants.ABILITY_ADMIN}},
					},
				}
				mockDB.DB().Create(&user)
				ctx.Set("currentUser", user)
				redisKey := fmt.Sprintf("cache:users/%v/abilitiesNames", int(user.ID))

				var abilities []string
				Expect(mockRDS.Get(redisKey, &abilities)).NotTo(BeNil())
				result := authService.AbleToAccessOtherUserData(ctx, int(user.ID), "required_ability")
				Expect(result).To(BeTrue())
			})
		})
	})

	Describe("AuthTokenWrapper", func() {
		var (
			w *httptest.ResponseRecorder
			c *gin.Context
			r *http.Request
		)

		BeforeEach(func() {
			gin.SetMode(gin.TestMode)
			w = httptest.NewRecorder()
			c, _ = gin.CreateTestContext(w)
		})

		It("should return 401 when Authorization header is missing", func() {
			r, _ = http.NewRequest(http.MethodGet, "/", nil)
			c.Request = r

			handler := authService.AuthTokenWrapper(func(c *gin.Context) {})
			handler(c)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).To(Equal("Authorization header is required"))
		})

		It("should return 401 when token is invalid", func() {
			r, _ = http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Authorization", "invalid_token")
			c.Request = r

			handler := authService.AuthTokenWrapper(func(c *gin.Context) {})
			handler(c)

			Expect(w.Code).To(Equal(http.StatusUnauthorized))
			var response map[string]string
			json.Unmarshal(w.Body.Bytes(), &response)
			Expect(response["error"]).To(Equal("Invalid token"))
		})

		It("should call next handler when token is valid", func() {
			user := &user_models.User{
				Email: faker.Email(),
				Role: &auth_models.Role{
					Abilities: []auth_models.Ability{{Name: "some_ability"}},
				},
			}

			mockDB.DB().Create(&user)
			validToken, _ := authService.GenerateToken(user.ID, "testuser")

			r, _ = http.NewRequest(http.MethodGet, "/", nil)
			r.Header.Set("Authorization", validToken)
			c.Request = r

			handlerCalled := false
			handler := authService.AuthTokenWrapper(func(c *gin.Context) {
				handlerCalled = true
			})
			handler(c)

			Expect(handlerCalled).To(BeTrue())
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
