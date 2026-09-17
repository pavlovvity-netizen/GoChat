package service_test

import (
	"errors"
	"testing"

	"GoChat/domain"
	"GoChat/service"
)

// MockMessageRepository реализует интерфейс MessageRepository через функции-поля
type MockMessageRepository struct {
	GetAllFn func() ([]domain.Message, error)
	CreateFn func(msg *domain.Message) error
	UpdateFn func(msg *domain.Message) error
	DeleteFn func(id int) error
}

func (m *MockMessageRepository) GetAll() ([]domain.Message, error) {
	return m.GetAllFn()
}

func (m *MockMessageRepository) Create(msg *domain.Message) error {
	return m.CreateFn(msg)
}

func (m *MockMessageRepository) Update(msg *domain.Message) error {
	return m.UpdateFn(msg)
}

func (m *MockMessageRepository) Delete(id int) error {
	return m.DeleteFn(id)
}

// === ТЕСТЫ ДЛЯ GETALL ===

func TestMessageService_GetAll(t *testing.T) {
	t.Run("Успешное получение списка сообщений", func(t *testing.T) {
		expectedMsgs := []domain.Message{
			{ID: 1, Text: "Первое"},
			{ID: 2, Text: "Второе"},
		}

		mockRepo := &MockMessageRepository{
			GetAllFn: func() ([]domain.Message, error) {
				return expectedMsgs, nil
			},
		}

		svc := service.NewMessageService(mockRepo)
		result, err := svc.GetAll()

		if err != nil {
			t.Fatalf("Не ожидали ошибку, получили: %v", err)
		}

		if len(result) != 2 {
			t.Fatalf("Ожидали 2 сообщения, получили %d", len(result))
		}

		if result[0].Text != "Первое" || result[1].Text != "Второе" {
			t.Errorf("Получены некорректные данные: %+v", result)
		}
	})

	t.Run("Возврат пустого списка", func(t *testing.T) {
		mockRepo := &MockMessageRepository{
			GetAllFn: func() ([]domain.Message, error) {
				return []domain.Message{}, nil
			},
		}

		svc := service.NewMessageService(mockRepo)
		result, err := svc.GetAll()

		if err != nil {
			t.Fatalf("Не ожидали ошибку, получили: %v", err)
		}

		if len(result) != 0 {
			t.Errorf("Ожидали пустой слайс, получили длиною %d", len(result))
		}
	})

	t.Run("Ошибка репозитория (БД недоступна)", func(t *testing.T) {
		expectedErr := errors.New("ошибка подключения к БД")

		mockRepo := &MockMessageRepository{
			GetAllFn: func() ([]domain.Message, error) {
				return nil, expectedErr
			},
		}

		svc := service.NewMessageService(mockRepo)
		result, err := svc.GetAll()

		if err == nil {
			t.Fatal("Ожидали ошибку от репозитория, но получили nil")
		}

		if !errors.Is(err, expectedErr) {
			t.Errorf("Ожидали ошибку '%v', получили '%v'", expectedErr, err)
		}

		if result != nil {
			t.Errorf("При ошибке слайс должен быть nil, получили %+v", result)
		}
	})
}

// === ТЕСТЫ ДЛЯ CREATE ===

func TestMessageService_Create(t *testing.T) {
	t.Run("Успешное создание", func(t *testing.T) {
		mockRepo := &MockMessageRepository{
			CreateFn: func(msg *domain.Message) error {
				msg.ID = 1 // Имитируем генерацию ID базой данных
				return nil
			},
		}

		svc := service.NewMessageService(mockRepo)
		input := &domain.Message{Text: "Привет, мир!"}

		result, err := svc.Create(input)

		if err != nil {
			t.Fatalf("Не ожидали ошибку, получили: %v", err)
		}
		if result.ID != 1 {
			t.Errorf("Ожидали ID=1, получили %d", result.ID)
		}
		if result.PostTime.IsZero() {
			t.Error("PostTime должен быть проставлен сервисом")
		}
	})

	t.Run("Ошибка: Передан nil", func(t *testing.T) {
		mockRepo := &MockMessageRepository{}
		svc := service.NewMessageService(mockRepo)

		_, err := svc.Create(nil)

		if err == nil {
			t.Error("Ожидали ошибку для nil сообщения, но получили nil")
		}
	})

	t.Run("Ошибка: Пустой текст или только пробелы", func(t *testing.T) {
		mockRepo := &MockMessageRepository{}
		svc := service.NewMessageService(mockRepo)

		invalidInputs := []string{"", "   ", "\t\n"}

		for _, text := range invalidInputs {
			_, err := svc.Create(&domain.Message{Text: text})
			if err == nil {
				t.Errorf("Ожидали ошибку валидации для текста '%s', но ее нет", text)
			}
		}
	})
}

// === ТЕСТЫ ДЛЯ UPDATE ===

func TestMessageService_Update(t *testing.T) {
	t.Run("Успешное обновление", func(t *testing.T) {
		mockRepo := &MockMessageRepository{
			UpdateFn: func(msg *domain.Message) error {
				return nil
			},
		}

		svc := service.NewMessageService(mockRepo)
		input := &domain.Message{ID: 10, Text: "Обновленный текст"}

		result, err := svc.Update(input)

		if err != nil {
			t.Fatalf("Не ожидали ошибку, получили: %v", err)
		}
		if result.Text != "Обновленный текст" {
			t.Errorf("Ожидали 'Обновленный текст', получили '%s'", result.Text)
		}
		if result.PostTime.IsZero() {
			t.Error("PostTime должен быть обновлен")
		}
	})

	t.Run("Ошибка: Некорректный ID (<= 0)", func(t *testing.T) {
		mockRepo := &MockMessageRepository{}
		svc := service.NewMessageService(mockRepo)

		invalidIDs := []int{0, -1, -100}

		for _, id := range invalidIDs {
			_, err := svc.Update(&domain.Message{ID: id, Text: "Текст"})
			if err == nil {
				t.Errorf("Ожидали ошибку для ID=%d, но получили nil", id)
			}
		}
	})
}

// === ТЕСТЫ ДЛЯ DELETE ===

func TestMessageService_Delete(t *testing.T) {
	t.Run("Успешное удаление", func(t *testing.T) {
		mockRepo := &MockMessageRepository{
			DeleteFn: func(id int) error {
				if id != 5 {
					t.Errorf("Ожидали передачу ID=5 в репозиторий, получили %d", id)
				}
				return nil
			},
		}

		svc := service.NewMessageService(mockRepo)
		err := svc.Delete(5)

		if err != nil {
			t.Fatalf("Не ожидали ошибку, получили: %v", err)
		}
	})

	t.Run("Ошибка: Ненайденная запись в репозитории", func(t *testing.T) {
		mockRepo := &MockMessageRepository{
			DeleteFn: func(id int) error {
				return domain.ErrNotFound
			},
		}

		svc := service.NewMessageService(mockRepo)
		err := svc.Delete(999)

		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("Ожидали domain.ErrNotFound, получили %v", err)
		}
	})

	t.Run("Ошибка: Некорректный ID", func(t *testing.T) {
		mockRepo := &MockMessageRepository{}
		svc := service.NewMessageService(mockRepo)

		err := svc.Delete(0)

		if err == nil {
			t.Error("Ожидали ошибку для ID=0, но ее нет")
		}
	})
}
