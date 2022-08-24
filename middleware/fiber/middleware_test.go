package fibermiddleware

import (
	"github.com/Joffref/opa-middleware/config"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/http"
	"reflect"
	"testing"
)

func TestMiddleware_Query(t *testing.T) {
	type fields struct {
		Config              *config.Config
		InputCreationMethod InputCreationMethod
	}
	type args struct {
		c *fiber.Ctx
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
		want   func(c *fiber.Ctx) error
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

func TestNewFiberMiddleware(t *testing.T) {
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
			got, err := NewFiberMiddleware(tt.args.cfg, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFiberMiddleware() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFiberMiddleware() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformFastHTTP(t *testing.T) {
	type args struct {
		ctx *fasthttp.RequestCtx
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformFastHTTP(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformFastHTTP() = %v, want %v", got, tt.want)
			}
		})
	}
}
