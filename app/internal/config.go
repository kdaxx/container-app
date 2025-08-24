package internal

import (
	"errors"
	"fmt"
	"github.com/kdaxx/container-app/app/api"
	"github.com/spf13/viper"
	"reflect"
)

// ConfigInjector injects it's config to all of api.Configuration in this app container
type ConfigInjector struct {
	viper *viper.Viper
}

func (a *ConfigInjector) PreInit() any {
	return func(ric []api.Configuration) error {
		v := viper.New()
		fileName := "application.yaml"
		v.SetConfigFile(fileName)
		v.SetConfigType("yaml")
		// empty treat to nil
		v.AllowEmptyEnv(false)

		if err := v.ReadInConfig(); err != nil {
			if errors.As(err, &viper.ConfigFileNotFoundError{}) {
				// Config file not found; ignore error if desired
				return fmt.Errorf("config file [%s] not found: %v\n", fileName, err)
			} else {
				// Config file was found but another error was produced
				return fmt.Errorf("read config file [%s] failed: %v\n", fileName, err)
			}
		}
		// Config file found and successfully parsed
		a.viper = v

		err := a.injectConfig(ric, v)
		if err != nil {
			return err
		}

		return nil
	}
}

func (a *ConfigInjector) injectConfig(ric []api.Configuration, v *viper.Viper) error {
	for _, bean := range ric {
		err := v.UnmarshalKey(bean.Prefix(), bean)
		if err != nil {
			return fmt.Errorf("unmarshal config %s to %v failed! %v\n",
				bean.Prefix(), reflect.TypeOf(bean), err)
		}
	}
	return nil
}

func NewAppConfigInjector() *ConfigInjector {
	return &ConfigInjector{}
}
