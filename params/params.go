package params

import "github.com/spf13/viper"

type Config struct {
	TableName  string `mapstructure:"TABLE_NAME"`
	AwsRegion  string `mapstructure:"AWS_REGION"`
	TopNews    string `mapstructure:"TOP_NEWS"`
	SingleNews string `mapstructure:"SINGLE_NEWS"`
	BatchSize  int    `mapstructure:"SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
