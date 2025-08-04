package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mh-go/internal/middleware"
	"mh-go/internal/services"

	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	// Render login page template
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login - MH-Go</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
    <div class="max-w-md w-full space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Sign in to your account
            </h2>
        </div>
        <form class="mt-8 space-y-6" hx-post="/login" hx-target="#error" hx-swap="innerHTML">
            <div class="rounded-md shadow-sm -space-y-px">
                <div>
                    <input name="email" type="email" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Email address">
                </div>
                <div>
                    <input name="password" type="password" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Password">
                </div>
            </div>
            <div id="error" class="text-red-600 text-sm"></div>
            <div>
                <button type="submit" 
                        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Sign in
                </button>
            </div>
            <div class="text-center">
                <a href="/register" class="text-indigo-600 hover:text-indigo-500">
                    Don't have an account? Sign up
                </a>
            </div>
        </form>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	token, err := h.authService.Login(email, password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
		return
	}

	// Set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		MaxAge:   7 * 24 * 60 * 60, // 7 days
	})

	w.Header().Set("HX-Redirect", "/dashboard")
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	// Render register page template
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Register - MH-Go</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
    <div class="max-w-md w-full space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Create your account
            </h2>
        </div>
        <form class="mt-8 space-y-6" hx-post="/register" hx-target="#error" hx-swap="innerHTML" hx-indicator="#loading">
            <div class="rounded-md shadow-sm -space-y-px">
                <div>
                    <input name="email" type="email" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Email address">
                </div>
                <div>
                    <input name="password" type="password" required minlength="6"
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Password (min 6 characters)">
                </div>
                <div>
                    <input name="confirm_password" type="password" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Confirm Password">
                </div>
            </div>
            <div id="error" class="text-red-600 text-sm"></div>
            <div id="loading" class="htmx-indicator text-blue-600 text-sm">Processing...</div>
            <div>
                <button type="submit" 
                        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Sign up
                </button>
            </div>
            <div class="text-center">
                <a href="/login" class="text-indigo-600 hover:text-indigo-500">
                    Already have an account? Sign in
                </a>
            </div>
        </form>
    </div>
</body>
</html>`
	w.Write([]byte(html))
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
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Forgot Password - MH-Go</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
    <div class="max-w-md w-full space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Reset your password
            </h2>
        </div>
        <form class="mt-8 space-y-6" hx-post="/forgot-password" hx-target="#message" hx-swap="innerHTML">
            <div>
                <input name="email" type="email" required 
                       class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                       placeholder="Email address">
            </div>
            <div id="message" class="text-sm"></div>
            <div>
                <button type="submit" 
                        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Send Reset Link
                </button>
            </div>
            <div class="text-center">
                <a href="/login" class="text-indigo-600 hover:text-indigo-500">
                    Back to sign in
                </a>
            </div>
        </form>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")

	_, err := h.authService.GenerateResetToken(email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("User not found"))
		return
	}

	w.Write([]byte("<div class='text-green-600'>Password reset link sent to your email</div>"))
}

func (h *AuthHandler) ResetPasswordPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Password - MH-Go</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <script src="https://unpkg.com/htmx.org@1.9.6"></script>
</head>
<body class="bg-gray-50 min-h-screen flex items-center justify-center">
    <div class="max-w-md w-full space-y-8">
        <div>
            <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
                Reset your password
            </h2>
        </div>
        <form class="mt-8 space-y-6" hx-post="/reset-password" hx-target="#message" hx-swap="innerHTML">
            <input type="hidden" name="token" value="` + r.URL.Query().Get("token") + `">
            <div class="rounded-md shadow-sm -space-y-px">
                <div>
                    <input name="password" type="password" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-t-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="New Password">
                </div>
                <div>
                    <input name="confirm_password" type="password" required 
                           class="relative block w-full px-3 py-2 border border-gray-300 placeholder-gray-500 text-gray-900 rounded-b-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                           placeholder="Confirm New Password">
                </div>
            </div>
            <div id="message" class="text-sm"></div>
            <div>
                <button type="submit" 
                        class="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                    Reset Password
                </button>
            </div>
        </form>
    </div>
</body>
</html>`
	w.Write([]byte(html))
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("token")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if password != confirmPassword {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>Passwords do not match</div>"))
		return
	}

	err := h.authService.ResetPassword(token, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("<div class='text-red-600'>" + err.Error() + "</div>"))
		return
	}

	w.Write([]byte("<div class='text-green-600'>Password reset successfully</div>"))
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
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
