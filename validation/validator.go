package validation

import (
	"errors"
	"take-a-pill/models"
)

// ValidateScheduleRequest проверяет корректность данных запроса на создание расписания
func ValidateScheduleRequest(req *models.ScheduleRequest) error {
	if req == nil {
		return errors.New("запрос не может быть пустым")
	}

	if req.UserID == "" {
		return errors.New("не указан идентификатор пользователя")
	}

	if req.MedicineName == "" {
		return errors.New("не указано название лекарства")
	}

	if req.Frequency < 1 || req.Frequency > 24 {
		return errors.New("частота приема должна быть от 1 до 24 раз в день")
	}

	if req.Duration < 0 {
		return errors.New("продолжительность лечения не может быть отрицательной")
	}

	return nil
}
