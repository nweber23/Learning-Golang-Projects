package handlers

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"image-process-service/middleware"
	"image-process-service/models"
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
		middleware.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		middleware.JSONError(w, http.StatusBadRequest, "Username and password required")
		return
	}

	existing, _ := h.store.FindUserByUsername(req.Username)
	if existing != nil {
		middleware.JSONError(w, http.StatusConflict, "Username already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error processing password")
		return
	}

	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	if err := h.store.SaveUser(user); err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error saving user")
		return
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error generating token")
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
		middleware.JSONError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.store.FindUserByUsername(req.Username)
	if err != nil || user == nil {
		middleware.JSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		middleware.JSONError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := h.generateJWT(user.ID)
	if err != nil {
		middleware.JSONError(w, http.StatusInternalServerError, "Error generating token")
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
