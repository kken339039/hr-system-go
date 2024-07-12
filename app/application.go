package app

import (
	"hr-system-go/app/plugins"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

type Application struct {
	fx      *fx.App
	plugins []interface{}
}

type AppLogger struct {
	fxevent.Logger
}

// todo: output fx event log by different action, eg: start, stop
func (AppLogger) LogEvent(event fxevent.Event) {}

func NewApplication() *Application {
	return &Application{
		plugins: plugins.Registry,
	}
}

func (a *Application) AddModule(module AppModuleInterface) {
	controllers := module.Controllers()
	provides := module.Provide()

	if len(controllers) > 0 {
		a.plugins = append(a.plugins, module.Controllers()...)
	}

	if len(provides) > 0 {
		a.plugins = append(a.plugins, module.Provide()...)
	}
}

func (app *Application) Run(funcs ...interface{}) {
	app.fx = fx.New(
		fx.WithLogger(func() fxevent.Logger { return AppLogger{} }),
		fx.Provide(
			app.plugins...,
		),
		fx.Invoke(funcs...),
	)
	app.fx.Run()
}
