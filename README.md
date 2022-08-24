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
	"github.com/Joffref/opa-middleware/config"
	"github.com/Joffref/opa-middleware/middleware/http"
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
	handler, err := httpmiddleware.NewHTTPMiddleware(
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
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}
}
```

### Remote based policy engine

```go
package main

import (
	"github.com/Joffref/opa-middleware/config"
	"github.com/Joffref/opa-middleware/middleware/http"
	"net/http"
)

type H struct {
	Name string
}

func (h *H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World : " + h.Name))
}

func main() {
	handler, err := httpmiddleware.NewHTTPMiddleware(
		&config.Config{
			URL:   "http://localhost:8181",
			Query: "data.policy.allow",
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
	err = http.ListenAndServe(":8080", handler)
	if err != nil {
		return
	}
}
```

## Usage with GIN
```go
package main

import (
    "github.com/Joffref/opa-middleware/config"
    ginmiddleware "github.com/Joffref/opa-middleware/middleware/gin"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    middleware, err := ginmiddleware.NewGinMiddleware(
        &config.Config{
            URL:   "https://opa.example.com/",
            Query: "data.policy.allow",
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
    err = r.Run(":8080")
    if err != nil {
        return
    }
}
```

## Usage with Fiber
```go
package main

import (
	"github.com/Joffref/opa-middleware/config"
	fibermiddleware "github.com/Joffref/opa-middleware/middleware/fiber"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})

	middleware, err := fibermiddleware.NewFiberMiddleware(&config.Config{
		URL:              "http://localhost:8080/",
		Query:            "data.policy.allow",
		DeniedStatusCode: 403,
		DeniedMessage:    "Forbidden",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		IgnoredHeaders: []string{
			"X-Request-Id",
		},
		Debug:          true,
		Logger:         log.New(log.Writer(), "", log.LstdFlags),
		ExceptedResult: true,
		Timeout:        5 * time.Second,
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
	err = app.Listen(":3000")
	if err != nil {
		return
	}
}
```