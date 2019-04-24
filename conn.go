package grmln

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

// OnResponse callback when a partial or complete response is received
type OnResponse func(resp *Response)

var noopOnResponse = func(resp *Response) {}

// Conn is a gremlin server connection
type Conn struct {
	mimeType string
	addr     string
	userName string
	password string
	headers  http.Header
	authArgs AuthenticationArgs
	ws       *websocket.Conn

	requestMutex sync.Mutex

	sendBufferPool *sendBufferPool
}

// SASL calculates sasl authentication args
func SASL(userName, password string) AuthenticationArgs {
	sasl := []byte{0}
	sasl = append(sasl, []byte(userName)...)
	sasl = append(sasl, 0)
	sasl = append(sasl, []byte(password)...)

	return AuthenticationArgs{
		SASL:          base64.StdEncoding.EncodeToString(sasl),
		SASLMechanism: "PLAIN",
	}
}

// Dial dials addresses
func Dial(ctx context.Context, addr, mimeType, userName, password string, headers http.Header) (*Conn, error) {
	dialer := websocket.Dialer{}

	ws, _, err := dialer.DialContext(ctx, addr, headers)
	if err != nil {
		return nil, err
	}

	return &Conn{
		mimeType:       mimeType,
		addr:           addr,
		userName:       userName,
		password:       password,
		headers:        headers,
		authArgs:       SASL(userName, password),
		ws:             ws,
		sendBufferPool: newSendBufferPool(mimeType),
	}, nil
}

// ProcessRequest can process a raw gremlin request
func (c *Conn) ProcessRequest(ctx context.Context, r Request, onResponse ...OnResponse) error {
	c.requestMutex.Lock()
	defer c.requestMutex.Unlock()

	err := c.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	return c.readResponse(ctx, onResponse...)
}

func (c *Conn) sendRequest(ctx context.Context, r Request) error {
	buf := c.sendBufferPool.get()
	defer c.sendBufferPool.put(buf)

	err := json.NewEncoder(buf).Encode(r)
	if err != nil {
		return err
	}

	dl, ok := ctx.Deadline()
	if ok {
		c.ws.SetWriteDeadline(dl)
	}

	return c.ws.WriteMessage(websocket.BinaryMessage, buf.Bytes())
}

func (c *Conn) readResponse(ctx context.Context, onResponse ...OnResponse) error {
	dl, ok := ctx.Deadline()
	if ok {
		c.ws.SetReadDeadline(dl)
	}
	for {
		var resp Response

		err := c.ws.ReadJSON(&resp)
		if err != nil {
			return err
		}

		if err := resp.Err(); err != nil {
			if !IsAuthenticate(err) {
				return err
			}

			return c.ProcessRequest(ctx, NewRequest(resp.RequestID, processorDefault, opAuthentication, c.authArgs), onResponse...)
		}

		for _, or := range onResponse {
			or(&resp)
		}

		if !resp.IsPartial() {
			return nil
		}
	}
}

// Close closes the connection (including the underlying websocket)
func (c *Conn) Close() error {
	c.requestMutex.Lock()
	defer c.requestMutex.Unlock()
	return c.ws.Close()
}
