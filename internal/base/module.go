package user

import (
	"hr-system-go/app"
)

type BaseModule struct {
	app.AppModuleInterface
}

func (m *BaseModule) Provide() []interface{} {
	return []interface{}{
		// models.NewBaseModel,
	}
}
