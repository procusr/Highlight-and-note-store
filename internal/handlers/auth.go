package handlers

import (
	"encoding/json"
	"net/http"

	"mh-go/internal/middleware"
	"mh-go/internal/services"
	"mh-go/templates"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := templates.LoginPage().Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	// First authenticate the user
	user, err := h.authService.Login(email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to generate token"))
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	w.Header().Set("HX-Redirect", "/overview")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := templates.RegisterPage().Render(r.Context(), w); err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form data"))
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validation
	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email and password are required"))
		return
	}

	if len(password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Password must be at least 6 characters long"))
		return
	}

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Passwords do not match"))
		return
	}

	// Register user
	_, err := h.authService.Register(email, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	// Success - redirect to login
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) ForgotPasswordPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	templates.ForgotPasswordPage().Render(r.Context(), w)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Failed to parse form data</div>"))
		return
	}

	email := r.FormValue("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Email is required</div>"))
		return
	}

	// Get user by email
	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>No account found with this email</div>"))
		return
	}

	// Generate reset token
	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("<div class='text-red-600'>Failed to generate reset token</div>"))
		return
	}

	// TODO: Send reset email with token
	// For now, just redirect to reset password page with token
	w.Header().Set("HX-Redirect", "/reset-password?token="+token)
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	token := r.URL.Query().Get("token")
	templates.ResetPasswordPage(token).Render(r.Context(), w)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Failed to parse form data</div>"))
		return
	}

	token := r.FormValue("token")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Validation
	if password == "" || confirmPassword == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Password is required</div>"))
		return
	}

	if len(password) < 6 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Password must be at least 6 characters long</div>"))
		return
	}

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Passwords do not match</div>"))
		return
	}

	// Reset password
	err := h.authService.ResetPassword(token, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>" + err.Error() + "</div>"))
		return
	}

	// Success - redirect to login
	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})

	w.Header().Set("HX-Redirect", "/login")
	w.WriteHeader(http.StatusOK)
}

func getUserID(r *http.Request) uuid.UUID {
	if userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID); ok {
		return userID
	}
	return uuid.Nil
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
