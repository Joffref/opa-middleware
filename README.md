# Open Policy Agent Middleware

This middleware integrates Open Policy Agent (OPA) to your http/gin/fiber app.
You can use it to enforce policies on endpoints.
You can use OPA as local policy engine, or as a remote policy engine.

## Installation

```bash
go get github.com/Joffref/opa-middleware
```

## Usage Generic with OPA and HTTP

### Local based policy engine

```go
package main

import (
	"github.com/Joffref/opa-middleware"
	"github.com/Joffref/opa-middleware/config"
	"net/http"
)

var Policy = `
package policy

default allow = false

allow {
	input.path = "/api/v1/users"
	input.method = "GET"
}`

type H struct {
	Name string
}

func (h *H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World : " + h.Name))
}

func main() {
	handler, err := opamiddleware.NewHTTPMiddleware(
		&config.Config{
			Policy: Policy,
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
		&H{
			Name: "John Doe",
		},
	)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler.ServeHTTP)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
```

### Remote based policy engine
The policy is the same as above, but the policy is stored in a remote server.
```go
package main

import (
	"github.com/Joffref/opa-middleware"
	"github.com/Joffref/opa-middleware/config"
	"net/http"
)

type H struct {
	Name string
}

func (h *H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World : " + h.Name))
}

func main() {
	handler, err := opamiddleware.NewHTTPMiddleware(
		&config.Config{
			URL: "http://localhost:8181/",
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
		&H{
			Name: "John Doe",
		},
	)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", handler.ServeHTTP)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
```

## Usage with GIN
```go
package main

import (
	"github.com/Joffref/opa-middleware"
	"github.com/Joffref/opa-middleware/config"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	middleware, err := opamiddleware.NewGinMiddleware(
		&config.Config{
			URL:           "http://localhost:8181/",
			Query:            "data.policy.allow",
			ExceptedResult:   true,
			DeniedStatusCode: 403,
			DeniedMessage:    "Forbidden",
		},
		func(c *gin.Context) (map[string]interface{}, error) {
			return map[string]interface{}{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			}, nil
		},
	)
	if err != nil {
		return
	}
	r.Use(middleware.Use())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}
```

## Usage with Fiber
```go
package main

import (
	"github.com/Joffref/opa-middleware"
	"github.com/Joffref/opa-middleware/config"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	middleware, err := opamiddleware.NewFiberMiddleware(
		&config.Config{
			URL:           "http://localhost:8181/",
			Query:            "data.policy.allow",
			ExceptedResult:   true,
			DeniedStatusCode: 403,
			DeniedMessage:    "Forbidden",
		},
		func(c *fiber.Ctx) (map[string]interface{}, error) {
			return map[string]interface{}{
				"path":   c.Path(),
				"method": c.Method(),
			}, nil
		},
	)
	if err != nil {
		return
	}
	app.Use(middleware.Use())
	app.Get("/ping", func(c *fiber.Ctx) error {
		err := c.JSON("pong")
		if err != nil {
			return err
		}
		return nil
	})
	app.Listen(":8080")
}
```