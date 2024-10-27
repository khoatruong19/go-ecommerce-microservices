package config

import (
	"github.com/khoatruong19/go-ecommerce-microservices/internal/pkg/config/environment"
	"go.uber.org/fx"
)

const ConfigFxName = "configfx"

var Module = fx.Module(ConfigFxName, fx.Provide(func() environment.Environment {
	return environment.ConfigAppEnv()
}))

var ModuleFunc = func(e environment.Environment) fx.Option {
	return fx.Module(ConfigFxName, fx.Provide(func() environment.Environment {
		return environment.ConfigAppEnv(e)
	}))
}
