package handler

import (
	"encoding/json"
	"net/http"

	"GoChat/domain"
)

type UserService interface {
	RegisterUser(user *domain.RegisterRequest) (*domain.User, error)
	GetUserByName(req *domain.LoginRequest) (*domain.User, error)
}

type UserHandler struct {
	svc UserService
}

func NewUserHandler(svc UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var registerReq domain.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	newUser, err := h.svc.RegisterUser(&registerReq)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(domain.AuthResponse{
		Status: "success",
		User:   *newUser,
		Token:  "",
	})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var loginReq domain.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Некорректный формат JSON"})
		return
	}

	user, err := h.svc.GetUserByName(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domain.AuthResponse{
		Status: "success",
		User:   *user,
		Token:  "",
	})
}
