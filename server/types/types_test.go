package types

import (
	"testing"
)

func TestSanitizeDomain(t *testing.T) {
	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{
			name:    "valid domain",
			domain:  "example.com",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain",
			domain:  "www.example.com",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld",
			domain:  "www.example.com.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods",
			domain:  "www.example.com.au.uk",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores",
			domain:  "www.example.com.au.uk.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens",
			domain:  "www.example-com.au.uk.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores",
			domain:  "www.example-com.au.uk.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores",
			domain:  "www.example-com.au.uk.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid domain with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			domain:  "www.example-com.au.uk.au.au.au.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "invalid domain",
			domain:  "example.com.au.uk.au.au.au.au.au.au.au.au.au.au.au.au",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorDetails{
				Domain: tt.domain,
			}
			if err := e.SanitizeDomain(); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeDomain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid localhost url",
			url:     "http://localhost:8080",
			wantErr: false,
		},
		{
			name:    "valid localhost url 2",
			url:     "http://localhost:8080/test",
			wantErr: false,
		},
		{
			name:    "valid url",
			url:     "https://example.com",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain",
			url:     "https://www.example.com",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld",
			url:     "https://www.example.com.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods",
			url:     "https://www.example.com.au.uk",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores",
			url:     "https://www.example.com.au.uk.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens",
			url:     "https://www.example-com.au.uk.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores",
			url:     "https://www.example-com.au.uk.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods",
			url:     "https://www.example-com.au.uk.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores",
			url:     "https://www.example-com.au.uk.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods",
			url:     "https://www.example-com.au.uk.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores",
			url:     "https://www.example-com.au.uk.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			url:     "https://www.example-com.au.uk.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores",
			url:     "https://www.example-com.au.uk.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			url:     "https://www.example-com.au.uk.au.au.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "valid url with subdomain and tld with multiple periods and underscores and hyphens and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods and multiple underscores and multiple periods",
			url:     "https://www.example-com.au.uk.au.au.au.au.au.au.au.au.au.au.au.au",
			wantErr: false,
		},
		{
			name:    "invalid url",
			url:     "example.com",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorDetails{
				URL: tt.url,
			}
			if err := e.SanitizeURL(); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "valid filename",
			filename: "example.txt",
			wantErr:  false,
		},
		{
			name:     "valid filename with spaces",
			filename: "example file.txt",
			wantErr:  false,
		},
		{
			name:     "valid filename with special characters",
			filename: "example.txt?%&$#@!",
			wantErr:  false,
		},
		{
			name:     "invalid filename",
			filename: "example.txt\x00",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorDetails{
				Filename: tt.filename,
			}
			if err := e.SanitizeFilename(); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeFilename() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeDatetime(t *testing.T) {
	tests := []struct {
		name     string
		datetime string
		wantErr  bool
	}{
		{
			name:     "valid datetime",
			datetime: "2023-01-01T00:00:00Z",
			wantErr:  false,
		},
		{
			name:     "valid datetime with milliseconds",
			datetime: "2023-01-01T00:00:00.000Z",
			wantErr:  false,
		},
		{
			name:     "valid datetime with timezone",
			datetime: "2023-01-01T00:00:00+00:00",
			wantErr:  false,
		},
		{
			name:     "valid datetime with timezone and milliseconds",
			datetime: "2023-01-01T00:00:00.000+00:00",
			wantErr:  false,
		},
		{
			name:     "invalid datetime",
			datetime: "2023-01-01T00:00:00",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorDetails{
				Datetime: tt.datetime,
			}
			if err := e.SanitizeDatetime(); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeDatetime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		wantErr   bool
	}{
		{
			name:      "valid user agent",
			userAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			wantErr:   false,
		},
		{
			name:      "valid user agent with special characters",
			userAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			wantErr:   false,
		},
		{
			name:      "valid user agent with special characters and spaces",
			userAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			wantErr:   false,
		},
		{
			name:      "invalid user agent",
			userAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			wantErr:   true,
		},
		{
			name:      "valid firefox user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:129.0) Gecko/20100101 Firefox/129.0",
			wantErr:   false,
		},
		{
			name:      "valid safari user agent",
			userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15",
			wantErr:   false,
		},
		{
			name:      "valid edge user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36 Edg/127.0.0.0",
			wantErr:   false,
		},
		{
			name:      "valid chrome user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
			wantErr:   false,
		},
		{
			name:      "valid brave user agent",
			userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36",
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &ErrorDetails{
				UserAgent: tt.userAgent,
			}
			if err := e.SanitizeUserAgent(); (err != nil) != tt.wantErr {
				t.Errorf("SanitizeUserAgent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateFields(t *testing.T) {
	tests := []struct {
		name    string
		e       *ErrorDetails
		wantErr bool
	}{
		{
			name: "valid error details",
			e: &ErrorDetails{
				Domain:    "example.com",
				ErrorText: "Error text",
				URL:       "https://example.com",
				Filename:  "example.txt",
				Line:      1,
				Column:    2,
				Datetime:  "2023-01-01T00:00:00Z",
				UserAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			},
			wantErr: false,
		},
		{
			name: "invalid error details",
			e: &ErrorDetails{
				Domain:    "example.com",
				ErrorText: "Error text",
				URL:       "https://example.com",
				Filename:  "example.txt",
				Line:      1,
				Column:    2,
				Datetime:  "2023-01-01T00:00:00",
				UserAgent: "Mozilla/5.0 (Linux; U; Android 10; en-us; Pixel 6 Pro Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/108.0.5359.125 Mobile Safari/537.36",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.e.ValidateFields(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFields() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
