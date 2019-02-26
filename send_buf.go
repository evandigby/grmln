package grmln

import (
	"bytes"
	"sync"
)

type sendBufferPool struct {
	mimeType     string
	mimeTypeLen  int
	mimeTypeData []byte
	truncateLen  int
	pool         sync.Pool
}

func newSendBufferPool(mimeType string) *sendBufferPool {
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

	return &p
}

func (p *sendBufferPool) mimeTypeDataCopy() []byte {
	buf := make([]byte, len(p.mimeTypeData))
	copy(buf, p.mimeTypeData)
	return buf
}

func (p *sendBufferPool) get() *bytes.Buffer {
	buf := p.pool.Get().(*bytes.Buffer)
	p.Reset(buf)
	return buf
}

func (p sendBufferPool) Reset(buf *bytes.Buffer) {
	buf.Truncate(p.truncateLen)
}

func (p *sendBufferPool) put(buf *bytes.Buffer) {
	p.pool.Put(buf)
}
