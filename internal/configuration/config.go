package configuration

import (
	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper

	S3 S3Config `mapstructure:"s3"`
	Mongo MongoDBConfig `mapstructure:"mongodb"`
}

type S3Config struct {
	Endpoint string
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	SslEnabled bool `mapstructure:"ssl"`
	Buckets S3BucketsConfig `mapstructre:"buckets"`
}

type S3BucketsConfig struct {
	ImagesBucket string `mapstructure:"images"`
	RecordingsBucket string `mapstructure:"recordings"`
	LogsBucket string `mapstructure:"logs"`
}

type MongoDBConfig struct {
	Uri string
	Database string
}

var appConfig *Config

func LoadConfig() error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	appConfig = &Config{}

	v.Unmarshal(appConfig)
	appConfig.v = v

	return nil
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("Please call LoadConfig before GetConfig!")
	}
	return appConfig
}
