package internal

import (
	"github.com/Joffref/opa-middleware/config"
	"net/http"
	"testing"
)

func TestQueryPolicy(t *testing.T) {
	policy := `
package policy

default allow = false

allow {
	input.path = "/api/v1/users"
	input.method = "GET"
}`

	type args struct {
		r    *http.Request
		cfg  *config.Config
		bind map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "query policy should return true",
			args: args{
				r: &http.Request{},
				cfg: &config.Config{
					Policy: policy,
					Query:  "data.policy.allow",
				},
				bind: map[string]interface{}{
					"path":   "/api/v1/users",
					"method": "GET",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "query policy should return false",
			args: args{
				r: &http.Request{},
				cfg: &config.Config{
					Policy: policy,
					Query:  "data.policy.allow",
				},
				bind: map[string]interface{}{
					"path":   "/api/v1/users",
					"method": "POST",
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryPolicy(tt.args.r, tt.args.cfg, tt.args.bind)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryPolicy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryPolicy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryURL(t *testing.T) {
	_ = `
package policy

default allow = false

allow {
	input.path = "/api/v1/users"
	input.method = "GET"
}`
	type args struct {
		r    *http.Request
		cfg  *config.Config
		bind map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		/*
			{
				name: "query url should return true",
				args: args{
					r: &http.Request{
						URL: &url.URL{
							Path: "/api/v1/users",
						},
					},
					cfg: &config.Config{
						URL:   "data.url.path",
						Query: "data.url.path == \"/api/v1/users\"",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryURL(tt.args.r, tt.args.cfg, tt.args.bind)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
