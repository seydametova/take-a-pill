package storage

import (
	"fmt"
	"log"
	"sort"
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
	log.Printf("=== Начало GetNextTakings ===")
	log.Printf("Текущее время: %v", now.Format("15:04:05"))

	for id, schedule := range s.schedules {
		log.Printf("\nПроверка расписания %s:", id)
		log.Printf("- Пользователь: %s", schedule.UserID)
		log.Printf("- Лекарство: %s", schedule.MedicineName)
		log.Printf("- Частота: %d раз в день", schedule.Frequency)
		log.Printf("- Длительность: %d дней", schedule.Duration)
		log.Printf("- Создано: %v", schedule.CreatedAt.Format("2006-01-02 15:04:05"))
		log.Printf("- Времена приема: %v", schedule.TakingTimes)

		if schedule.UserID != userID {
			log.Printf("Пропускаем: расписание принадлежит другому пользователю")
			continue
		}

		// Проверяем, не истек ли срок действия расписания
		endDate := schedule.CreatedAt.AddDate(0, 0, schedule.Duration)
		if now.After(endDate) {
			log.Printf("Пропускаем: расписание истекло (конец: %v)", endDate.Format("2006-01-02 15:04:05"))
			continue
		}

		for _, t := range schedule.TakingTimes {
			// Создаем время приема на сегодня
			takingTime := time.Date(now.Year(), now.Month(), now.Day(), t.Hour, t.Minute, 0, 0, now.Location())

			log.Printf("\nПроверка времени приема %v:", takingTime.Format("15:04"))

			// Если время уже прошло, пропускаем
			if now.After(takingTime) {
				log.Printf("- Пропускаем: время уже прошло")
				continue
			}

			log.Printf("- Добавляем в список приемов")
			nextTakings = append(nextTakings, models.NextTaking{
				ScheduleID:     id,
				MedicineName:   schedule.MedicineName,
				NextTakingTime: t,
			})
		}
	}

	// Сортируем по времени приема
	sort.Slice(nextTakings, func(i, j int) bool {
		timeI := time.Date(now.Year(), now.Month(), now.Day(),
			nextTakings[i].NextTakingTime.Hour,
			nextTakings[i].NextTakingTime.Minute, 0, 0, now.Location())
		timeJ := time.Date(now.Year(), now.Month(), now.Day(),
			nextTakings[j].NextTakingTime.Hour,
			nextTakings[j].NextTakingTime.Minute, 0, 0, now.Location())
		return timeI.Before(timeJ)
	})

	log.Printf("\n=== Результаты ===")
	log.Printf("Всего найдено приемов: %d", len(nextTakings))
	for _, t := range nextTakings {
		log.Printf("- %s в %02d:%02d", t.MedicineName, t.NextTakingTime.Hour, t.NextTakingTime.Minute)
	}
	log.Printf("=== Конец GetNextTakings ===\n")

	return nextTakings
}
