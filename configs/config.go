package configs

import "github.com/spf13/viper"

type conf struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	TokenLimit int    `mapstructure:"TOKEN_LIMIT"`
	IPLimit    int    `mapstructure:"IP_LIMIT"`
}

type Conf struct {
	TokenLimit int `mapstructure:"TOKEN_LIMIT"`
	IPLimit    int `mapstructure:"IP_LIMIT"`
}

func NewConf(tokenLimit, ipLimit int) *Conf {
	return &Conf{
		TokenLimit: tokenLimit,
		IPLimit:    ipLimit,
	}

}

func LoadConfig(path string) (*conf, error) {
	var c *conf
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		panic(err)
	}
	return c, err
}

func GetTokenLimit() int {
	c, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return c.TokenLimit
}

func GetIpLimit() int {
	c, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return c.IPLimit
}
