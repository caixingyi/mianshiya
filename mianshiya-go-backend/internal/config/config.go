package config

// Config 定义了应用程序的配置项，包括数据库连接配置和 Redis 连接配置，后续可以扩展为更多的配置项
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
	AI       AIConfig
	ES       ESConfig
}

// DatabaseConfig 定义了数据库连接的配置项
type DatabaseConfig struct {
	DBUser     string `mapstructure:"db_user" json:"db_user" yaml:"db_user"`
	DBPassword string `mapstructure:"db_password" json:"db_password" yaml:"db_password"`
	DBHost     string `mapstructure:"db_host" json:"db_host" yaml:"db_host"`
	DBPort     int    `mapstructure:"db_port" json:"db_port" yaml:"db_port"`
	DBName     string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
}

// RedisConfig 定义了 Redis 连接的配置项
type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

// Load 从配置文件或环境变量中加载配置项，目前直接返回硬编码的配置，后续可以扩展为从文件或环境变量加载
func Load() (*Config, error) {
	return &Config{
		Database: DatabaseConfig{
			DBUser:     "root",
			DBPassword: "12345678",
			DBHost:     "localhost",
			DBPort:     3306,
			DBName:     "mianshiya",
		},
		Redis: RedisConfig{
			Host:     "localhost",
			Port:     6379,
			Password: "",
			DB:       0,
		},
		AI: AIConfig{
			APIKey:  "ark-7a54e8b4-aeac-4aa4-9df3-8532b683b968-fa9cf",
			BaseURL: "https://ark.cn-beijing.volces.com/api/v3",
			Model:   "deepseek-v4-flash-260425",
		},
		ES: ESConfig{
			Addresses: []string{"http://127.0.0.1:9200"},
		},
	}, nil
}

// AIConfig 定义了 AI 调用的配置项（火山引擎 DeepSeek）
type AIConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

type ESConfig struct {
	Addresses []string // ES 地址列表，单机就是 ["http://localhost:9200"]
}
