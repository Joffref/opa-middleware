package httpmiddleware

import (
	"github.com/Joffref/opa-middleware/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHTTPMiddleware_Query(t *testing.T) {
	type fields struct {
		Config *config.Config
		Next   http.Handler
	}
	type args struct {
		req *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPMiddleware{
				Config: tt.fields.Config,
				Next:   tt.fields.Next,
			}
			got, err := h.query(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Query() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHTTPMiddleware_ServeHTTP(t *testing.T) {
	policy := `
package policy

default allow = false

allow {
	input.path = "/api/v1/users"
	input.method = "GET"
}`
	type fields struct {
		Config *config.Config
		Next   http.Handler
	}
	type args struct {
		rw  http.ResponseWriter
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test HTTPMiddleware_ServeHTTP",
			fields: fields{
				Config: &config.Config{
					Policy: policy,
					Query:  "data.policy.allow",
					InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
						return map[string]interface{}{
							"path":   r.URL.Path,
							"method": r.Method,
						}, nil
					},
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			args: args{
				rw: &httptest.ResponseRecorder{},
				req: &http.Request{
					URL: &url.URL{
						Path: "/api/v1/users",
					},
					Method: "GET",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HTTPMiddleware{
				Config: tt.fields.Config,
				Next:   tt.fields.Next,
			}
			h.ServeHTTP(tt.args.rw, tt.args.req)
		})
	}
}

func TestNewHttpMiddleware(t *testing.T) {
	type args struct {
		cfg  *config.Config
		next http.Handler
	}
	tests := []struct {
		name    string
		args    args
		want    *HTTPMiddleware
		wantErr bool
	}{
		{
			name: "Test NewHTTPMiddleware",
			args: args{
				cfg: &config.Config{
					Policy: "policy",
					Query:  "data.query",
					InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
						return map[string]interface{}{
							"path":   r.URL.Path,
							"method": r.Method,
						}, nil
					},
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			want: &HTTPMiddleware{
				Config: &config.Config{
					Policy: "policy",
					Query:  "data.query",
					InputCreationMethod: func(r *http.Request) (map[string]interface{}, error) {
						return map[string]interface{}{
							"path":   r.URL.Path,
							"method": r.Method,
						}, nil
					},
					ExceptedResult:   true,
					DeniedStatusCode: 403,
					DeniedMessage:    "Forbidden",
				},
				Next: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewHTTPMiddleware(tt.args.cfg, tt.args.next)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHTTPMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
