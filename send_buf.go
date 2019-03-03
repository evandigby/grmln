package grmln

import (
	"bytes"
	"sync"
)

var sendBufferPools = map[string]*sendBufferPool{}
var sendBufferPoolsMutex sync.RWMutex

type sendBufferPool struct {
	mimeType     string
	mimeTypeLen  int
	mimeTypeData []byte
	truncateLen  int
	pool         sync.Pool
}

func newSendBufferPool(mimeType string) *sendBufferPool {
	// These should be shared by mime type across all users in the same process
	sendBufferPoolsMutex.RLock()
	if p, ok := sendBufferPools[mimeType]; ok {
		sendBufferPoolsMutex.RUnlock()
		return p
	}
	sendBufferPoolsMutex.RUnlock()

	sendBufferPoolsMutex.Lock()
	defer sendBufferPoolsMutex.Unlock()

	// Double check in case it was created between RUnlock and Lock
	if p, ok := sendBufferPools[mimeType]; ok {
		return p
	}

	p := sendBufferPool{
		mimeType:    mimeType,
		mimeTypeLen: len(mimeType),
	}

	p.mimeTypeData = append(
		[]byte{
			byte(p.mimeTypeLen),
		},
		[]byte(p.mimeType)...,
	)

	p.truncateLen = len(p.mimeTypeData)

	p.pool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(p.mimeTypeDataCopy())
		},
	}

	sendBufferPools[mimeType] = &p
	return &p
}

func (p *sendBufferPool) mimeTypeDataCopy() []byte {
	buf := make([]byte, len(p.mimeTypeData))
	copy(buf, p.mimeTypeData)
	return buf
}

func (p *sendBufferPool) get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	p.reset(buf)
	return buf
}

func (p sendBufferPool) reset(buf *bytes.Buffer) {
	buf.Truncate(p.truncateLen)
}

func (p *sendBufferPool) put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}
