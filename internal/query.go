package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Joffref/opa-middleware/config"
	"github.com/open-policy-agent/opa/rego"
	"net/http"
)

// QueryPolicy is a helper function to query a local policy evaluation.
// It takes the request, the configuration and the bind variables.
func QueryPolicy(r *http.Request, cfg *config.Config, bind map[string]interface{}) (bool, error) {
	ctx, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
	defer cancel()
	q, err := rego.New(
		rego.Query(cfg.Query),
		rego.Module("policy.rego", cfg.Policy),
	).PrepareForEval(ctx)
	result, err := q.Eval(ctx, rego.EvalInput(bind))
	if err != nil {
		return false, err
	}
	return result.Allowed(), nil
}

// QueryURL is a helper function to query the policy engine using its URL.
// It takes the request, the configuration and the bind variables.
func QueryURL(r *http.Request, cfg *config.Config, bind map[string]interface{}) (bool, error) {
	_, cancel := context.WithTimeout(r.Context(), cfg.Timeout)
	defer cancel()
	jsonBind, err := json.Marshal(bind)
	if err != nil {
		return false, err
	}
	u, err := buildURL(cfg)
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(jsonBind))
	if err != nil {
		return false, err
	}
	req.Header, err = buildHeaders(r, cfg)
	if err != nil {
		return false, err
	}
	post, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer post.Body.Close()
	result := make(map[string]interface{})
	err = json.NewDecoder(post.Body).Decode(&result)
	if err != nil {
		return false, err
	}
	return result["result"].(bool), nil
}
