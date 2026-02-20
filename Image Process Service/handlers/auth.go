package handlers

import (
	"net/http"

	"image-process-service/middleware"
	"image-process-service/models"
)

type AuthHandler struct {
	store *models.Store
}

func NewAuthHandler(store *models.Store) *AuthHandler {
	return &AuthHandler{store: store}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := DecodeJSON(r, &req); err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// TODO: Validate username and password
	// TODO: Hash password with bcrypt
	// TODO: Check if username already exists
	// TODO: Create user and save to store
	// TODO: Generate JWT token
	// TODO: Return AuthResponse
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := DecodeJSON(r, &req); err != nil {
		middleware.JSONError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	// TODO: Find user by username
	// TODO: Verify password hash with bcrypt
	// TODO: Generate JWT token
	// TODO: Return AuthResponse
}
