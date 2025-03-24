package storage

import (
	"fmt"
	"sync"
	"take-a-pill/models"
	"time"

	"take-a-pill/validation"

	"github.com/google/uuid"
)

// Структура для хранения расписаний в памяти
type MemoryStorage struct {
	// Карта для хранения расписаний, где ключ - это ID расписания
	schedules map[string]*models.Schedule
	// Мьютекс для безопасной работы с картой
	mu sync.RWMutex
}

// Создаем новое хранилище
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		schedules: make(map[string]*models.Schedule),
	}
}

// Создаем новое расписание
func (s *MemoryStorage) CreateSchedule(req *models.ScheduleRequest) (*models.Schedule, error) {
	// Валидация запроса
	if err := validation.ValidateScheduleRequest(req); err != nil {
		return nil, err
	}

	// Блокируем доступ к карте для записи
	s.mu.Lock()
	defer s.mu.Unlock()

	// Создаем новое расписание
	schedule := &models.Schedule{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		MedicineName: req.MedicineName,
		Frequency:    req.Frequency,
		Duration:     req.Duration,
		CreatedAt:    time.Now(),
		TakingTimes:  models.CalculateTakingTimes(req.Frequency),
	}

	// Сохраняем расписание в карту
	s.schedules[schedule.ID] = schedule

	return schedule, nil
}

// GetSchedulesByUserID возвращает список ID расписаний пользователя
func (s *MemoryStorage) GetSchedulesByUserID(userID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var scheduleIDs []string
	for id, schedule := range s.schedules {
		if schedule.UserID == userID {
			scheduleIDs = append(scheduleIDs, id)
		}
	}
	return scheduleIDs
}

// GetScheduleByID возвращает расписание по его ID
func (s *MemoryStorage) GetScheduleByID(scheduleID string) (*models.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if schedule, ok := s.schedules[scheduleID]; ok {
		return schedule, nil
	}
	return nil, fmt.Errorf("расписание не найдено")
}

// CalculateTakingTimes рассчитывает времена приёма лекарств
func (s *MemoryStorage) CalculateTakingTimes(schedule *models.Schedule) []models.TakingTime {
	dayStart := 8 // 8:00
	dayEnd := 22  // 22:00
	totalHours := dayEnd - dayStart

	// Рассчитываем интервал между приёмами в часах
	interval := totalHours / schedule.Frequency
	if interval < 1 {
		interval = 1 // Минимальный интервал - 1 час
	}

	var times []models.TakingTime
	for hour := dayStart; hour <= dayEnd-interval; hour += interval {
		// Округляем минуты до ближайших 15
		minute := 0
		times = append(times, models.TakingTime{
			Hour:   hour,
			Minute: minute,
		})
	}

	return times
}

// GetNextTakings возвращает ближайшие приёмы лекарств для пользователя
func (s *MemoryStorage) GetNextTakings(userID string) []models.NextTaking {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var nextTakings []models.NextTaking
	now := time.Now()
	currentHour := now.Hour()

	for id, schedule := range s.schedules {
		if schedule.UserID != userID {
			continue
		}

		times := s.CalculateTakingTimes(schedule)
		for _, t := range times {
			if t.Hour == currentHour || t.Hour == currentHour+1 {
				nextTakings = append(nextTakings, models.NextTaking{
					ScheduleID:     id,
					MedicineName:   schedule.MedicineName,
					NextTakingTime: t,
				})
			}
		}
	}

	return nextTakings
}
