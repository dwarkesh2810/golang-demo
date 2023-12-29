package conf

import (
	"github.com/spf13/viper"
)

type Config struct {
	JwtSecret      string `mapstructure:"JWT_SEC_KEY"`
	FromEmail      string `mapstructure:"EMAIL"`
	Password       string `mapstructure:"PASSWORD"`
	SmtpHost       string `mapstructure:"SMTP_HOST"`
	SmtpPort       string `mapstructure:"SMTP_PORT"`
	BaseUploadPath string `mapstructure:"BASE_UPLOAD_PATH"`
	RateLimiter    int    `mapstructure:"RATELIMITER"`
	BlockTime      int64  `mapstructure:"BLOCKTIME"`
	ConnString     string `mapstructure:"CONN"`
	DbDriver       string `mapstructure:"DBDRIVER"`
}

func LoadConfig(path, configName, configType string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
