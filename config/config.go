package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	}
	Supabase struct {
		URL    string `mapstructure:"url"`
		APIKey string `mapstructure:"api_key"`
	}
}

var AppConfig Config

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// 设置默认值
	viper.SetDefault("server.port", "8080")

	// 环境变量映射
	viper.BindEnv("supabase.url", "SUPABASE_URL")
	viper.BindEnv("supabase.api_key", "SUPABASE_KEY")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果找不到配置文件，尝试从环境变量读取
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(&AppConfig)
}
