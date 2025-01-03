package config

type Config struct {
	Address  string `mapstructure:"address" validate:"required"`
	Username string `mapstructure:"username" validate:"required"`
	Password string `json:"-" mapstructure:"password"`
	DB       int    `default:"0" mapstructure:"db"`
}
