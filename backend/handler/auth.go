package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// --- Request / Response types ---

type tokenRequest struct {
	Token string `json:"token"`
}

type googleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Error         string `json:"error"`
	ErrorDesc     string `json:"error_description"`
}

type userResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type errorResponse struct {
	Message string `json:"message"`
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Message: message})
}

func verifyGoogleToken(idToken string) (*googleTokenInfo, error) {
	apiURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("gagal menghubungi Google: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca respons: %w", err)
	}

	var info googleTokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("gagal parse respons Google: %w", err)
	}

	if info.Error != "" {
		return nil, fmt.Errorf("%s", info.ErrorDesc)
	}
	if info.Email == "" {
		return nil, fmt.Errorf("token tidak mengandung email")
	}

	return &info, nil
}

// --- Handlers ---

// GoogleLogin — POST /auth/google
// Verifikasi Google ID Token dan kembalikan data user.
func GoogleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Request body tidak valid")
		return
	}
	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "Token tidak boleh kosong")
		return
	}

	info, err := verifyGoogleToken(req.Token)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Token tidak valid: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, userResponse{
		Name:    info.Name,
		Email:   info.Email,
		Picture: info.Picture,
	})
}

// GetMe — GET /auth/me
// Verifikasi token dari header Authorization dan kembalikan info user.
func GetMe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	authHeader := r.Header.Get("Authorization")
	if len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		writeError(w, http.StatusUnauthorized, "Authorization header tidak valid")
		return
	}

	info, err := verifyGoogleToken(authHeader[7:])
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Token tidak valid: "+err.Error())
		return
	}

	writeJSON(w, http.StatusOK, userResponse{
		Name:    info.Name,
		Email:   info.Email,
		Picture: info.Picture,
	})
}

// Logout — POST /auth/logout
// Server tidak menyimpan sesi; cukup konfirmasi logout berhasil.
func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method tidak diizinkan")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "Logout berhasil"})
}
