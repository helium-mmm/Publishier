package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/helium-mmm/Publishier/internal/api/middleware"
	"github.com/helium-mmm/Publishier/internal/api/response"
	"github.com/helium-mmm/Publishier/internal/domain"
	"github.com/helium-mmm/Publishier/internal/service"
)

type PostHandler struct {
	service *service.PostService
}

type createPostRequest struct {
	Content string `json:"content"`
}

type postResponse struct {
	ID          string  `json:"id"`
	Content     string  `json:"content"`
	Status      string  `json:"status"`
	CreatedAt   string  `json:"created_at"`
	PublishedAt *string `json:"published_at,omitempty"`
}

func NewPostHandler(service *service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Content == "" {
		response.Error(w, http.StatusBadRequest, "content is required")
		return
	}

	post, err := h.service.Create(r.Context(), userID, req.Content)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "internal server error")
		return
	}

	response.JSON(w, http.StatusCreated, toPostResponse(post))
}

func (h *PostHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.service.Get(r.Context(), userID, postID)
	if err != nil {
		h.writePostError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toPostResponse(post))
}

func (h *PostHandler) Publish(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	postID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid post id")
		return
	}

	post, err := h.service.Publish(r.Context(), userID, postID)
	if err != nil {
		h.writePostError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, toPostResponse(post))
}

func (h *PostHandler) writePostError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		response.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidStatus):
		response.Error(w, http.StatusConflict, err.Error())
	case errors.Is(err, domain.ErrAccountNotFound):
		response.Error(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domain.ErrPublicationFailed):
		response.Error(w, http.StatusBadGateway, "failed to publish to telegram")
	default:
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func toPostResponse(post *domain.Post) postResponse {
	resp := postResponse{
		ID:        post.ID.String(),
		Content:   post.Content,
		Status:    string(post.Status),
		CreatedAt: post.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	if post.PublishedAt != nil {
		formatted := post.PublishedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.PublishedAt = &formatted
	}
	return resp
}
