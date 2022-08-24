package httpmiddleware

import (
	"github.com/Joffref/opa-middleware/config"
	"github.com/Joffref/opa-middleware/internal"
	"net/http"
)

// HTTPMiddleware is the middleware for http requests
type HTTPMiddleware struct {
	Config *config.Config
	// Next is the next handler in the request chain.
	Next http.Handler
}

// NewHTTPMiddleware returns a new HTTPMiddleware
func NewHTTPMiddleware(cfg *config.Config, next http.Handler) (*HTTPMiddleware, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	if next == nil {
		return nil, err
	}
	return &HTTPMiddleware{
		Config: cfg,
		Next:   next,
	}, nil
}

// ServeHTTP serves the http request. Act as Use acts in other frameworks.
func (h *HTTPMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if h.Config.Debug {
		h.Config.Logger.Printf("[opa-middleware-http] Request: %s", req.URL.String())
	}
	result, err := h.query(req)
	if err != nil {
		if h.Config.Debug {
			h.Config.Logger.Printf("[opa-middleware-http] Error: %s", err.Error())
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	if h.Config.Debug {
		h.Config.Logger.Printf("[opa-middleware-http] Result: %t", result)
	}
	if result != h.Config.ExceptedResult {
		http.Error(rw, h.Config.DeniedMessage, h.Config.DeniedStatusCode)
		return
	}
	h.Next.ServeHTTP(rw, req)
}

func (h *HTTPMiddleware) query(req *http.Request) (bool, error) {
	bind, err := h.Config.InputCreationMethod(req)
	if err != nil {
		return !h.Config.ExceptedResult, err
	}
	if h.Config.URL != "" {
		input := make(map[string]interface{})
		input["input"] = bind
		return internal.QueryURL(req, h.Config, input)
	}
	return internal.QueryPolicy(req, h.Config, bind)
}
