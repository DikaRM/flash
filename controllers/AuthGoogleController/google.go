package AuthGoogleController

import (
	"database/sql"
	"encoding/json"
	"flash/config"
	"flash/entities"
	"log"
	"net/http"
	"net/url"
)

type Auth struct {
	ClientID     string
	ClientSecret string
	db           *sql.DB
}

func NewAuthController(cfg *config.Config, database *sql.DB) *Auth {
	return &Auth{
		ClientID:     cfg.GoogleClient,
		ClientSecret: cfg.GoogleSecret,
		db:           database,
	}
}
func (c *Auth) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	// 1. TUKAR CODE JADI ACCESS TOKEN
	tokenURL := "https://oauth2.googleapis.com/token"
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	data.Set("redirect_uri", "http://localhost:8000/auth/google/callback")
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		http.Error(w, "Gagal tukar token", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	// 2. AMBIL DATA USER DARI GOOGLE PAKAI ACCESS TOKEN
	userInfoURL := "https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + tokenResp.AccessToken
	userResp, _ := http.Get(userInfoURL)
	defer userResp.Body.Close()

	var userInfo struct {
		ID       string `json:"id"`    // Google ID
		Email    string `json:"email"` // Email user
		Username string `json:"name"`  // Nama user
	}
	json.NewDecoder(userResp.Body).Decode(&userInfo)

	// 3. VALIDASI dengan DATABASE (pakai email & google_id dari JSON)
	var user entities.Users
	err = c.db.QueryRow(`
        SELECT id,username,email, google_id  FROM users
        WHERE email = ? OR google_id = ?
    `, userInfo.Email, userInfo.ID).Scan(&user.ID, &user.Username, &user.Email, &user.GoogleID)

	if err == sql.ErrNoRows {
		// REGISTER: user belum ada, simpan dari data JSON
		result, err := c.db.Exec(`
            INSERT INTO users (username,email, google_id, created_at)
            VALUES (?, ?, ?, datetime('now'))
        `, userInfo.Username, userInfo.Email, userInfo.ID)
		if err != nil {
			http.Error(w, "erroror Insert"+err.Error(), http.StatusInternalServerError)
			return
		}
		userID, err := result.LastInsertId()
		user.ID = userID
		if err != nil {
			http.Error(w, "Error Id "+err.Error(), http.StatusInternalServerError)
		}

	} else if err != nil {
		http.Error(w, "Error database"+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. UPDATE google_id jika perlu
	if user.GoogleID == nil || *user.GoogleID != userInfo.ID {
		c.db.Exec(`UPDATE users SET google_id = ? WHERE email = ?`, userInfo.ID, userInfo.Email)
	}

	// 5. SIMPAN SESSION & REDIRECT
	session, err := config.GetSession(r)
	if err != nil {
		log.Println("Gagal Nyimpan Session")
		return
	}
	session.Values["username"] = user.Username
	session.Values["userEmail"] = user.Email
	session.Values["logged_in"] = true
	if err := session.Save(r, w); err != nil {
		log.Println("Gagal Nyimpannny", err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	return
}

func (c *Auth) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	// Build URL manual
	url := "https://accounts.google.com/o/oauth2/v2/auth?" +
		"client_id=" + c.ClientID +
		"&redirect_uri=http://localhost:8000/auth/google/callback" +
		"&response_type=code" +
		"&scope=email profile" +
		"&access_type=online" +
		"&state=random-state-token"

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
