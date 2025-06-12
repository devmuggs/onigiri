package auth

import (
	"encoding/json"
	"net/http"

	"github.com/devmuggs/onigiri/server/internal/features/users"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type AuthHandler struct {
	// service     AuthService
	userService users.Service
	logger      *zap.Logger
}

func NewHandler(userService users.Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{logger: logger, userService: userService}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/sign-up", h.SignUp)
	r.Post("/login", h.Login)
	r.Post("/login-out", h.SignUp)
	r.Get("/me", h.SignUp)
	return r
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	var input users.CreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	existingUser, err := h.userService.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}
	if existingUser != nil {
		http.Error(w, "email already in use.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(input.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := &users.User{
		DisplayName: input.DisplayName,
		Email:       input.Email,
	}

	input.Password = hashedPassword

	if err := h.userService.CreateUser(r.Context(), &input); err != nil {
		h.logger.Error("createUser error", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	existingUser, err := h.userService.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		http.Error(w, "invalid email", http.StatusBadRequest)
		return
	}
	if existingUser == nil {
		http.Error(w, "incorrect email or password.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(existingUser.HashedPassword)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	isCorrectPassword := CheckPassword(existingUser.HashedPassword, hashedPassword)
	if !isCorrectPassword {
		http.Error(w, "incorrect email or password.", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingUser.User)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("logout"))
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("me"))
}
