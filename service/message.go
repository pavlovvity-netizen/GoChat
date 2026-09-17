package service

import (
	"errors"
	"strings"
	"time"

	"GoChat/domain"
)

type MessageRepository interface {
	GetAll() ([]domain.Message, error)
	Create(msg *domain.Message) error
	Update(msg *domain.Message) error
	Delete(id int) error
}

type MessageService struct {
	repo MessageRepository
}

func NewMessageService(repo MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) GetAll() ([]domain.Message, error) {
	msgs, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (s *MessageService) Create(msg *domain.Message) (*domain.Message, error) {
	if msg == nil {
		return nil, errors.New("Сообщение не может быть nil")
	}
	if strings.TrimSpace(msg.Text) == "" {
		return nil, errors.New("Сообщение не может быть пустым")
	}

	msg.PostTime = time.Now()

	err := s.repo.Create(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) Update(msg *domain.Message) (*domain.Message, error) {
	if msg == nil {
		return nil, errors.New("Сообщение не может быть nil")
	}
	if msg.ID <= 0 {
		return nil, errors.New("Некорректный id")
	}
	if strings.TrimSpace(msg.Text) == "" {
		return nil, errors.New("Сообщение не может быть пустым")
	}

	err := s.repo.Update(msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s *MessageService) Delete(id int) error {
	if id <= 0 {
		return errors.New("Некорректный id")
	}

	return s.repo.Delete(id)
}
