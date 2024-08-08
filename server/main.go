package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ErrorDetails struct {
	ID            int            `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Domain        string         `gorm:"not null" json:"domain"`
	Error         string         `gorm:"not null" json:"errortext"`
	URL           string         `gorm:"not null" json:"url"`
	Line          string         `gorm:"not null" json:"line"`
	Datetime      string         `gorm:"not null" json:"datetime"`
	OS            string         `gorm:"not null" json:"os"`
	Browser       string         `gorm:"not null" json:"browser"`
	WebPropertyID int            `gorm:"not null" json:"web_property_id"`
}

var APP_CONFIG map[string]string
var GlobalMailer Mailer

func main() {
	envFile, _ := godotenv.Read(".env")

	// get the values from the environment variables from .env file
	APP_CONFIG = envFile
	smtpUsername := envFile["SMTP_USERNAME"]
	smtpPassword := envFile["SMTP_PASSWORD"]
	smtpHost := envFile["SMTP_HOST"]
	smtpPortRaw := envFile["SMTP_PORT"]
	smtpPort, err := strconv.Atoi(smtpPortRaw)
	if err != nil {
		log.Fatal(err)
	}

	GlobalMailer = Mailer{
		Host:     smtpHost,
		Port:     int(smtpPort),
		Username: smtpUsername,
		Password: smtpPassword,
	}
	GlobalMailer.Initialize(smtpUsername, smtpPassword, smtpHost)

	db, err := gorm.Open(sqlite.Open("database/user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&ErrorDetails{})

	user_db := newUserDB(db)

	// Initialize the user database with a default admin user
	var adminUser User
	user_db.DB.First(&adminUser, "username = ?", "admin")
	if adminUser.Username == "" {
		hashedPassword, err := HashPassword("admin")
		if err != nil {
			panic("failed to hash password")
		}
		adminUser = User{
			Username:      "admin",
			Password:      hashedPassword,
			UserRole:      "administrator",
			EmailVerified: true,
		}
		user_db.DB.Create(&adminUser)
	}

}

type Router struct {
	DB              *gorm.DB
	UserDB          *UserDB
	Mux             *http.ServeMux
	Context         context.Context
	APIRouter       http.Handler
	DashboardRouter http.Handler
	PublicRouter    http.Handler
}

func NewRouter(context context.Context, db *gorm.DB) *Router {
	userDB := newUserDB(db)
	r := &Router{
		DB:      db,
		Mux:     http.NewServeMux(),
		Context: context,
		UserDB:  &userDB,
	}
	r.routes()

	return r
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Router.ServeHTTP" + req.URL.Path)
	router.Mux.ServeHTTP(w, req)
}

func (router *Router) routes() {
	// API routes
	router.Mux.HandleFunc("POST /api/auth/login", router.api_auth_login)
	router.Mux.HandleFunc("POST /api/auth/logout", router.api_auth_logout)
	router.Mux.HandleFunc("GET /api/auth/status", router.api_auth_status)
	router.Mux.HandleFunc("POST /api/auth/register", router.api_auth_register)
	router.Mux.HandleFunc("GET /api/auth/verify-email", router.api_auth_verify_email)

}

func (router *Router) api_report_error(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data ErrorDetails
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}
	log.Printf("Received: %+v", data)

	// Check token in the data matches the one in the db

	// Insert the error into the database
	router.DB.Create(&data)

	// Tell the client that the error was successfully logged
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Success\"}"))
}

func (router *Router) api_auth_login(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "overlord-session")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Authenticate user
	userDB := router.UserDB

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	// Unmarshal the JSON data into the struct
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	// Use the data
	log.Printf("Received: %+v", data)
	log.Println(data["username"])
	log.Println(data["password"])

	//csrf := data["csrf"].(string)
	//nonce := data["nonce"].(string)
	user := userDB.FindByUsername(data["username"].(string))
	passwordVerified := user.CheckPassword(data["password"].(string))

	if !passwordVerified {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set user as authenticated
	session.Values["userID"] = user.ID
	session.Values["authenticated"] = true
	log.Println(session.Values)
	session.Save(r, w)

	// Update user's last login time
	user.LastLoginAt = time.Now()
	userDB.UpdateUser(user)

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Success\"}"))
}

func (router *Router) api_auth_logout(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "overlord-session")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Success\"}"))
}

func (router *Router) api_auth_status(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "overlord-session")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"message\": \"Unauthorized\"}"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Authenticated\"}"))
}

func (router *Router) api_auth_verify_email(w http.ResponseWriter, r *http.Request) {
	session, _ := Store.Get(r, "overlord-session")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the token from the URL query parameter
	username := r.URL.Query().Get("username")
	token := r.URL.Query().Get("token")

	// Authenticate user
	userDB := router.UserDB

	user := userDB.FindByUsername(username)

	tokenIsValid := user.CheckEmailToken(token)

	if !tokenIsValid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set user as authenticated
	session.Values["userID"] = user.ID
	session.Values["authenticated"] = true
	log.Println(session.Values)
	session.Save(r, w)

	// Update user's last login time
	user.LastLoginAt = time.Now()
	user.EmailVerified = true
	userDB.UpdateUser(user)

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Success\"}"))
}

func (router *Router) api_auth_register(w http.ResponseWriter, r *http.Request) {
	userDB := router.UserDB

	type registerForm struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
		Forename  string `json:"forename"`
		Surname   string `json:"surname"`
		Phone     string `json:"phone"`
		Birthdate string `json:"birthdate"`
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// Unmarshal the JSON data into the struct
	var data registerForm
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Use the data
	log.Printf("Received: %+v", data)

	// Check if user already exists
	var user User
	userDB.DB.First(&user, "username = ?", data.Username)
	if user.Username != "" {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	token, err := GenerateToken(data.Email)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	hashedEmailToken, err := HashEmailToken(data.Email)
	if err != nil {
		http.Error(w, "Error hashing email token", http.StatusInternalServerError)
		return
	}

	parsedBirthdate, err := time.Parse("2006-01-02", data.Birthdate)
	if err != nil {
		http.Error(w, "Error parsing birthdate", http.StatusInternalServerError)
		return
	}

	user = User{
		Username:      data.Username,
		Password:      hashedPassword,
		Email:         data.Email,
		EmailToken:    hashedEmailToken,
		Forename:      data.Forename,
		Surname:       data.Surname,
		PhoneNumber:   data.Phone,
		Birthdate:     parsedBirthdate,
		EmailVerified: false,
		PhoneVerified: false,
		UserRole:      "user",
	}
	user.ID, err = userDB.CreateUser(&user)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// URL encode the token
	token = url.QueryEscape(token)

	// Send verification email
	email := Email{
		To:      []string{user.Email},
		From:    APP_CONFIG["SMTP_USERNAME"],
		Subject: "Verify your email address",
		Body:    "Please verify your email address by clicking the link below:\n\n" + APP_CONFIG["SITE_URL"] + "/verify?token=" + token + "&username=" + user.Username,
	}

	err = GlobalMailer.Send(email)
	if err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	session, _ := Store.Get(r, "overlord-session")

	// Set user as authenticated
	session.Values["userID"] = user.ID
	session.Values["authenticated"] = true
	log.Println(session.Values)
	session.Save(r, w)

	// Update user's last login time
	user.LastLoginAt = time.Now()
	userDB.UpdateUser(user)

	// return success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{\"message\": \"Success\"}"))
}
