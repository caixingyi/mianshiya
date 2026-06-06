package config

// Config 定义了应用程序的配置项，包括数据库连接配置和 Redis 连接配置，后续可以扩展为更多的配置项
type Config struct {
	Database DatabaseConfig
	Redis    RedisConfig
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
	}, nil
}
