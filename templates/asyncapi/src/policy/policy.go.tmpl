// Package policy enthält OPA-WASM Policy-Evaluation für Channels
package policy

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"github.com/open-policy-agent/opa/rego"
)

// EvalPolicy wertet die Policy für given input aus.
func EvalPolicy(ctx context.Context, module string, input interface{}) (bool, error) {
	r := rego.New(
		rego.Query("data."+module+".allow"),
		rego.Load([]string{"policies/"+module+".rego"}, nil),
	)
	rs, err := r.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		log.Error().Err(err).Msg("Policy evaluation failed")
		return false, err
	}
	if len(rs) == 0 {
		return false, errors.New("no policy result")
	}
	allowed, ok := rs[0].Expressions[0].Value.(bool)
	if !ok {
		return false, errors.New("unexpected policy result type")
	}
	return allowed, nil
}
