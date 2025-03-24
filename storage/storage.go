package storage

import (
	"fmt"
	"sync"
	"take-a-pill/models"

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
	// Проверяем, что user_id не пустой
	if req.UserID == "" {
		return nil, fmt.Errorf("не указан user_id")
	}

	// Блокируем доступ к карте для записи
	s.mu.Lock()
	// Разблокируем после завершения функции
	defer s.mu.Unlock()

	// Создаем новое расписание
	schedule := &models.Schedule{
		ID:           uuid.New().String(), // Генерируем уникальный ID
		UserID:       req.UserID,
		MedicineName: req.MedicineName,
		Frequency:    req.Frequency,
		Duration:     req.Duration,
	}

	// Сохраняем расписание в карту
	s.schedules[schedule.ID] = schedule

	// Возвращаем созданное расписание
	return schedule, nil
}
