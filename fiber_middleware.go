package opamiddleware

import (
	"errors"
	"github.com/Joffref/opa-middleware/config"
	"github.com/Joffref/opa-middleware/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"net/http"
)

// FiberInputCreationMethod is the method that is used to create the input for the policy.
type FiberInputCreationMethod func(c *fiber.Ctx) (map[string]interface{}, error)

type FiberMiddleware struct {
	// Config is the configuration for the middleware.
	Config *config.Config
	// InputCreationMethod is a function that returns the value to be sent to the OPA server.
	InputCreationMethod FiberInputCreationMethod `json:"binding_method,omitempty"`
}

// NewFiberMiddleware is the constructor for the opa fiber middleware.
func NewFiberMiddleware(cfg *config.Config, input FiberInputCreationMethod) (*FiberMiddleware, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	if input == nil && cfg.InputCreationMethod == nil {
		return nil, errors.New("[opa-middleware-fiber] InputCreationMethod must be provided")
	}
	return &FiberMiddleware{
		Config:              cfg,
		InputCreationMethod: input,
	}, nil
}

// Use returns the handler for the middleware that is used by fiber to evaluate the request against the policy.
func (g *FiberMiddleware) Use() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if g.Config.Debug {
			g.Config.Logger.Printf("[opa-middleware-fiber] Request received")
		}
		result, err := g.query(c)
		if err != nil {
			if g.Config.Debug {
				g.Config.Logger.Printf("[opa-middleware-fiber] Error: %s", err.Error())
			}
			c.Status(http.StatusInternalServerError)

		}
		if g.Config.Debug {
			g.Config.Logger.Printf("[opa-middleware-fiber] Result: %t", result)
		}
		if result != g.Config.ExceptedResult {
			c.Status(g.Config.DeniedStatusCode)
			return errors.New("[opa-middleware-fiber] Access denied")
		}
		err = c.Next()
		if err != nil {
			return err
		}
		return nil
	}
}

func (g *FiberMiddleware) query(c *fiber.Ctx) (bool, error) {
	bind, err := g.InputCreationMethod(c)
	if err != nil {
		return !g.Config.ExceptedResult, err
	}
	if g.Config.URL != "" {
		input := make(map[string]interface{})
		input["input"] = bind
		return internal.QueryURL(transformFastHTTP(c.Context()), g.Config, input)
	}
	return internal.QueryPolicy(transformFastHTTP(c.Context()), g.Config, bind)
}

func transformFastHTTP(ctx *fasthttp.RequestCtx) *http.Request {
	req := &http.Request{
		Header: make(http.Header),
	}
	headers := make(map[string][]string)
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		headers[string(key)] = append(headers[string(key)], string(value))
	})
	for k, v := range headers {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	req = req.WithContext(ctx)
	return req
}
