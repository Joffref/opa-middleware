package config

import (
	"errors"
	"github.com/open-policy-agent/opa/rego"
	"log"
	"net/http"
	"time"
)

// Config is the configuration for the middleware.
type Config struct {
	// URL is the URL of the OPA server, if not set Middleware will use Policy field.
	// You must set either URL or Policy.
	URL string `json:"url,omitempty"`

	// Policy is the policy document to use, if not set Middleware will use URL field.
	// You must set either URL or Policy.
	Policy string `json:"policy,omitempty"`

	instantiatedPolicy *rego.PreparedEvalQuery

	// Query is the name of the policy to query.
	Query string `json:"query,omitempty"`

	// InputCreationMethod is a function that returns the value to be sent to the OPA server.
	InputCreationMethod func(r *http.Request) (map[string]interface{}, error) `json:"binding_method,omitempty"`

	// ExceptedResult is the result that the OPA server should return if the request is allowed.
	// Default to true. At the moment only boolean values are supported.
	ExceptedResult bool `json:"excepted_result,omitempty"`

	// DeniedStatusCode is the status code that should be returned if the request is denied.
	DeniedStatusCode int `json:"denied_status,omitempty"`

	// DeniedMessage is the message that should be returned if the request is denied.
	DeniedMessage string `json:"denied_message,omitempty"`

	// Headers is a list of headers to send to the OPA server.
	// All headers are sent to the OPA server except those in the IgnoredHeaders list.
	Headers map[string][]string `json:"headers,omitempty"`

	// IgnoredHeaders is a list of headers to ignore when sending to the OPA server.
	IgnoredHeaders []string `json:"excepted_headers,omitempty"`

	// Debug is a flag that enables debug mode.
	Debug bool `json:"debug,omitempty"`

	// Logger is the logger that is use when debug mode is enabled.
	// If not set, the default logger is used.
	Logger *log.Logger `json:"logger,omitempty"`

	// Timeout is the timeout for the request policy evaluation.
	// If not set, the default is 10 seconds.
	Timeout time.Duration `json:"timeout,omitempty"`
}

func (c *Config) Validate() error {
	if c.Debug {
		if c.Logger == nil {
			c.Logger = log.Default()
		}
	}
	if c.Timeout == 0 {
		c.Timeout = 10 * time.Second
	}
	if c.URL == "" && c.Policy == "" {
		return errors.New("[opa-middleware] You must set either URL or Policy")
	}
	if c.URL != "" && c.Policy != "" {
		return errors.New("[opa-middleware] You must set either URL or Policy")
	}
	if c.ExceptedResult != true && c.ExceptedResult != false {
		return errors.New("[opa-middleware] You must set ExceptedResult")
	}
	if c.Query == "" {
		return errors.New("[opa-middleware] You must set Query")
	}
	return nil
}
