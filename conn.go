package grmln

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// OnResponse callback when a partial or complete response is received
type OnResponse func(resp *Response)

// Conn is a gremlin server connection
type Conn struct {
	ws *websocket.Conn

	// DefaultScriptEvaluationTimeout is the default script evaluation timeout. Defaults to 3000ms
	DefaultScriptEvaluationTimeout time.Duration

	// DefaultEvalLanguage is the default eval language for use. Defaults to "gremlin-groovy"
	DefaultEvalLanguage string

	// DefaultBatchSize is the default size of batched responses. 0 uses server default
	DefaultBatchSize int
}

// Dial dials addresses
func Dial(addr string) (*Conn, error) {
	dialer := websocket.Dialer{}

	ws, _, err := dialer.Dial(addr, http.Header{})
	if err != nil {
		return nil, err
	}

	return &Conn{
		ws:                             ws,
		DefaultScriptEvaluationTimeout: 3000 * time.Millisecond,
		DefaultEvalLanguage:            LanguageGremlinGroovy,
		DefaultBatchSize:               0, // 0 uses server default
	}, nil
}

const (
	opEval = "eval"
)

// Eval languages
const (
	LanguageGremlinGroovy = "gremlin-groovy"
)

// Eval evaluates a gremlin statement
func (c *Conn) Eval(args EvalArgs, onResponse OnResponse) error {
	r := request{
		RequestID: uuid.New().String(),
		Operation: opEval,
		Arguments: args,
	}

	return c.processRequest(r, onResponse)
}

// EvalDefault is a helper that calls Eval using the default argument values
func (c *Conn) EvalDefault(gremlin string, onResponse OnResponse) error {
	return c.Eval(
		EvalArgs{
			OpArgs: OpArgs{
				BatchSize: c.DefaultBatchSize,
			},
			Gremlin:                   gremlin,
			Language:                  c.DefaultEvalLanguage,
			ScriptEvaluationTimeoutMS: int64(c.DefaultScriptEvaluationTimeout / time.Millisecond),
		},
		onResponse,
	)
}

func (c *Conn) processRequest(r request, onResponse OnResponse) error {
	err := c.sendRequest(r)
	if err != nil {
		return err
	}

	return c.readResponse(onResponse)
}

func (c *Conn) sendRequest(r request) error {
	const appType = "application/vnd.gremlin-v2.0+json"
	msg := []byte{
		byte(len(appType)),
	}

	rdata, err := json.Marshal(r)
	if err != nil {
		return err
	}

	msg = append(msg, []byte(appType)...)
	msg = append(msg, rdata...)

	return c.ws.WriteMessage(websocket.BinaryMessage, msg)
}

func (c *Conn) readResponse(onResponse OnResponse) error {
	for {
		var resp Response
		err := c.ws.ReadJSON(&resp)
		if err != nil {
			return err
		}

		if err := resp.Err(); err != nil {
			return err
		}

		onResponse(&resp)

		if !resp.IsPartial() {
			return nil
		}
	}
}

// Close closes the connection (including the underlying websocket)
func (c *Conn) Close() error {
	return c.ws.Close()
}
