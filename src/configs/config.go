package configs

import (
	"fmt"
	"log"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

var AppConfig App

type App struct {
	Host    string `json:"host" default:"0.0.0.0" envconfig:"HOST"`
	Port    int    `default:"8080" envconfig:"PORT"`
	RunMode string `default:"debug" envconfig:"RUN_MODE"`
	Env     string `default:"debug" envconfig:"ENV"`

	Postgres Postgres
	Config   Config
}

// AddressListener returns address listener of HTTP server.
func (c *App) AddressListener() string {
	return fmt.Sprintf("%v:%v", c.Host, c.Port)
}

type Config struct {
	DatabaseURL      string `mapstructure:"database.url"`
	DatabasePort     int    `mapstructure:"database.port"`
	DatabaseUsername string `mapstructure:"database.username"`
	DatabasePassword string `mapstructure:"database.password"`
}

type Postgres struct {
	Username          string `default:"vin_id" envconfig:"POSTGRES_USER"`
	Password          string `default:"vin_id" envconfig:"POSTGRES_PASS"`
	Host              string `default:"127.0.0.1" envconfig:"POSTGRES_HOST"`
	Port              int    `default:"3306" envconfig:"POSTGRES_PORT"`
	Database          string `default:"gamezone" envconfig:"POSTGRES_DB"`
	MaxOpenConnection int    `default:"10" envconfig:"POSTGRES_MAX_OPEN"`
	MaxIdleConnection int    `default:"10" envconfig:"POSTGRES_MAX_IDLE"`
	MaxLifeTime       int    `default:"24" envconfig:"POSTGRES_MAX_LIFETIME"`
}

func Load() (*App, error) {
	AppConfig := App{}
	if err := envconfig.Process("", &AppConfig); err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	// v.SetConfigType("yaml")

	// if err := v.ReadConfig(bytes.NewBuffer([]byte(defaultValue))); err != nil {
	// 	return nil, err
	// }

	v.ReadInConfig()

	err := v.Unmarshal(&AppConfig.Config, func(c *mapstructure.DecoderConfig) {
		c.TagName = "json"

	})

	log.Printf("host: %v", v.GetString("host"))

	log.Printf("AppConfig: %+v", AppConfig)

	// v.SetConfigFile("config.yaml")
	// v.ReadInConfig()
	// v.Unmarshal(&cfg, func(c *envconfig.DecoderConfig) {
	// 	c.TagName = "json"
	// })

	return &AppConfig, err
}
