package opamiddleware

import (
	"errors"
	"github.com/Joffref/opa-middleware/config"
	"github.com/Joffref/opa-middleware/internal"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EchoInputCreationMethod func(c echo.Context) (map[string]interface{}, error)

type EchoMiddleware struct {
	Config *config.Config
	InputCreationMethod EchoInputCreationMethod `json:"binding_method,omitempty"`
}

func NewEchoMiddleware(cfg *config.Config, input EchoInputCreationMethod) (*EchoMiddleware, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}
	if input == nil {
		if cfg.InputCreationMethod == nil {
			return nil, errors.New("[opa-middleware-echo] InputCreationMethod must be provided")
		}
		input = func(c echo.Context) (map[string]interface{}, error) {
			bind, err := cfg.InputCreationMethod(c.Request())
			if err != nil {
				return nil, err
			}
			return bind, nil
		}
	}
	return &EchoMiddleware{
		Config:              cfg,
		InputCreationMethod: input,
	}, nil
}

func (e *EchoMiddleware) Use() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if e.Config.Debug {
				e.Config.Logger.Printf("[opa-middleware-echo] Request received")
			}
			result, err := e.query(c)
			if err != nil {
				if e.Config.Debug {
					e.Config.Logger.Printf("[opa-middleware-echo] Error: %s", err.Error())
				}
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
			}
			if e.Config.Debug {
				e.Config.Logger.Printf("[opa-middleware-echo] Result: %t", result)
			}
			if result != e.Config.ExceptedResult {
				return c.JSON(e.Config.DeniedStatusCode, map[string]interface{}{"error": e.Config.DeniedMessage})
			}
			return next(c)
		}
	}
}

func (e *EchoMiddleware) query(c echo.Context) (bool, error) {
	bind, err := e.InputCreationMethod(c)
	if err != nil {
		return !e.Config.ExceptedResult, err
	}
	if e.Config.URL != "" {
		input := make(map[string]interface{})
		input["input"] = bind
		return internal.QueryURL(c.Request(), e.Config, input)
	}
	return internal.QueryPolicy(c.Request(), e.Config, bind)
}