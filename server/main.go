package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"tjseabury/overlord/types"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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

	// Create the data directory if it doesn't exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	db, err := gorm.Open(sqlite.Open("data/app.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&types.ErrorDetails{})

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

	// Start the mock error generator
	if APP_CONFIG["MOCK_MODE"] == "TRUE" {
		ctx := context.Background()
		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					// If the context is done, return from the goroutine
					return
				default:
					// If the context is not done, continue with the loop
					mockError := random_mock_error()
					db.Create(mockError)
					fmt.Println("Inserting mock error.")
					time.Sleep(time.Second * time.Duration(rand.Intn(30)))
				}
			}
		}(ctx)
	}

	// Create the router
	router := NewRouter(context.Background(), db)

	// Read in the shadow watcher script
	shadowWatcherScript, err := os.ReadFile("../client/dist/ShadowWatcher.js")
	if err != nil {
		log.Fatal(err)
	}
	// serve the shadow watcher script
	router.Mux.HandleFunc("GET /ShadowWatcher", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(shadowWatcherScript)
	})

	router.Mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<title>Overlord</title>
	<script src="/ShadowWatcher" data-token="test-token"></script>
</head>
<body style="background-color: #444; color: #fff;">
	<h1>Overlord</h1>
	<p>This is the Overlord dashboard.</p>
	<script>
	// This is a test script that will produce an error
	const test = function() {
		console.log("Test function called");
		throw new Error("This is a test error");
	};
	test();
	</script>
</body>
</html>`))
	})

	log.Fatal(http.ListenAndServe(":8080", router))
}

type Router struct {
	DB              *gorm.DB
	UserDB          *UserController
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
	router.Mux.HandleFunc("POST /api/report-error", router.api_report_error)
	router.Mux.HandleFunc("GET /", router.handle_dashboard)
}

func (router *Router) handle_dashboard(w http.ResponseWriter, r *http.Request) {
	errors := make([]types.ErrorDetails, 0)
	router.DB.Find(&errors)

	// Read in the dashboard template file
	dashboardTemplateFile, err := os.ReadFile("templates/dashboard.html")
	if err != nil {
		log.Fatal(err)
	}

	dashboard_template := template.Must(template.New("dashboard").Parse(
		string(dashboardTemplateFile),
	))

	w.Header().Set("Content-Type", "text/html")
	dashboard_template.Execute(w, errors)
}

func (router *Router) api_report_error(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")                                // Allow any origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allowed methods
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")     // Allowed headers

	// If it's a preflight OPTIONS request, send an OK status and return
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var data types.ErrorDetails
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error parsing JSON body", http.StatusBadRequest)
		return
	}

	// Validate the data
	// err := data.ValidateFields()
	// if err != nil {
	// 	log.Println(data, err)
	// 	http.Error(w, "Invalid data", http.StatusBadRequest)
	// 	return
	// }
	// Check token in the data matches the one in the db

	// Insert the error into the database
	router.DB.Create(&data)

	log.Printf("Inserted: %+v", data)

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
