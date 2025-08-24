package app

import (
	"github.com/kdaxx/container-app/app/conf"
	"github.com/kdaxx/container-app/app/internal"
	"github.com/kdaxx/container/v2/api"
)

type Registrar struct {
}

func (r *Registrar) RegisterBeans(beanRegister api.BeanRegister) {
	beanRegister.RegisterBeans([]any{
		conf.NewAppConfig(),
		conf.NewLoggerConfig(),

		internal.NewAppLogger(),
		internal.NewAppConfigInjector(),
	})
}
