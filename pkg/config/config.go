package config

import "os"

type Config struct {
	ServerPort string
	AIApiKey   string `mapstructure:"AI_API_KEY"` // 可从环境变量或配置文件映射
	AIBaseURL  string `mapstructure:"AI_BASE_URL"`
}

func Load() *Config {
	// 可以从环境变量、配置文件等加载配置，这里使用默认值+环境变量
	return &Config{
		ServerPort: ":8080",
		AIApiKey:   os.Getenv("AI_API_KEY"),
		AIBaseURL:  "https://api.openai.com/v1/chat/completions",
	}
}
