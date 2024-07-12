package auth

import (
	"hr-system-go/app"
	"hr-system-go/app/plugins/logger"
	"hr-system-go/internal/auth/services"
)

type AuthModule struct {
	app.AppModuleInterface
}

func (m *AuthModule) Controllers() []interface{} {
	return []interface{}{
		func(
			logger *logger.Logger,
		) *AuthModule {
			logger.Info("= Auth module init")
			return m
		},
	}
}

func (m *AuthModule) Provide() []interface{} {
	return []interface{}{
		services.NewAuthService,
	}
}
