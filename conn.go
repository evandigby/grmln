package grmln

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// OnResponse callback when a partial or complete response is received
type OnResponse func(resp *Response)

// Conn is a gremlin server connection
type Conn struct {
	mimeType string
	addr     string
	ws       *websocket.Conn

	requestMutex sync.Mutex

	sendBufferPool *sendBufferPool
}

// Dial dials addresses
func Dial(ctx context.Context, addr string, mimeType string) (*Conn, error) {
	dialer := websocket.Dialer{}

	ws, _, err := dialer.DialContext(ctx, addr, http.Header{})
	if err != nil {
		return nil, err
	}

	return &Conn{
		mimeType:       mimeType,
		addr:           addr,
		ws:             ws,
		sendBufferPool: newSendBufferPool(mimeType),
	}, nil
}

// ProcessRequest can process a raw gremlin request
func (c *Conn) ProcessRequest(ctx context.Context, r Request, onResponse OnResponse) error {
	c.requestMutex.Lock()
	defer c.requestMutex.Unlock()

	err := c.sendRequest(ctx, r)
	if err != nil {
		return err
	}

	return c.readResponse(ctx, onResponse)
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

func (c *Conn) readResponse(ctx context.Context, onResponse OnResponse) error {
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
	c.requestMutex.Lock()
	defer c.requestMutex.Unlock()
	return c.ws.Close()
}
