package config

import (
	"fmt"
	"time"

	"github.com/DannyTuanAnh/end-to-end_encrypted_messaging_app/internal/utils"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type ServiceConfig struct {
	AuthServiceAddr string
}

type Config struct {
	DB      DatabaseConfig
	Server  ServerConfig
	Service ServiceConfig
}

func NewConfigServer() *Config {
	return &Config{
		Server: ServerConfig{
			Port:            utils.GetEnv("SV_PORT", "8080"),
			ReadTimeout:     utils.GetEnvTime("SV_ReadTimeout", 5) * time.Second,     // thời gian tối đa để đọc yêu cầu từ client
			WriteTimeout:    utils.GetEnvTime("SV_WriteTimeout", 10) * time.Second,   // thời gian tối đa để gửi phản hồi cho một yêu cầu
			IdleTimeout:     utils.GetEnvTime("SV_IdleTimeout", 120) * time.Second,   // thời gian chờ tối đa cho một kết nối không hoạt động (giữ kết nối tối đa 2 phút)
			ShutdownTimeout: utils.GetEnvTime("SV_ShutdownTimeout", 5) * time.Second, // thời gian tối đa để server hoàn thành các yêu cầu đang xử lý trước khi tắt
		},

		Service: ServiceConfig{
			AuthServiceAddr: utils.GetEnv("AUTH_SERVICE_ADDR", ":50051"),
		},
	}
}

func NewConfigDB() *Config {
	return &Config{
		DB: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			DBName:   utils.GetEnv("DB_NAME", "myapp"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}

func (c *Config) DB_DNS() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.DB.User, c.DB.Password, c.DB.Host, c.DB.Port, c.DB.DBName, c.DB.SSLMode)
}
