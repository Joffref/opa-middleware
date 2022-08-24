# Open Policy Agent Gin Middleware

This middleware integrates Open Policy Agent (OPA) to your gin app.
You can use it to enforce policies on endpoints.
You can use OPA as local policy engine, or as a remote policy engine.

## Installation

```bash
go get github.com/Joffref/opa-middleware
```

## Usage
### Local policy engine
```go
package main

import (
    "github.com/Joffref/gin-opa-middleware"
    "github.com/gin-gonic/gin"
)

var policy = `
package example.authz

default allow := false

allow {
	input.method == "GET"
}`

func main() {
    r := gin.Default()
	r.Use(opa.Middleware(context.Background(), &opa.Config{
		Policy: policy,
		Query: "data.example.authz.allow",
		InputCreationMethod: func(c *gin.Context) (map[string]interface{}, error) {
			return map[string]interface{}{
				"method": c.Request.Method,
			}, nil
		},
		ExceptedResult:   true,
		DeniedStatusCode: 403,
		Debug:            true,
		Logger:           log.New(gin.DefaultWriter, "[opa] ", log.LstdFlags),
	}))
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello World!")
    })
    r.Run()
}
```
### Remote policy engine
```go
package main

import (
    "github.com/Joffref/gin-opa-middleware"
    "github.com/gin-gonic/gin"
)
func main() {
    r := gin.Default()
    r.Use(opa.Middleware(context.Background(), &opa.Config{
        URL:   "http://localhost:8181",
        Query: "data.example.authz.allow",
        InputCreationMethod: func(c *gin.Context) (map[string]interface{}, error) {
        return map[string]interface{}{
        "method": c.Request.Method,
        }, nil
        },
        ExceptedResult:   true,
        DeniedStatusCode: 403,
        Debug:            true,
        Logger:           log.New(gin.DefaultWriter, "[opa] ", log.LstdFlags),
    }))
    r.GET("/", func(c *gin.Context) {
        c.String(200, "Hello World!")
    })
    err := r.Run(":8080")
    if err != nil {
		return
    }
}
```