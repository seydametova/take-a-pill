package models

// Структура для запроса на создание расписания
type ScheduleRequest struct {
	// ID пользователя
	UserID string `json:"user_id"`
	// Название лекарства
	MedicineName string `json:"medicine_name"`
	// Сколько раз в день принимать
	Frequency int `json:"frequency"`
	// Сколько дней принимать
	Duration int `json:"duration"`
}

// Структура для хранения расписания
type Schedule struct {
	// Уникальный ID расписания
	ID string
	// ID пользователя
	UserID string
	// Название лекарства
	MedicineName string
	// Сколько раз в день принимать
	Frequency int
	// Сколько дней принимать
	Duration int
}
