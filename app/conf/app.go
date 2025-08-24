package conf

// AppConfig loads app config
type AppConfig struct {
	Mode string `mapstructure:"mode"`
}

func (a *AppConfig) Prefix() string {
	return "app"
}

func NewAppConfig() *AppConfig {
	return &AppConfig{}
}
