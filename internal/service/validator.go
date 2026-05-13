package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

var ErrInvalidURL = errors.New("invalid url")

func normalizeAndValidateURL(rawURL string) (string, error) {
	if rawURL == "" {
		return "", ErrInvalidURL
	}

	// Prepend http:// if no scheme is provided
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("%w: invalid scheme %s", ErrInvalidURL, u.Scheme)
	}
	return rawURL, nil
}
