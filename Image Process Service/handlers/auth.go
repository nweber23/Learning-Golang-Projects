package handlers

import (
	"net/http"
	"time"
	"image-process-service/middleware"
	"image-process-service/models"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	store     *models.Store
	jwtSecret string
}

func NewAuthHandler(store *models.Store, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := DecodeJSON(r, &req); err != nil {
		middleware.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		middleware.JSONError(w, "Username and password required", http.StatusBadRequest)
		return
	}

	existing, _ := h.store.FindUserByUsername(req.Username)
	if existing != nil {
		middleware.JSONError(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.JSONError(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := h.store.SaveUser(user); err != nil {
		middleware.JSONError(w, "Error saving user", http.StatusInternalServerError)
		return
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		middleware.JSONError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		JWT:      token,
	}

	WriteJSON(w, http.StatusCreated, response)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := DecodeJSON(r, &req); err != nil {
		middleware.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.store.FindUserByUsername(req.Username)
	if err != nil || user == nil {
		middleware.JSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		middleware.JSONError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		middleware.JSONError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := models.AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		JWT:      token,
	}

	WriteJSON(w, http.StatusOK, response)
}

func (h *AuthHandler) generateJWT(userID string) (string, error) {
	claims := middleware.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.jwtSecret))
}
