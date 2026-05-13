package service

import (
	"testing"
)

func TestNormalizeAndValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{"valid http", "http://google.com", "http://google.com", false},
		{"valid https", "https://google.com", "https://google.com", false},
		{"missing scheme", "google.com", "http://google.com", false},
		{"invalid scheme", "ftp://google.com", "", true},
		{"empty string", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeAndValidateURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("normalizeAndValidateURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeAndValidateURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
