package grmln

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	opEval = "eval"
)

// Eval languages
const (
	LanguageGremlinGroovy = "gremlin-groovy"
)

// RequestProcessor processes raw requests
type RequestProcessor interface {
	ProcessRequest(ctx context.Context, r Request, onResponse OnResponse) error
}

// NewOperator creates a new gremlin operator
func NewOperator(p RequestProcessor) *Operator {
	return &Operator{
		p:                              p,
		DefaultScriptEvaluationTimeout: 3000 * time.Millisecond,
		DefaultEvalLanguage:            LanguageGremlinGroovy,
		DefaultBatchSize:               0, // 0 uses server default
	}
}

// Operator is a helper to build gremlin operations
type Operator struct {
	p RequestProcessor

	// DefaultScriptEvaluationTimeout is the default script evaluation timeout. Defaults to 3000ms
	DefaultScriptEvaluationTimeout time.Duration

	// DefaultEvalLanguage is the default eval language for use. Defaults to "gremlin-groovy"
	DefaultEvalLanguage string

	// DefaultBatchSize is the default size of batched responses. 0 uses server default
	DefaultBatchSize int
}

// Eval evaluates a gremlin statement
func (o *Operator) Eval(ctx context.Context, args EvalArgs, onResponse OnResponse) error {
	r := Request{
		RequestID: uuid.New().String(),
		Operation: opEval,
		Arguments: args,
	}

	return o.p.ProcessRequest(ctx, r, onResponse)
}

// EvalDefault is a helper that calls Eval using the default argument values
func (o *Operator) EvalDefault(ctx context.Context, gremlin string, onResponse OnResponse) error {
	return o.Eval(
		ctx,
		EvalArgs{
			OpArgs: OpArgs{
				BatchSize: o.DefaultBatchSize,
			},
			Gremlin:                   gremlin,
			Language:                  o.DefaultEvalLanguage,
			ScriptEvaluationTimeoutMS: int64(o.DefaultScriptEvaluationTimeout / time.Millisecond),
		},
		onResponse,
	)
}
