package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"GoChat/domain"
	"GoChat/handler"
)

type MockMessageService struct {
	GetAllFn func() ([]domain.Message, error)
	CreateFn func(msg *domain.Message) (*domain.Message, error)
	UpdateFn func(msg *domain.Message) (*domain.Message, error)
	DeleteFn func(id int) error
}

func (m *MockMessageService) GetAll() ([]domain.Message, error) {
	return m.GetAllFn()
}

func (m *MockMessageService) Create(msg *domain.Message) (*domain.Message, error) {
	return m.CreateFn(msg)
}

func (m *MockMessageService) Update(msg *domain.Message) (*domain.Message, error) {
	return m.UpdateFn(msg)
}

func (m *MockMessageService) Delete(id int) error {
	return m.DeleteFn(id)
}

func TestMessageHandler_GetAll(t *testing.T) {
	t.Run("Успешное получение списка сообщений (200 OK)", func(t *testing.T) {
		expectedMessages := []domain.Message{
			{ID: 1, Text: "Первое сообщение"},
			{ID: 2, Text: "Второе сообщение"},
		}

		mockSvc := &MockMessageService{
			GetAllFn: func() ([]domain.Message, error) {
				return expectedMessages, nil
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/messages", nil)
		rec := httptest.NewRecorder()

		h.GetAll(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Ожидали статус 200, получили %d", rec.Code)
		}

		var actualMessages []domain.Message
		if err := json.Unmarshal(rec.Body.Bytes(), &actualMessages); err != nil {
			t.Fatalf("Не удалось распарсить ответ хэндлера: %v", err)
		}

		if len(actualMessages) != 2 {
			t.Errorf("Ожидали 2 сообщения, получили %d", len(actualMessages))
		}
	})

	t.Run("Внутренняя ошибка сервера (500 Internal Server Error)", func(t *testing.T) {
		mockSvc := &MockMessageService{
			GetAllFn: func() ([]domain.Message, error) {
				return nil, errors.New("ошибка базы данных")
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/messages", nil)
		rec := httptest.NewRecorder()

		h.GetAll(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("Ожидали статус 500, получили %d", rec.Code)
		}
	})
}

func TestMessageHandler_Create(t *testing.T) {
	t.Run("Успешное создание (201 Created)", func(t *testing.T) {
		mockSvc := &MockMessageService{
			CreateFn: func(msg *domain.Message) (*domain.Message, error) {
				return &domain.Message{
					ID:   1,
					Text: msg.Text,
				}, nil
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		reqBody := `{"text": "Привет, мир!"}`
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Ожидали статус 201, получили %d", rec.Code)
		}

		var resp domain.Response
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Не удалось распарсить ответ хэндлера: %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Ожидали status='success', получили '%s'", resp.Status)
		}

		if resp.Message.ID != 1 || resp.Message.Text != "Привет, мир!" {
			t.Errorf("Некорректное сообщение в ответе: %+v", resp.Message)
		}
	})

	t.Run("Некорректный JSON (400 Bad Request)", func(t *testing.T) {
		mockSvc := &MockMessageService{}
		h := handler.NewMessageHandler(mockSvc)

		invalidJSON := `{"text": "Привет`
		req := httptest.NewRequest(http.MethodPost, "/messages", bytes.NewBufferString(invalidJSON))
		rec := httptest.NewRecorder()

		h.Create(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Ожидали статус 400, получили %d", rec.Code)
		}
	})
}

func TestMessageHandler_Update(t *testing.T) {
	t.Run("Успешное обновление (200 OK)", func(t *testing.T) {
		mockSvc := &MockMessageService{
			UpdateFn: func(msg *domain.Message) (*domain.Message, error) {
				if msg.ID != 1 {
					t.Errorf("Ожидали ID = 1, получили %d", msg.ID)
				}
				return &domain.Message{
					ID:   msg.ID,
					Text: msg.Text,
				}, nil
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		reqBody := `{"text": "Обновленный текст"}`
		req := httptest.NewRequest(http.MethodPut, "/messages/1", bytes.NewBufferString(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Ожидали статус 200, получили %d", rec.Code)
		}

		// Декодируем в domain.Response вместо domain.Message
		var resp domain.Response
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Не удалось распарсить ответ: %v", err)
		}

		if resp.Status != "success" {
			t.Errorf("Ожидали status='success', получили '%s'", resp.Status)
		}

		if resp.Message.Text != "Обновленный текст" {
			t.Errorf("Ожидали текст 'Обновленный текст', получили '%s'", resp.Message.Text)
		}
	})

	t.Run("Запись не найдена (404 Not Found)", func(t *testing.T) {
		mockSvc := &MockMessageService{
			UpdateFn: func(msg *domain.Message) (*domain.Message, error) {
				return nil, domain.ErrNotFound
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		reqBody := `{"text": "Неизвестное сообщение"}`
		req := httptest.NewRequest(http.MethodPut, "/messages/999", bytes.NewBufferString(reqBody))
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Ожидали статус 404, получили %d", rec.Code)
		}
	})

	t.Run("Некорректный JSON (400 Bad Request)", func(t *testing.T) {
		mockSvc := &MockMessageService{}
		h := handler.NewMessageHandler(mockSvc)

		invalidJSON := `{"text": "Сломанный JSON`
		req := httptest.NewRequest(http.MethodPut, "/messages/1", bytes.NewBufferString(invalidJSON))
		req.SetPathValue("id", "1")
		rec := httptest.NewRecorder()

		h.Update(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("Ожидали статус 400, получили %d", rec.Code)
		}
	})
}

func TestMessageHandler_Delete(t *testing.T) {
	t.Run("Успешное удаление (200 OK / 204 No Content)", func(t *testing.T) {
		// Arrange: мок возвращает nil (нет ошибки)
		mockSvc := &MockMessageService{
			DeleteFn: func(id int) error {
				if id != 1 {
					t.Errorf("Ожидали ID = 1, получили %d", id)
				}
				return nil
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/messages/1", nil)
		req.SetPathValue("id", "1") // для встроенного роутера Go 1.22+
		rec := httptest.NewRecorder()

		// Act
		h.Delete(rec, req)

		// Assert: ждем 200 OK или 204 No Content (в зависимости от твоей реализации)
		if rec.Code != http.StatusOK && rec.Code != http.StatusNoContent {
			t.Errorf("Ожидали статус 200 или 204, получили %d", rec.Code)
		}
	})

	t.Run("Запись не найдена (404 Not Found)", func(t *testing.T) {
		// Arrange: мок возвращает ошибку domain.ErrNotFound
		mockSvc := &MockMessageService{
			DeleteFn: func(id int) error {
				return domain.ErrNotFound
			},
		}

		h := handler.NewMessageHandler(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/messages/999", nil)
		req.SetPathValue("id", "999")
		rec := httptest.NewRecorder()

		// Act
		h.Delete(rec, req)

		// Assert
		if rec.Code != http.StatusNotFound {
			t.Errorf("Ожидали статус 404, получили %d", rec.Code)
		}
	})
}
