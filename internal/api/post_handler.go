package api

import (
	"net/http"

	"github.com/helium-mmm/Publishier/internal/service"
)

type PostHandler struct {
	service *service.PostService
}

type PostRequest struct {
	Content string `json:"content"`
}

func NewPostHandler(service *service.PostService) *PostHandler { 
	return &PostHandler{service: service}
}

func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) { 
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", )
	}
}

func (h *PostHandler) Publish( ) { 

}

func (h *PostHandler) Get( ) { 

}