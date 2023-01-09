package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	v *viper.Viper
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

	appConfig = &Config{v: v}

	return nil
}

func GetConfig() *Config {
	if appConfig == nil {
		panic("Please call LoadConfig before GetConfig!")
	}
	return appConfig
}

type S3Config struct {
	Endpoint string
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	SslEnabled bool `mapstructure:"ssl"`
}

func (config *Config) GetS3Config() (*S3Config, error) {
	var s3Config S3Config
	if err := config.v.UnmarshalKey("s3", &s3Config); err != nil {
		return nil, err
	}

	return &s3Config, nil
}

type S3BucketsConfig struct {
	ImagesBucket string `mapstructure:"images"`
	RecordingsBucket string `mapstructure:"recordings"`
	LogsBucket string `mapstructure:"logs"`
}

func (config *Config) GetS3BucketsConfig() (*S3BucketsConfig, error) {
	var bucketsConfig S3BucketsConfig
	if err := config.v.UnmarshalKey("s3.buckets", &bucketsConfig); err != nil {
		return nil, err
	}

	return &bucketsConfig, nil
}