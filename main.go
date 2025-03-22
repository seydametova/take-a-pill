package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"take-a-pill/models"
	"take-a-pill/storage"
)

// Структура для хранения данных сервера
type Server struct {
	// Хранилище для расписаний
	db *storage.MemoryStorage
}

// Создаем новый сервер
func NewServer() *Server {
	return &Server{
		db: storage.NewMemoryStorage(),
	}
}

// Обработчик для создания расписания
func (s *Server) createSchedule(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != "POST" {
		http.Error(w, "Нужно использовать POST метод", http.StatusMethodNotAllowed)
		return
	}

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

func main() {
	// Создаем сервер
	server := NewServer()

	// Настраиваем маршруты
	http.HandleFunc("/schedule", server.createSchedule)

	// Запускаем сервер
	port := ":8081"
	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v\n", err)
	}
}
