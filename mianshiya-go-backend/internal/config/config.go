package config

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	DBUser     string `mapstructure:"db_user" json:"db_user" yaml:"db_user"`
	DBPassword string `mapstructure:"db_password" json:"db_password" yaml:"db_password"`
	DBHost     string `mapstructure:"db_host" json:"db_host" yaml:"db_host"`
	DBPort     int    `mapstructure:"db_port" json:"db_port" yaml:"db_port"`
	DBName     string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
}

func Load() (*Config, error) {
	return &Config{
		Database: DatabaseConfig{
			DBUser:     "root",
			DBPassword: "12345678",
			DBHost:     "localhost",
			DBPort:     3306,
			DBName:     "mianshiya",
		},
	}, nil
}
