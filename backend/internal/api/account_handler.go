package api

import (
	"encoding/json"
	"net/http"

	"github.com/helium-mmm/Publishier/internal/api/middleware"
	"github.com/helium-mmm/Publishier/internal/api/response"
	"github.com/helium-mmm/Publishier/internal/service"
)

type AccountHandler struct {
	service *service.AccountService
}

type connectTelegramRequest struct {
	BotToken string `json:"bot_token"`
	ChatID   string `json:"chat_id"`
}

type telegramStatusResponse struct {
	Connected bool   `json:"connected"`
	ChatID    string `json:"chat_id,omitempty"`
	Platform  string `json:"platform,omitempty"`
}

func NewAccountHandler(service *service.AccountService) *AccountHandler {
	return &AccountHandler{service: service}
}

func (h *AccountHandler) GetTelegram(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	status, err := h.service.GetTelegramStatus(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	resp := telegramStatusResponse{Connected: status.Connected}
	if status.Connected {
		resp.ChatID = status.ChatID
		resp.Platform = string(status.Platform)
	}

	response.JSON(w, http.StatusOK, resp)
}

func (h *AccountHandler) ConnectTelegram(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req connectTelegramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.ConnectTelegram(r.Context(), userID, req.BotToken, req.ChatID); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "connected"})
}
