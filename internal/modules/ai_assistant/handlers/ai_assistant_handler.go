package handlers

import (
	"encoding/json"
	"net/http"

	"dental_clinic/internal/config"
	"dental_clinic/internal/middleware"
	"dental_clinic/internal/modules/ai_assistant/dto"
	"dental_clinic/internal/modules/ai_assistant/services"
	"dental_clinic/internal/utils"
)

type AIAssistantHandler struct {
	service *services.AIAssistantService
	cfg     config.Config
}

func NewAIAssistantHandler(s *services.AIAssistantService, cfg config.Config) *AIAssistantHandler {
	return &AIAssistantHandler{
		service: s,
		cfg:     cfg,
	}
}

// Chat godoc
// @Summary AI assistant chat
// @Description Processes a user chat message, updates booking state, and creates an appointment when state is complete
// @Tags AI Assistant
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.ChatRequest true "Chat message"
// @Success 200 {object} dto.ChatResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/ai/chat [post]
func (h *AIAssistantHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req dto.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userID, err := middleware.GetUserID(r, h.cfg.JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	tokenStr := utils.GetToken(r)
	response, err := h.service.ProcessMessage(userID, tokenStr, req, r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

// Reset godoc
// @Summary Reset AI assistant booking flow
// @Description Clears the current booking state and starts a new chat session
// @Tags AI Assistant
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.ChatResponse
// @Failure 401 {object} map[string]string
// @Router /api/ai/chat/reset [post]
func (h *AIAssistantHandler) Reset(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserID(r, h.cfg.JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response, err := h.service.ResetBooking(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
