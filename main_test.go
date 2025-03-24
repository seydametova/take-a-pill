package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"take-a-pill/models"
)

// Тест на создание расписания
func TestCreateSchedule(t *testing.T) {
	// Создаем сервер для тестов
	server := NewServer()

	// Тестовые данные
	data := models.ScheduleRequest{
		UserID:       "test123",
		MedicineName: "Аспирин",
		Frequency:    3,
		Duration:     7,
	}

	// Делаем JSON из данных
	jsonData, _ := json.Marshal(data)

	// Делаем тестовый запрос
	req := httptest.NewRequest("POST", "/schedule", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	server.createSchedule(w, req)

	// Проверяем что все ок
	if w.Code != http.StatusOK {
		t.Error("Должен быть код 200")
	}

	// Проверяем что есть schedule_id
	var result map[string]string
	json.NewDecoder(w.Body).Decode(&result)
	if result["schedule_id"] == "" {
		t.Error("Нет schedule_id в ответе")
	}
}

// Тест на пустой user_id
func TestEmptyUserID(t *testing.T) {
	server := NewServer()

	// Отправляем запрос без user_id
	data := models.ScheduleRequest{
		UserID:       "",
		MedicineName: "Аспирин",
		Frequency:    3,
		Duration:     7,
	}

	jsonData, _ := json.Marshal(data)
	req := httptest.NewRequest("POST", "/schedule", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.createSchedule(w, req)

	// Должна быть ошибка
	if w.Code != http.StatusInternalServerError {
		t.Error("Должен быть код 500")
	}
}

// Тест на GET запрос
func TestGetRequest(t *testing.T) {
	server := NewServer()

	// Делаем GET запрос вместо POST
	req := httptest.NewRequest("GET", "/schedule", nil)
	w := httptest.NewRecorder()
	server.createSchedule(w, req)

	// Должна быть ошибка
	if w.Code != http.StatusMethodNotAllowed {
		t.Error("Должен быть код 405")
	}
	if w.Body.String() != "Нужно использовать POST метод\n" {
		t.Error("Неправильная ошибка")
	}
}

// Тест на плохой JSON
func TestBadJSON(t *testing.T) {
	server := NewServer()

	// Отправляем плохой JSON
	req := httptest.NewRequest("POST", "/schedule", bytes.NewBufferString("{плохой json}"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	server.createSchedule(w, req)

	// Должна быть ошибка
	if w.Code != http.StatusBadRequest {
		t.Error("Должен быть код 400")
	}
	if w.Body.String() != "Ошибка при чтении данных\n" {
		t.Error("Неправильная ошибка")
	}
}
