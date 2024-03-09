package configs

import "github.com/spf13/viper"

type Conf struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	TokenLimit int    `mapstructure:"TOKEN_LIMIT"`
	TimeToken  int    `mapstructure:"TIME_TOKEN"`
	IPLimit    int    `mapstructure:"IP_LIMIT"`
	TimeIP     int    `mapstructure:"TIME_IP"`
}

func NewConf(tokenLimit, ipLimit int) *Conf {
	return &Conf{
		TokenLimit: tokenLimit,
		IPLimit:    ipLimit,
	}

}

func LoadConfig(path string) (*Conf, error) {
	var c *Conf
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

func GetTokenLimit() (int, int) {
	c, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return c.TokenLimit, c.TimeToken
}

func GetIpLimit() (int, int) {
	c, err := LoadConfig(".")
	if err != nil {
		panic(err)
	}
	return c.IPLimit, c.TimeIP
}
