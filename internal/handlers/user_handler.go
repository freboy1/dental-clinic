package handlers

import (
	"database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/services"
	"dental_clinic/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service *services.UserService
	cfg config.Config
}

func NewUserHandler(s *services.UserService, cfg config.Config) *UserHandler {
	return &UserHandler{
		service: s,
		cfg: cfg,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	_ = r
	users, err := h.service.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req services.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.UpdateUser(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	tokenStr := getToken(r)
	err := h.service.DeleteUser(id, tokenStr)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func (h *UserHandler) VerifyAccountByLink(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	err := h.service.VerifyUserEmail(token)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req services.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	ip := r.RemoteAddr
	user, err := h.service.Login(req, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	token, _ := utils.GenerateJWT(user.Id.String(), user.Email, user.Role, h.cfg.JWTSecret)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) UpdatePassword (w http.ResponseWriter, r *http.Request) {
	var req services.UpdatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenStr := getToken(r)
	err := h.service.UpdatePassword(tokenStr, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"updated": "successfully"})

}

func getToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	parts := strings.Split(authHeader, " ")
	tokenStr := parts[1]
	return tokenStr
}

func (h *UserHandler) UpdateEmail (w http.ResponseWriter, r *http.Request) {
	var req services.UpdateEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	tokenStr := getToken(r)
	err := h.service.UpdateEmail(tokenStr, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
}

func (h *UserHandler) VerifyNewEmail(w http.ResponseWriter, r *http.Request) {
    token := r.URL.Query().Get("token")
    if token == "" {
        http.Error(w, "Token missing", http.StatusBadRequest)
        return
    }

    err := h.service.VerifyEmailToken(token)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"updated": "successfully"})
}
