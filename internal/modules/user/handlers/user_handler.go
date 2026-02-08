package handlers

import (
	"database/sql"
	"dental_clinic/internal/config"
	"dental_clinic/internal/modules/user/services"
	"dental_clinic/internal/modules/user/dto"
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
// Register godoc
// @Summary Register new user
// @Description Creates a new user account
// @Tags Users
// @Accept  json
// @Produce  json
// @Param request body dto.RegisterRequest true "User registration data"
// @Success 200 {object} dto.RegisterResponse
// @Failure 400 {object} dto.RegisterResponse
// @Router /api/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	response := dto.RegisterResponse{
		Success: "0",
		Message: "",
		User_id:  "",
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Message = "Invalid request body"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := h.service.Register(req)
	if err != nil {
		response.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	response.Success = "1"
	response.Message = "successfully created"
	response.User_id = user.Id.String()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}


// GetAllUsers godoc
// @Summary Get all users
// @Description Returns a list of all users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	_ = r
	tokenStr := getToken(r)
	users, err := h.service.GetAllUsers(tokenStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(services.ToUserResponseList(users))
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

	var req dto.RegisterRequest
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

// Login godoc
// @Summary Login
// @Description to login
// @Tags Users
// @Accept  json
// @Produce  json
// @Param request body dto.LoginRequest true "User login data"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.LoginResponse
// @Router /api/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	response := dto.LoginResponse{
		Success: "0",
		Token: "",
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	ip := r.RemoteAddr
	user, err := h.service.Login(req, ip)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	token, _ := utils.GenerateJWT(user.Id.String(), user.Email, user.Role, h.cfg.JWTSecret)
	response.Token = token
	response.Success = "1"
	// w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(map[string]string{"success": "1", "token": token})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdatePassword (w http.ResponseWriter, r *http.Request) {
	var req dto.UpdatePasswordRequest
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
	var req dto.UpdateEmailRequest
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


func (h *UserHandler) GetUserByIDAdmin(w http.ResponseWriter, r *http.Request) {
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