package ginmiddleware

import (
	"github.com/Joffref/opa-middleware/config"
	"github.com/gin-gonic/gin"
	"reflect"
	"testing"
)

func TestMiddleware_Query(t *testing.T) {
	type fields struct {
		Config              *config.Config
		InputCreationMethod InputCreationMethod
	}
	type args struct {
		c *gin.Context
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
			g := &Middleware{
				Config:              tt.fields.Config,
				InputCreationMethod: tt.fields.InputCreationMethod,
			}
			got, err := g.query(tt.args.c)
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

func TestMiddleware_Use(t *testing.T) {
	type fields struct {
		Config              *config.Config
		InputCreationMethod InputCreationMethod
	}
	tests := []struct {
		name   string
		fields fields
		want   func(c *gin.Context)
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Middleware{
				Config:              tt.fields.Config,
				InputCreationMethod: tt.fields.InputCreationMethod,
			}
			if got := g.Use(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Use() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGinMiddleware(t *testing.T) {
	type args struct {
		cfg   *config.Config
		input InputCreationMethod
	}
	tests := []struct {
		name    string
		args    args
		want    *Middleware
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGinMiddleware(tt.args.cfg, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGinMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGinMiddleware() got = %v, want %v", got, tt.want)
			}
		})
	}
}
