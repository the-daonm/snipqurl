package service

import (
	"errors"
	"fmt"
	"net/url"
)

var ErrInvalidURL = errors.New("invalid url")

func validateURL(rawURL string) error {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: invalid scheme %s", ErrInvalidURL, u.Scheme)
	}
	return nil
}
