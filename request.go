package grmln

import (
	"github.com/google/uuid"
)

// NewRequest creates a new request
func NewRequest(id, processor, operation string, arguments interface{}) Request {
	if id == "" {
		id = uuid.New().String()
	}
	return Request{
		RequestID: id,
		Processor: processor,
		Operation: operation,
		Arguments: arguments,
	}
}

// Request represents a raw gremlin request
type Request struct {
	RequestID string      `json:"requestId"`
	Operation string      `json:"op"`
	Processor string      `json:"processor"`
	Arguments interface{} `json:"args"`
}

// OpArgs are args available to all operations
type OpArgs struct {
	BatchSize int `json:"batchSize,omitempty"`
}

type Bindings map[string]interface{}

// EvalArgs args required for eval ops
type EvalArgs struct {
	OpArgs
	Gremlin                   string            `json:"gremlin"`
	Bindings                  Bindings          `json:"bindings,omitempty"`
	Language                  string            `json:"language"`
	Aliases                   map[string]string `json:"aliases,omitempty"`
	ScriptEvaluationTimeoutMS int64             `json:"scriptEvaluationTimeout"`
}

// AuthenticationArgs args required for authentication ops
type AuthenticationArgs struct {
	SASL          string `json:"sasl"`
	SASLMechanism string `json:"saslMechanism,omitempty"`
}

// SessionArgs are args for a session
type SessionArgs struct {
	Session string `json:"session"`
}

// TransactionEvalArgs are args specific to evals within a session with transaction management
type TransactionEvalArgs struct {
	EvalArgs
	ManageTransaction bool `json:"manageTransaction"`
}

// SessionEvalArgs are args specific to evals within a session
type SessionEvalArgs struct {
	SessionArgs
	TransactionEvalArgs
}

// SessionAuthenticationArgs are args specific to authentication within a session
type SessionAuthenticationArgs struct {
	SessionArgs
	AuthenticationArgs
}

type CloseArgs struct {
	Force bool `json:"force,omitempty"`
}

// SessionCloseArgs args for session close
type SessionCloseArgs struct {
	SessionArgs
	CloseArgs
}
