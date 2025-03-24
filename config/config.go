package config

import "time"

// Config содержит все настройки сервиса
type Config struct {
	ServerPort       string
	NextTakingPeriod time.Duration
	DayStartHour     int
	DayEndHour       int
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		ServerPort:       ":8081",
		NextTakingPeriod: time.Hour,
		DayStartHour:     8,
		DayEndHour:       22,
	}
}
