package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"GoChat/domain"
)

type MessageService interface {
	GetAll() ([]domain.Message, error)
	Create(msg *domain.Message) (*domain.Message, error)
	Update(msg *domain.Message) (*domain.Message, error)
	Delete(id int) error
}

type MessageHandler struct {
	svc MessageService
}

func NewMessageHandler(svc MessageService) *MessageHandler {
	return &MessageHandler{
		svc: svc,
	}
}

func (h *MessageHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	msgs, err := h.svc.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msgs)
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req domain.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	msg := &domain.Message{
		Text:   req.Text,
		UserID: req.UserID,
	}

	createdMsg, err := h.svc.Create(msg)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain.Response{
		Status:  "success",
		Message: *createdMsg,
	})
}

func (h *MessageHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный или отсутствующий ID в URL"})
		return
	}

	var req domain.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	msg := &domain.Message{
		ID:   id,
		Text: req.Text,
	}

	updateMsg, err := h.svc.Update(msg)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{
		Status:  "success",
		Message: *updateMsg,
	})
}

func (h *MessageHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный или отсутствующий ID в URL"})
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.Response{
		Status: "success",
	})
}
