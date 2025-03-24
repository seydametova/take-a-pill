package models

import (
	"time"
)

// Структура для запроса на создание расписания
type ScheduleRequest struct {
	// ID пользователя
	UserID string `json:"user_id"`
	// Название лекарства
	MedicineName string `json:"medicine_name"`
	// Сколько раз в день принимать (от 1 до 24 раз в день)
	Frequency int `json:"frequency"`
	// Сколько дней принимать (0 - постоянный прием, >0 - количество дней)
	Duration int `json:"duration"`
}

// Структура для хранения расписания
type Schedule struct {
	// Уникальный ID расписания
	ID string `json:"id"`
	// ID пользователя
	UserID string `json:"user_id"`
	// Название лекарства
	MedicineName string `json:"medicine_name"`
	// Сколько раз в день принимать
	Frequency int `json:"frequency"`
	// Сколько дней принимать
	Duration int `json:"duration"`
	// Время создания расписания
	CreatedAt time.Time `json:"created_at"`
	// Рассчитанные времена приема
	TakingTimes []TakingTime `json:"taking_times"`
}

// Структура для хранения времени приема
type TakingTime struct {
	// Час приема (0-23)
	Hour int `json:"hour"`
	// Минута приема (0-59)
	Minute int `json:"minute"`
}

// RoundToQuarter округляет минуты до ближайших 15
func (t *TakingTime) RoundToQuarter() {
	minutes := t.Minute
	remainder := minutes % 15
	if remainder < 8 {
		t.Minute = minutes - remainder
	} else {
		t.Minute = minutes + (15 - remainder)
	}

	// Если после округления получилось 60 минут
	if t.Minute == 60 {
		t.Minute = 0
		t.Hour++
	}
}

// IsWithinDayHours проверяет, попадает ли время в диапазон 8:00-22:00
func (t *TakingTime) IsWithinDayHours() bool {
	return t.Hour >= 8 && t.Hour < 22
}

// CalculateTakingTimes рассчитывает времена приема для заданной частоты
func CalculateTakingTimes(frequency int) []TakingTime {
	dayStart := 8 // 8:00
	dayEnd := 22  // 22:00
	totalMinutes := (dayEnd - dayStart) * 60

	// Если прием 1 раз в день - фиксированное время 9:00
	if frequency == 1 {
		return []TakingTime{{Hour: 9, Minute: 0}}
	}

	// Рассчитываем интервал между приемами в минутах
	interval := totalMinutes / frequency

	var times []TakingTime
	for i := 0; i < frequency; i++ {
		// Рассчитываем время в минутах от начала дня
		totalMinutes := (dayStart * 60) + (i * interval)

		// Переводим в часы и минуты
		hour := totalMinutes / 60
		minute := totalMinutes % 60

		// Округляем минуты до ближайших 15
		minute = ((minute + 7) / 15) * 15

		// Корректируем час если минуты стали 60
		if minute == 60 {
			minute = 0
			hour++
		}

		// Проверяем что время попадает в диапазон 8:00-22:00
		if hour >= dayStart && hour < dayEnd {
			times = append(times, TakingTime{
				Hour:   hour,
				Minute: minute,
			})
		}
	}

	return times
}

// Структура для ответа со списком расписаний
type SchedulesResponse struct {
	ScheduleIDs []string `json:"schedule_ids"`
}

// Структура для ответа с деталями расписания
type ScheduleDetailsResponse struct {
	Schedule     *Schedule    `json:"schedule"`
	TakingsTimes []TakingTime `json:"takings_times"`
}

// Структура для ответа с ближайшими приёмами
type NextTakingsResponse struct {
	Takings []NextTaking `json:"takings"`
}

// Структура для информации о ближайшем приёме
type NextTaking struct {
	ScheduleID     string     `json:"schedule_id"`
	MedicineName   string     `json:"medicine_name"`
	NextTakingTime TakingTime `json:"next_taking_time"`
}
