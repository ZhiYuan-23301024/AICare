package config

type Config struct {
	ServerPort string
	AIApiKey   string // 可从环境变量或配置文件映射
	AIBaseURL  string
}

func Load() *Config {
	// 可以从环境变量、配置文件等加载配置，这里使用默认值+环境变量
	return &Config{
		ServerPort: ":8080",
		AIApiKey:   "sk-bxbuwkbwyqooyvbhhsfpddbnfrydysffoxgulwsphhkjimxp",
		AIBaseURL:  "https://api.siliconflow.cn/v1/chat/completions",
	}
}
