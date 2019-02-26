package grmln

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

// EvalArgs args required for eval ops
type EvalArgs struct {
	OpArgs
	Gremlin                   string                 `json:"gremlin"`
	Bindings                  map[string]interface{} `json:"bindings,omitempty"`
	Language                  string                 `json:"language"`
	Aliases                   map[string]string      `json:"aliases,omitempty"`
	ScriptEvaluationTimeoutMS int64                  `json:"scriptEvaluationTimeout"`
}

