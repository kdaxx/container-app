package conf

import "github.com/sirupsen/logrus"

// LogConfig loads log config
type LogConfig struct {
	Filepath   string `mapstructure:"filepath"`
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	MaxBackups int    `mapstructure:"max-backups"`
	MaxAge     int    `mapstructure:"max-age"`
	MaxSize    int64  `mapstructure:"max-size"`
	Compress   bool   `mapstructure:"compress"`
	Mode       string `mapstructure:"mode"`
}

func (c *LogConfig) Prefix() string {
	return "log"
}

func (c *LogConfig) ConfigurationBean() any {
	return c
}

func NewLoggerConfig() *LogConfig {
	return &LogConfig{
		Level:      logrus.InfoLevel.String(),
		Format:     "2006-01-02-15-04-05.000",
		Filepath:   "log/app.log",
		MaxBackups: 10,
		MaxAge:     30,
		MaxSize:    10 * 1024 * 1024,
		Compress:   false,
	}
}
