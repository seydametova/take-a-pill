package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"take-a-pill/models"
	"take-a-pill/storage"

	"github.com/gorilla/mux"
)

// Структура для хранения данных сервера
type Server struct {
	// Хранилище для расписаний
	db *storage.MemoryStorage
	// Роутер
	router *mux.Router
}

// Создаем новый сервер
func NewServer() *Server {
	s := &Server{
		db:     storage.NewMemoryStorage(),
		router: mux.NewRouter(),
	}

	// Настраиваем маршруты
	s.routes()

	return s
}

// Настройка маршрутов
func (s *Server) routes() {
	// Добавляем логирование для всех запросов
	s.router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	s.router.HandleFunc("/schedule", s.createSchedule).Methods("POST")
	s.router.HandleFunc("/schedule", s.getScheduleDetails).Methods("GET")
	s.router.HandleFunc("/schedules", s.getSchedules).Methods("GET")
	s.router.HandleFunc("/next_takings", s.getNextTakings).Methods("GET")
}

// Обработчик для создания расписания
func (s *Server) createSchedule(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	var request models.ScheduleRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Ошибка при чтении данных", http.StatusBadRequest)
		return
	}

	// Создаем расписание
	schedule, err := s.db.CreateSchedule(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	response := fmt.Sprintf(`{"schedule_id": "%s"}`, schedule.ID)
	w.Write([]byte(response))
}

// Обработчик для получения списка расписаний пользователя
func (s *Server) getSchedules(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "не указан user_id", http.StatusBadRequest)
		return
	}

	// Получаем список расписаний
	scheduleIDs := s.db.GetSchedulesByUserID(userID)

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]string{
		"schedule_ids": scheduleIDs,
	})
}

// Обработчик для получения деталей расписания
func (s *Server) getScheduleDetails(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры запроса
	userID := r.URL.Query().Get("user_id")
	scheduleID := r.URL.Query().Get("schedule_id")

	log.Printf("Получен запрос на детали расписания: user_id=%s, schedule_id=%s", userID, scheduleID)

	if userID == "" {
		log.Println("Не указан user_id")
		http.Error(w, "не указан user_id", http.StatusBadRequest)
		return
	}

	if scheduleID == "" {
		log.Println("Не указан schedule_id")
		http.Error(w, "не указан schedule_id", http.StatusBadRequest)
		return
	}

	// Получаем расписание из хранилища
	schedule, err := s.db.GetScheduleByID(scheduleID)
	if err != nil {
		log.Printf("Ошибка при получении расписания: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Проверяем, что расписание принадлежит пользователю
	if schedule.UserID != userID {
		log.Printf("Расписание %s не принадлежит пользователю %s", scheduleID, userID)
		http.Error(w, "расписание не найдено", http.StatusNotFound)
		return
	}

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(schedule); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	log.Printf("Успешно отправлены детали расписания %s", scheduleID)
}

// Обработчик для получения следующих приемов
func (s *Server) getNextTakings(w http.ResponseWriter, r *http.Request) {
	// Получаем user_id из параметров запроса
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		log.Println("Не указан user_id")
		http.Error(w, "не указан user_id", http.StatusBadRequest)
		return
	}

	// Получаем следующие приемы
	nextTakings := s.db.GetNextTakings(userID)

	// Отправляем ответ
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string][]models.NextTaking{
		"takings": nextTakings,
	}); err != nil {
		log.Printf("Ошибка при отправке ответа: %v", err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	log.Printf("Успешно отправлены следующие приемы для пользователя %s", userID)
}

func main() {
	// Создаем сервер
	server := NewServer()

	// Запускаем сервер
	port := ":8081"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, server.router)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
