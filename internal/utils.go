package internal

import (
	"github.com/Joffref/opa-middleware/config"
	"net/http"
	"net/url"
	"strings"
)

func buildHeaders(r *http.Request, cfg *config.Config) (http.Header, error) {
	headers := r.Header.Clone()
	for _, header := range cfg.IgnoredHeaders {
		headers.Del(header)
	}
	for header, values := range cfg.Headers {
		for _, value := range values {
			headers.Set(header, value)
		}
	}
	return headers, nil
}

func buildURL(cfg *config.Config) (string, error) {
	// Remove trailing slash from URL if present.
	cfg.URL = strings.TrimSuffix(cfg.URL, "/")
	u := cfg.URL + "/v1/" + strings.Replace(cfg.Query, ".", "/", -1)
	uCleaned, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return uCleaned.String(), nil
}
