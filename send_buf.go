package grmln

import (
	"bytes"
	"sync"
)

const (
	appType    = "application/vnd.gremlin-v2.0+json"
	appTypeLen = len(appType)
)

var appTypeData = append(
	[]byte{
		byte(appTypeLen),
	},
	[]byte(appType)...,
)

var truncateLen = len(appTypeData)

func appTypeDataCopy() []byte {
	buf := make([]byte, len(appTypeData))
	copy(buf, appTypeData)
	return buf
}

var sendBufPool = sync.Pool{
	New: func() interface{} {
		return sendBuf{bytes.NewBuffer(appTypeDataCopy())}
	},
}

func getSendBuf() sendBuf {
	buf := sendBufPool.Get().(sendBuf)
	buf.Reset()
	return buf
}

type sendBuf struct {
	*bytes.Buffer
}

func (s sendBuf) Close() {
	sendBufPool.Put(s)
}

func (s sendBuf) Reset() {
	s.Buffer.Truncate(truncateLen)
}
