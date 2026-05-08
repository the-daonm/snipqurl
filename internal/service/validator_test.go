package service

import (
	"testing"
)

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://google.com", false},
		{"valid https", "https://google.com", false},
		{"missing scheme", "google.com", true},
		{"invalid scheme", "ftp://google.com", true},
		{"empty string", "", true},
		{"invalid url", "http://192.168.0.%31", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
