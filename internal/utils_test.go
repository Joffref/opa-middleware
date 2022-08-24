package internal

import (
	"github.com/Joffref/opa-middleware/config"
	"net/http"
	"reflect"
	"testing"
)

func Test_buildHeaders(t *testing.T) {
	type args struct {
		r   *http.Request
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    http.Header
		wantErr bool
	}{
		{
			name: "build headers",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Accept": []string{"application/json"},
					},
				},
				cfg: &config.Config{
					IgnoredHeaders: []string{"Accept"},
				},
			},
			want:    http.Header{},
			wantErr: false,
		},
		{
			name: "build headers with IgnoredHeaders that are not in the request",
			args: args{
				r: &http.Request{
					Header: http.Header{
						"Accept": []string{"application/json"},
					},
				},
				cfg: &config.Config{
					IgnoredHeaders: []string{"Accept", "Accept-Encoding"},
				},
			},
			want:    http.Header{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildHeaders(tt.args.r, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildHeaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildHeaders() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildURL(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "build url",
			args: args{
				cfg: &config.Config{
					URL:   "http://localhost:8080",
					Query: "data.query",
				},
			},
			want:    "http://localhost:8080/v1/data/query",
			wantErr: false,
		},
		{
			name: "build url with trailing slash",
			args: args{
				cfg: &config.Config{
					URL:   "http://localhost:8080/",
					Query: "data.query",
				},
			},
			want:    "http://localhost:8080/v1/data/query",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildURL(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("buildURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
