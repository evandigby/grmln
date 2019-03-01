package grmln

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	opEval           = "eval"
	opAuthentication = "authentication"
	opClose          = "close"
)

const (
	processorDefault = ""
	processorSession = "session"
)

// Eval languages
const (
	LanguageGremlinGroovy = "gremlin-groovy"
)

// RequestProcessor processes raw requests
type RequestProcessor interface {
	ProcessRequest(ctx context.Context, r Request, onResponse OnResponse) error
}

// OperatorConfig is configuration required by all operators
type OperatorConfig struct {
	// DefaultScriptEvaluationTimeout is the default script evaluation timeout. Defaults to 3000ms
	DefaultScriptEvaluationTimeout time.Duration

	// DefaultEvalLanguage is the default eval language for use. Defaults to "gremlin-groovy"
	DefaultEvalLanguage string

	// DefaultBatchSize is the default size of batched responses. 0 uses server default
	DefaultBatchSize int
}

func (o OperatorConfig) evalArgs(gremlin string) EvalArgs {
	return EvalArgs{
		OpArgs: OpArgs{
			BatchSize: o.DefaultBatchSize,
		},
		Gremlin:                   gremlin,
		Language:                  o.DefaultEvalLanguage,
		ScriptEvaluationTimeoutMS: int64(o.DefaultScriptEvaluationTimeout / time.Millisecond),
	}
}

// NewOperator creates a new gremlin operator
func NewOperator(p RequestProcessor) *Operator {
	return &Operator{
		p: p,
		OperatorConfig: OperatorConfig{
			DefaultScriptEvaluationTimeout: 3000 * time.Millisecond,
			DefaultEvalLanguage:            LanguageGremlinGroovy,
			DefaultBatchSize:               0, // 0 uses server default
		},
	}
}

// Operator is a helper to build gremlin operations
type Operator struct {
	p RequestProcessor

	OperatorConfig
}

// NewSession creates a new session
func (o *Operator) NewSession() *SessionOperator {
	return &SessionOperator{
		p:              o.p,
		OperatorConfig: o.OperatorConfig,
		session:        uuid.New().String(),
	}
}

// Eval evaluates a gremlin statement
func (o *Operator) Eval(ctx context.Context, args EvalArgs, onResponse OnResponse) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorDefault, opEval, args), onResponse)
}

// EvalDefault is a helper that calls Eval using the default argument values
func (o *Operator) EvalDefault(ctx context.Context, gremlin string, onResponse OnResponse) error {
	return o.Eval(ctx, o.evalArgs(gremlin), onResponse)
}

// Authentication evaluates an authentication statement
func (o *Operator) Authentication(ctx context.Context, args AuthenticationArgs, onResponse OnResponse) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorDefault, opAuthentication, args), onResponse)
}

// SessionOperator is a helper to build gremlin operations
type SessionOperator struct {
	p RequestProcessor

	OperatorConfig

	session string
}

func (o *SessionOperator) sessionArgs() SessionArgs {
	return SessionArgs{
		Session: o.session,
	}
}

// Eval evaluates a gremlin statement
func (o *SessionOperator) Eval(ctx context.Context, args TransactionEvalArgs, onResponse OnResponse) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorSession, opEval, SessionEvalArgs{
		SessionArgs:         o.sessionArgs(),
		TransactionEvalArgs: args,
	}), onResponse)
}

// EvalDefault is a helper that calls Eval using the default argument values
func (o *SessionOperator) EvalDefault(ctx context.Context, gremlin string, onResponse OnResponse) error {
	return o.Eval(
		ctx,
		TransactionEvalArgs{
			EvalArgs: o.evalArgs(gremlin),
		},
		onResponse,
	)
}

// Authentication evaluates an authentication statement
func (o *SessionOperator) Authentication(ctx context.Context, args AuthenticationArgs, onResponse OnResponse) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorSession, opAuthentication,
		SessionAuthenticationArgs{
			SessionArgs:        o.sessionArgs(),
			AuthenticationArgs: args,
		}), onResponse)
}

// Close closes the session
func (o *SessionOperator) Close(ctx context.Context, args CloseArgs) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorSession, opAuthentication,
		SessionCloseArgs{
			SessionArgs: o.sessionArgs(),
			CloseArgs:   args,
		}), noopOnResponse)
}

// CloseDefault closes the session with default optionss
func (o *SessionOperator) CloseDefault(ctx context.Context) error {
	return o.p.ProcessRequest(ctx, NewRequest(processorSession, opAuthentication,
		SessionCloseArgs{
			SessionArgs: o.sessionArgs(),
		}), noopOnResponse)
}

