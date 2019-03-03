package grmln

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
)

// TestSendBufReset ensures resetting a send buf will produce the correct results
func TestSendBufReset(t *testing.T) {
	const numTests = 25

	p := newSendBufferPool(DefaultMimeType)

	b := p.get()

	for i := 0; i < numTests; i++ {
		randLen := rand.Int63n(127) + 1
		data := make([]byte, randLen)
		rand.Read(data)

		t.Run(fmt.Sprintf("pass %d", i), func(t *testing.T) {
			expectedResult := append(p.mimeTypeDataCopy(), data...)
			_, err := b.Write(data)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if actualResult := b.Bytes(); bytes.Compare(expectedResult, actualResult) != 0 {
				t.Fatalf("expected %v but got %v", expectedResult, actualResult)
			}
		})

		p.reset(b)
	}
}

func TestMimeTypeDataCopy(t *testing.T) {
	p := newSendBufferPool(DefaultMimeType)

	expected, actual := p.mimeTypeData, p.mimeTypeDataCopy()

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected %v but got %v", expected, actual)
	}
}

func BenchmarkBufPool(b *testing.B) {
	data := make([]byte, 1024)
	rand.Read(data)

	p := newSendBufferPool(DefaultMimeType)
	b.Run("pooled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			func() {
				buf := p.get()
				defer p.put(buf)
				buf.Write(data)
			}()
		}
	})

	b.Run("unpooled", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = append(
				[]byte{
					byte(p.mimeTypeLen),
				},
				append([]byte(p.mimeTypeData), data...)...,
			)
		}
	})

	b.Run("unpooled using copy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = append(p.mimeTypeDataCopy(), data...)
		}
	})
}
