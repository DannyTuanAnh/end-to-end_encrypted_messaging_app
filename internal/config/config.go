package config

import "time"

type Config struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

func NewConfig() *Config {
	return &Config{
		Port:            ":8080",
		ReadTimeout:     5 * time.Second,   // thời gian tối đa để đọc yêu cầu từ client
		WriteTimeout:    10 * time.Second,  // thời gian tối đa để gửi phản hồi cho một yêu cầu
		IdleTimeout:     120 * time.Second, // thời gian chờ tối đa cho một kết nối không hoạt động (giữ kết nối tối đa 2 phút)
		ShutdownTimeout: 5 * time.Second,   // thời gian tối đa để server hoàn thành các yêu cầu đang xử lý trước khi tắt
	}
}
