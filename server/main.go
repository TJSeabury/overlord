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

func random_mock_error() types.ErrorDetails {

	// Define slices with a large set of mock data
	var domains = []string{
		"auth.example.com",
		"api.example.org",
		"frontend.example.net",
		"payments.example.com",
		"notifications.example.org",
		"shop.example.co",
		"blog.example.com",
		"docs.example.org",
		"support.example.net",
		"media.example.com",
		"auth.example.com",
		"api.example.org",
		"frontend.example.net",
		"payments.example.com",
		"notifications.example.org",
		"shop.example.co",
		"blog.example.com",
		"docs.example.org",
		"support.example.net",
		"media.example.com",
		"chat.example.com",
		"data.example.org",
		"events.example.net",
		"maps.example.com",
		"media.example.co",
		"files.example.com",
		"analytics.example.org",
		"dashboard.example.net",
		"streaming.example.com",
		"search.example.org",
		"wiki.example.com",
		"jobs.example.net",
		"forum.example.org",
		"news.example.com",
		"community.example.net",
		"research.example.org",
		"api.example.net",
		"platform.example.com",
		"ads.example.org",
		"updates.example.com",
		"resources.example.net",
		"support.example.com",
		"checkout.example.org",
		"gallery.example.com",
		"events.example.com",
		"services.example.net",
		"blogging.example.org",
		"store.example.com",
		"newsletter.example.net",
		"feedback.example.org",
		"tools.example.com",
		"training.example.net",
		"knowledge.example.org",
		"play.example.com",
		"docs.example.net",
		"portal.example.org",
		"security.example.com",
		"integration.example.net",
		"webinars.example.org",
		"applications.example.com",
		"projects.example.net",
		"forums.example.org",
		"mediahub.example.com",
		"solutions.example.net",
		"platforms.example.org",
		"members.example.com",
		"resources.example.org",
		"services.example.co",
	}

	var error_texts = []string{
		"Invalid credentials provided.",
		"Failed to parse JSON response.",
		"Uncaught TypeError: Cannot read property 'foo' of undefined",
		"Transaction failed due to insufficient funds.",
		"Failed to send email notification.",
		"Unhandled promise rejection.",
		"404 Not Found error while fetching resource.",
		"Server timed out while processing request.",
		"Cross-Origin Resource Sharing (CORS) policy error.",
		"Cannot connect to WebSocket server.",
		"Failed to load resource: net::ERR_CONNECTION_REFUSED",
		"SyntaxError: Unexpected token < in JSON at position 0",
		"TypeError: Cannot read property 'bar' of null",
		"ReferenceError: foo is not defined",
		"NetworkError: Failed to fetch",
		"Uncaught Error: Assertion failed",
		"Error: Unable to parse the response",
		"RangeError: Maximum call stack size exceeded",
		"TypeError: Failed to execute 'appendChild' on 'Node': parameter 1 is not of type 'Node'.",
		"TimeoutError: Timeout waiting for page to load",
		"URIError: URI malformed",
		"Error: Cannot find module 'express'",
		"SyntaxError: Unexpected end of JSON input",
		"TypeError: Cannot set property 'value' of undefined",
		"Error: Request failed with status code 500",
		"NetworkError: Unable to connect to server",
		"ReferenceError: event is not defined",
		"DOMException: The operation is insecure.",
		"TypeError: 'undefined' is not a function",
		"Error: Invalid arguments",
		"TypeError: Cannot read property 'length' of undefined",
		"Error: Could not load content for URL",
		"SecurityError: The operation is insecure",
		"Error: Unsupported URL type",
		"Error: No data received from server",
		"Error: Could not parse response",
		"Error: Unsupported media type",
		"ReferenceError: process is not defined",
		"TypeError: Cannot read property 'clientWidth' of null",
		"Error: Fetch API cannot load",
		"Error: Server responded with a status of 403 (Forbidden)",
		"TypeError: Cannot convert undefined or null to object",
		"Error: The user canceled the request",
		"Error: Page not found (404)",
		"DOMException: The operation is not supported",
		"Error: Failed to parse XML",
		"TypeError: Cannot read property 'innerHTML' of null",
		"Error: Unexpected token in JSON",
		"TypeError: Cannot read property 'get' of undefined",
		"ReferenceError: Promise is not defined",
		"Error: Failed to load script",
		"TypeError: The 'data' argument must be of type string or an instance of Buffer or ArrayBuffer. Received type undefined",
		"Error: Unsupported operation",
		"Error: Authentication required",
		"DOMException: Failed to execute 'postMessage' on 'Window': The target origin provided ('http://example.com') does not match the recipient window's origin ('http://example.org').",
		"Error: Invalid JSON response",
		"TypeError: Cannot read property 'toUpperCase' of undefined",
		"SyntaxError: Unexpected end of input",
		"TypeError: Failed to execute 'send' on 'XMLHttpRequest': The 'data' parameter is not a string or a FormData object.",
		"Error: The remote server returned an error: (404) Not Found.",
		"ReferenceError: window is not defined",
		"Error: Failed to parse HTML",
		"Error: Request aborted",
		"Error: Could not find file",
		"Error: Failed to get the data",
		"DOMException: Failed to execute 'setItem' on 'Storage': Setting the value of 'key' exceeded the quota.",
		"TypeError: Failed to fetch",
		"Error: CORS header ‘Access-Control-Allow-Origin’ missing",
		"ReferenceError: fetch is not defined",
		"Error: No such file or directory",
		"Error: JSON parse error",
		"DOMException: The operation is not allowed",
		"TypeError: Cannot read property 'style' of null",
		"Error: Network request failed",
		"SyntaxError: Invalid or unexpected token",
		"TypeError: Cannot read property 'scrollIntoView' of undefined",
		"Error: Network timeout",
		"TypeError: Failed to execute 'querySelector' on 'Document': '[object Object]' is not a valid selector.",
		"Error: Failed to decode JSON",
		"ReferenceError: localStorage is not defined",
		"TypeError: Cannot call method 'split' of undefined",
		"Error: Expected response to be JSON but got text",
		"Error: Failed to parse input data",
	}

	var urls = []string{
		"https://auth.example.com/login",
		"https://api.example.org/data",
		"https://frontend.example.net/dashboard",
		"https://payments.example.com/checkout",
		"https://notifications.example.org/api/send",
		"https://shop.example.co/product/123",
		"https://blog.example.com/post/456",
		"https://docs.example.org/getting-started",
		"https://support.example.net/ticket/789",
		"https://media.example.com/video/101112",
		"https://shop.example.com/products/987",
		"https://blog.example.com/article/1234",
		"https://docs.example.org/user-guide",
		"https://support.example.net/tickets/567",
		"https://media.example.com/assets/images/logo.png",
		"https://chat.example.com/conversations/789",
		"https://data.example.org/stats",
		"https://events.example.net/2024/conference",
		"https://maps.example.com/location/456",
		"https://media.example.co/videos/101",
		"https://files.example.com/uploads/abc123",
		"https://analytics.example.org/reports/2024",
		"https://dashboard.example.net/overview",
		"https://streaming.example.com/live",
		"https://search.example.org/query?term=error",
		"https://wiki.example.com/page/45",
		"https://jobs.example.net/apply",
		"https://forum.example.org/discussion/567",
		"https://news.example.com/latest",
		"https://community.example.net/groups",
		"https://research.example.org/publications",
		"https://api.example.net/v1/resources",
		"https://platform.example.com/tools",
		"https://ads.example.org/campaigns",
		"https://updates.example.com/release-notes",
		"https://resources.example.net/guide",
		"https://checkout.example.org/complete",
		"https://gallery.example.com/albums/12",
		"https://events.example.com/webinars",
		"https://services.example.net/support",
		"https://blogging.example.org/new-post",
		"https://store.example.com/cart",
		"https://newsletter.example.net/archive",
		"https://feedback.example.org/surveys",
		"https://tools.example.com/utility",
		"https://training.example.net/sessions",
		"https://knowledge.example.org/faq",
		"https://play.example.com/games",
		"https://docs.example.net/reference",
		"https://portal.example.org/dashboard",
		"https://security.example.com/alerts",
		"https://integration.example.net/api",
		"https://webinars.example.org/schedule",
		"https://applications.example.com/software",
		"https://projects.example.net/overview",
		"https://forums.example.org/topics",
		"https://mediahub.example.com/library",
		"https://solutions.example.net/consulting",
		"https://platforms.example.org/services",
		"https://members.example.com/profile",
		"https://resources.example.org/tutorials",
		"https://events.example.org/summit",
		"https://shop.example.com/sale",
		"https://blog.example.net/tips",
		"https://docs.example.com/setup",
		"https://support.example.org/help",
		"https://media.example.net/podcasts",
		"https://chat.example.org/groups",
		"https://data.example.net/analytics",
		"https://events.example.com/speakers",
		"https://maps.example.net/maps",
		"https://media.example.com/photos",
		"https://files.example.org/docs",
		"https://analytics.example.net/insights",
		"https://dashboard.example.com/status",
		"https://streaming.example.org/shows",
		"https://search.example.net/results",
		"https://wiki.example.org/updates",
		"https://jobs.example.com/listings",
		"https://forum.example.net/topics",
		"https://news.example.org/headlines",
		"https://community.example.com/forums",
		"https://research.example.net/papers",
		"https://api.example.org/v2/queries",
		"https://platform.example.net/launch",
		"https://ads.example.com/promotions",
		"https://updates.example.org/news",
		"https://resources.example.com/articles",
		"https://checkout.example.net/review",
		"https://gallery.example.org/collections",
		"https://events.example.com/expos",
		"https://services.example.org/client",
		"https://blogging.example.net/archive",
		"https://store.example.org/orders",
		"https://newsletter.example.com/current",
		"https://feedback.example.net/forms",
		"https://tools.example.org/resources",
		"https://training.example.com/courses",
		"https://knowledge.example.net/guides",
		"https://play.example.org/activities",
		"https://docs.example.com/manual",
		"https://portal.example.net/access",
		"https://security.example.org/threats",
		"https://integration.example.com/plugins",
		"https://webinars.example.net/events",
		"https://applications.example.org/apps",
		"https://projects.example.com/cases",
		"https://forums.example.net/chat",
		"https://mediahub.example.org/videos",
		"https://solutions.example.com/strategies",
		"https://platforms.example.net/tools",
		"https://members.example.org/benefits",
		"https://resources.example.com/downloads",
	}

	var filenames = []string{
		"auth_controller.go",
		"data_service.go",
		"app.js",
		"payment_processor.go",
		"notification_service.py",
		"main.go",
		"server.py",
		"utils.js",
		"database_helper.py",
		"router.ts",
		"auth_service.js",
		"data_handler.go",
		"main_controller.py",
		"payment_gateway.ts",
		"notification_manager.js",
		"app_config.json",
		"index.html",
		"api_router.go",
		"utils.js",
		"error_logger.py",
		"server.js",
		"database_migrations.sql",
		"request_handler.ts",
		"user_profile.html",
		"payment_service.go",
		"error_monitor.js",
		"admin_dashboard.py",
		"auth_routes.ts",
		"payment_processor.js",
		"api_client.go",
		"websocket_server.py",
		"email_sender.go",
		"session_manager.js",
		"cache_service.py",
		"logger.ts",
		"file_uploader.js",
		"api_server.go",
		"form_validator.py",
		"data_exporter.ts",
		"authenticator.js",
		"request_validator.go",
		"database_utils.py",
		"image_processor.ts",
		"error_handling.js",
		"email_service.go",
		"web_app.html",
		"file_reader.py",
		"client_interface.ts",
		"payment_handler.js",
		"server_config.go",
		"session_controller.py",
		"data_processor.ts",
		"api_utils.js",
		"response_handler.go",
		"config_manager.py",
		"auth_controller.ts",
		"file_writer.js",
		"payment_api.go",
		"notification_service.py",
		"admin_controller.ts",
		"logger.js",
		"data_syncer.go",
		"websocket_client.py",
		"api_manager.ts",
		"session_service.js",
		"email_handler.go",
		"error_reporter.py",
		"web_config.json",
		"server_manager.ts",
		"request_manager.js",
		"api_connector.go",
		"file_manager.py",
		"data_cleaner.ts",
		"auth_manager.js",
		"cache_controller.go",
		"response_generator.py",
		"payment_utils.ts",
		"websocket_server.go",
		"image_resizer.py",
		"api_interceptor.js",
		"error_tracker.go",
		"notification_controller.py",
		"form_processor.ts",
		"client_service.js",
		"server_utils.go",
		"session_manager.py",
		"data_collector.ts",
	}

	var user_agents = []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Firefox/115.0",
		"Safari/537.36",
		"PostmanRuntime/7.28.4",
		"curl/7.64.1",
		"Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Mobile Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:90.0) Gecko/20100101 Firefox/90.0",
		"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.12.388 Version/12.18",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:91.0) Gecko/20100101 Firefox/91.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36",
		"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:89.0) Gecko/20100101 Firefox/89.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 11_1_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.48",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; Trident/7.0; AS; .NET CLR 4.0.30319; .NET4.0E; .NET4.0C; MAS; MAM; MEI; MSE; MSIE 11.0; MSE; MSE; MAE) like Gecko",
		"Mozilla/5.0 (Android 11; Mobile; rv:90.0) Gecko/90.0 Firefox/90.0",
		"Mozilla/5.0 (Linux; Android 10; SM-G960U) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/537.36 (KHTML, like Gecko) Version/13.1.2 Safari/537.36",
	}

	var stack_traces = []string{
		"at AuthController.login (auth_controller.go:42)\nat main.main (main.go:20)",
		"at DataService.FetchData (data_service.go:87)\nat main.main (main.go:45)",
		"at app.js:102\nat Object.<anonymous> (app.js:150)",
		"at PaymentProcessor.process (payment_processor.go:63)\nat main.main (main.go:78)",
		"at NotificationService.sendEmail (notification_service.py:34)\nat main (server.py:90)",
		"at UserController.getUserProfile (user_controller.js:55)\nat main (app.js:72)",
		"at CheckoutService.processPayment (checkout_service.js:98)\nat main (index.js:25)",
		"at FileUploader.upload (file_uploader.go:12)\nat main.main (main.go:88)",
		"at ApiClient.sendRequest (api_client.py:67)\nat main (server.py:101)",
		"at WebSocketServer.handleConnection (websocket_server.js:44)\nat main (server.js:29)",
		"Error: Failed to fetch\n    at fetch (node_modules/whatwg-fetch/dist/fetch.umd.js:1:875)\n    at Object.getData (src/api.js:14:5)\n    at async src/components/Widget.js:22:12\n    at async src/app.js:5:7",
		"TypeError: Cannot read property 'map' of undefined\n    at processData (src/utils.js:10:7)\n    at src/components/ItemList.js:33:15\n    at Array.map (<anonymous>)\n    at src/components/ItemList.js:31:12\n    at src/app.js:9:5",
		"ReferenceError: foo is not defined\n    at src/utils.js:8:9\n    at src/app.js:14:3\n    at async src/components/MyComponent.js:20:7",
		"Error: Network request failed\n    at fetchData (src/api.js:8:5)\n    at src/components/MyComponent.js:25:12\n    at async src/app.js:7:5",
		"SyntaxError: Unexpected token < in JSON at position 0\n    at src/api.js:12:7\n    at processResponse (src/utils.js:20:5)\n    at src/components/Widget.js:35:10",
		"TypeError: Cannot read property 'length' of null\n    at processItems (src/utils.js:17:5)\n    at src/components/ItemList.js:22:10\n    at src/app.js:13:7",
		"Error: Invalid JSON response\n    at src/api.js:5:10\n    at src/components/ItemList.js:30:5\n    at src/app.js:8:3",
		"ReferenceError: localStorage is not defined\n    at saveData (src/storage.js:10:7)\n    at src/components/Settings.js:22:5\n    at src/app.js:12:8",
		"TypeError: 'undefined' is not a function\n    at processForm (src/forms.js:15:5)\n    at src/components/FormComponent.js:20:10\n    at src/app.js:6:5",
		"Error: Request failed with status code 500\n    at fetchData (src/api.js:10:7)\n    at src/components/MyComponent.js:29:12\n    at src.app.js:11:5",
		"TypeError: Cannot read property 'clientWidth' of null\n    at updateLayout (src/layout.js:8:7)\n    at src/components/Sidebar.js:16:10\n    at src/app.js:10:5",
		"Error: Server responded with a status of 404 (Not Found)\n    at fetchData (src/api.js:12:7)\n    at src/components/MyComponent.js:26:10\n    at src.app.js:9:5",
		"ReferenceError: event is not defined\n    at handleEvent (src/events.js:5:7)\n    at src/components/EventComponent.js:12:5\n    at src/app.js:7:3",
		"TypeError: Cannot read property 'style' of undefined\n    at applyStyles (src/styles.js:14:5)\n    at src/components/MyComponent.js:18:10\n    at src/app.js:6:5",
		"SyntaxError: Unexpected end of JSON input\n    at parseResponse (src/utils.js:20:7)\n    at src/api.js:15:5\n    at src/components/Widget.js:28:12",
		"TypeError: Cannot set property 'value' of undefined\n    at setValue (src/forms.js:12:7)\n    at src/components/FormComponent.js:22:5\n    at src.app.js:11:5",
		"Error: Cannot find module 'express'\n    at Function.Module._resolveFilename (node:internal/modules/cjs/loader:907:15)\n    at Function.Module._load (node:internal/modules/cjs/loader:752:27)\n    at Module.require (node:internal/modules/cjs/loader:975:19)\n    at require (node:internal/modules/cjs/helpers:92:18)\n    at Object.<anonymous> (src/server.js:1:13)",
		"ReferenceError: Promise is not defined\n    at fetchData (src/api.js:7:7)\n    at src.app.js:5:3\n    at src/components/MyComponent.js:21:7",
		"TypeError: Failed to execute 'send' on 'XMLHttpRequest': The 'data' parameter is not a string or a FormData object\n    at sendRequest (src/api.js:9:7)\n    at src/components/RequestComponent.js:19:10\n    at src.app.js:12:5",
		"Error: The remote server returned an error: (403) Forbidden\n    at fetchData (src/api.js:11:7)\n    at src/components/MyComponent.js:28:10\n    at src.app.js:13:5",
		"TypeError: Cannot convert undefined or null to object\n    at src/utils.js:15:7\n    at src/components/MyComponent.js:20:10\n    at src.app.js:9:5",
		"Error: Authentication required\n    at fetchData (src/api.js:10:7)\n    at src/components/AuthComponent.js:22:5\n    at src.app.js:14:5",
		"TypeError: Cannot read property 'scrollIntoView' of undefined\n    at scrollToElement (src/utils.js:12:7)\n    at src/components/Scroller.js:20:10\n    at src.app.js:7:5",
		"Error: Unsupported URL type\n    at fetchData (src/api.js:6:7)\n    at src/components/MyComponent.js:23:5\n    at src.app.js:10:5",
		"SyntaxError: Invalid or unexpected token\n    at src/api.js:8:7\n    at src/components/RequestComponent.js:17:10\n    at src.app.js:8:5",
		"ReferenceError: window is not defined\n    at updateWindowSize (src/utils.js:7:7)\n    at src/components/WindowComponent.js:15:10\n    at src.app.js:6:5",
		"TypeError: Cannot read property 'get' of undefined\n    at fetchData (src/api.js:12:7)\n    at src/components/MyComponent.js:29:10\n    at src.app.js:13:5",
		"Error: The user canceled the request\n    at fetchData (src/api.js:14:7)\n    at src/components/RequestComponent.js:18:10\n    at src.app.js:11:5",
		"TypeError: Failed to fetch\n    at fetchData (src/api.js:9:7)\n    at src/components/MyComponent.js:24:10\n    at src.app.js:7:5",
		"SyntaxError: Unexpected token in JSON\n    at parseResponse (src/utils.js:18:7)\n    at src/api.js:16:5\n    at src/components/Widget.js:25:12",
		"Error: Failed to load resource: net::ERR_CONNECTION_REFUSED\n    at src/api.js:11:7\n    at src/components/MyComponent.js:20:5\n    at src.app.js:12:8",
		"TypeError: Cannot read property 'innerHTML' of null\n    at updateContent (src/utils.js:10:7)\n    at src/components/MyComponent.js:23:10\n    at src.app.js:9:5",
		"ReferenceError: fetch is not defined\n    at fetchData (src/api.js:7:7)\n    at src/components/MyComponent.js:18:10\n    at src.app.js:10:5",
		"Error: Request aborted\n    at fetchData (src/api.js:13:7)\n    at src/components/MyComponent.js:27:10\n    at src.app.js:14:5",
		"TypeError: Cannot read property 'split' of undefined\n    at splitData (src/utils.js:13:7)\n    at src/components/DataComponent.js:22:5\n    at src.app.js:12:5",
		"SyntaxError: Unexpected end of input\n    at src/api.js:16:7\n    at src/components/MyComponent.js:30:10\n    at src.app.js:15:5",
		"TypeError: Cannot call method 'trim' of undefined\n    at trimData (src/utils.js:11:7)\n    at src/components/FormComponent.js:19:10\n    at src.app.js:8:5",
		"ReferenceError: document is not defined\n    at updateDocument (src/utils.js:8:7)\n    at src/components/DocumentComponent.js:15:10\n    at src.app.js:6:5",
		"Error: Fetch API cannot load\n    at fetchData (src/api.js:9:7)\n    at src/components/MyComponent.js:28:10\n    at src.app.js:13:5",
		"TypeError: Cannot read property 'parentNode' of null\n    at updateParent (src/utils.js:15:7)\n    at src/components/MyComponent.js:20:10\n    at src.app.js:12:5",
		"Error: Invalid arguments\n    at src/api.js:6:7\n    at src/components/RequestComponent.js:22:10\n    at src.app.js:10:5",
	}

	// Generate random index for each slice
	randomIndex := func(sliceLen int) int {
		return rand.Intn(sliceLen)
	}

	// Generate a random ErrorDetails struct
	return types.ErrorDetails{
		Domain:     domains[randomIndex(len(domains))],
		ErrorText:  error_texts[randomIndex(len(error_texts))],
		URL:        urls[randomIndex(len(urls))],
		Filename:   filenames[randomIndex(len(filenames))],
		Line:       rand.Intn(100) + 1, // Random line number between 1 and 100
		Column:     rand.Intn(50) + 1,  // Random column number between 1 and 50
		Datetime:   time.Now().Format(time.RFC3339),
		UserAgent:  user_agents[randomIndex(len(user_agents))],
		StackTrace: stack_traces[randomIndex(len(stack_traces))],
	}
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
