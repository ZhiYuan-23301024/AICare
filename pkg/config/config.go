// ...existing code...
package config

import "github.com/spf13/viper"

type Config struct {
	ServerPort string `mapstructure:"server_port"`
	AIApiKey   string `mapstructure:"ai_api_key"`
	AIBaseURL  string `mapstructure:"ai_base_url"`
}

func Load() *Config {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("configs")

	var cfg Config

	// 读取配置文件，若文件不存在或读取失败则返回零值（字段为空）
	if err := v.ReadInConfig(); err != nil {
		return &cfg
	}

	// 反序列化到结构体，未提供的字段则保持零值（空字符串）
	_ = v.Unmarshal(&cfg)
	return &cfg
}

// ...existing code...
