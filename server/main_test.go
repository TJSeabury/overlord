package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"tjseabury/overlord/types"

	"github.com/joho/godotenv"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestPostEndpoint(t *testing.T) {
	envFile, _ := godotenv.Read(".env")

	// get the values from the environment variables from .env file
	APP_CONFIG = envFile

	// Create the data directory if it doesn't exist
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	db, err := gorm.Open(sqlite.Open("data/app_test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&types.ErrorDetails{})

	// Setup
	router := NewRouter(context.Background(), db)

	server := httptest.NewServer(router)
	defer server.Close()

	var tests = []struct {
		name    string
		data    types.ErrorDetails
		wantErr bool
	}{
		{
			name: "test-1",
			data: types.ErrorDetails{
				Domain:    "whatever.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "https://vwhatever.com/path/to/resource",
				Filename:  "app.js",
				Line:      42,
				Column:    7,
				Datetime:  "2023-10-02T15:04:05Z",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0",
			},
			wantErr: false,
		},
		{
			name: "test-2",
			data: types.ErrorDetails{
				Domain:    "butter.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "https://butter.com/path/to/resource",
				Filename:  "app.js",
				Line:      42,
				Column:    7,
				Datetime:  "invalid datetime",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
			},
			wantErr: true,
		},
		{
			name: "test-3",
			data: types.ErrorDetails{
				Domain:    "invalid.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "htp:/invalid-url",
				Filename:  "app.js",
				Line:      42,
				Column:    7,
				Datetime:  "2023-10-02T15:04:05Z",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
			},
			wantErr: true,
		},
		{
			name: "test-4",
			data: types.ErrorDetails{
				Domain:    "invalid domain.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "https://invalid-domain.com/path/to/resource",
				Filename:  "app.js",
				Line:      42,
				Column:    7,
				Datetime:  "2023-10-02T15:04:05Z",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
			},
			wantErr: true,
		},
		{
			name: "test-5",
			data: types.ErrorDetails{
				Domain:    "greengecko.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "https://greengecko.com/path/to/resource",
				Filename:  "invalid|filename",
				Line:      42,
				Column:    7,
				Datetime:  "2023-10-02T15:04:05Z",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.63 Safari/537.36",
			},
			wantErr: true,
		},
		{
			name: "test-6",
			data: types.ErrorDetails{
				Domain:    "complexuser.com",
				ErrorText: "Memory allocation error at line 150",
				URL:       "https://complexuser.com/path/to/resource",
				Filename:  "app.js",
				Line:      42,
				Column:    7,
				Datetime:  "2023-10-02T15:04:05Z",
				UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36 Edg/127.0.0.0",
			},
			wantErr: false,
		},
		// {
		// 	name: "test-7",
		// 	data: types.ErrorDetails{
		// 		Domain:    "redbluegreen.com",
		// 		ErrorText: "Memory allocation error at line 150",
		// 		URL:       "https://redbluegreen.com/path/to/resource",
		// 		Filename:  "app.js",
		// 		Line:      42,
		// 		Column:    7,
		// 		Datetime:  "2023-10-02T15:04:05Z",
		// 		UserAgent: "InvalidUserAgent/1.0",
		// 	},
		// 	wantErr: true,
		// },
		{
			name:    "test-8",
			data:    types.ErrorDetails{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.data)
			if err != nil {
				t.Fatalf("Failed to marshal JSON: %v", err)
			}

			// Create a POST request
			req, err := http.NewRequest("POST", server.URL+"/api/report-error", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Failed to send request: %v", err)
			}
			defer resp.Body.Close()

			// Assert the response
			if tt.wantErr {
				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status BadRequest; got %v", resp.Status)
				}
			} else {
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status OK; got %v", resp.Status)
				}
			}
		})
	}

	// remove the test database
	os.Remove("data/app_test.db")

}
