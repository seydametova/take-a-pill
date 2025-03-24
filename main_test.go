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
	server.router.ServeHTTP(w, req)

	// Проверяем статус
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	// Проверяем что вернулся schedule_id
	var response map[string]string
	json.NewDecoder(w.Body).Decode(&response)
	if response["schedule_id"] == "" {
		t.Error("Не получен schedule_id")
	}
}

// Тест на пустой user_id
func TestCreateScheduleWithEmptyUserID(t *testing.T) {
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
	server.router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Ожидался статус 500, получен %d", w.Code)
	}
}

// Тест на GET запрос
func TestGetRequest(t *testing.T) {
	server := NewServer()

	// Делаем GET запрос вместо POST
	req := httptest.NewRequest("GET", "/schedule", nil)
	w := httptest.NewRecorder()

	// Используем router вместо прямого вызова обработчика
	server.router.ServeHTTP(w, req)

	// Должна быть ошибка
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался код 405, получен %d", w.Code)
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

func TestGetSchedules(t *testing.T) {
	server := NewServer()

	// Сначала создаем расписание
	createData := models.ScheduleRequest{
		UserID:       "test123",
		MedicineName: "Аспирин",
		Frequency:    3,
		Duration:     7,
	}
	jsonData, _ := json.Marshal(createData)
	createReq := httptest.NewRequest("POST", "/schedule", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, createReq)

	// Теперь получаем список расписаний
	req := httptest.NewRequest("GET", "/schedules?user_id=test123", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Проверяем статус
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	// Проверяем что в ответе есть список расписаний
	var response map[string][]string
	json.NewDecoder(w.Body).Decode(&response)
	if len(response["schedule_ids"]) == 0 {
		t.Error("Список расписаний пуст")
	}
}

func TestGetScheduleDetails(t *testing.T) {
	server := NewServer()

	// Сначала создаем расписание
	createData := models.ScheduleRequest{
		UserID:       "test123",
		MedicineName: "Аспирин",
		Frequency:    3,
		Duration:     7,
	}
	jsonData, _ := json.Marshal(createData)
	createReq := httptest.NewRequest("POST", "/schedule", bytes.NewBuffer(jsonData))
	createReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	server.router.ServeHTTP(w, createReq)

	// Получаем schedule_id из ответа
	var createResponse map[string]string
	json.NewDecoder(w.Body).Decode(&createResponse)
	scheduleID := createResponse["schedule_id"]

	// Получаем детали расписания
	req := httptest.NewRequest("GET", "/schedule?user_id=test123&schedule_id="+scheduleID, nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	// Проверяем статус
	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", w.Code)
	}

	// Проверяем данные расписания
	var schedule models.Schedule
	json.NewDecoder(w.Body).Decode(&schedule)
	if schedule.ID != scheduleID {
		t.Error("Неверный ID расписания")
	}
	if schedule.UserID != "test123" {
		t.Error("Неверный ID пользователя")
	}
	if schedule.MedicineName != "Аспирин" {
		t.Error("Неверное название лекарства")
	}
}
