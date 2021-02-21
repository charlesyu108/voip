package main

import (
	"bytes"
)

type RingBuffer struct {
	internalBuffer *bytes.Buffer
}

func (ring *RingBuffer) Write(p []byte) (n int, err error) {
	toAdd := len(p)
	size := ring.internalBuffer.Len()
	capacity := ring.internalBuffer.Cap()

	if toAdd+size >= capacity {
		unread := capacity - size
		ring.internalBuffer.Truncate(unread)
	}

	return ring.internalBuffer.Write(p)

}

func (ring *RingBuffer) Read(p []byte) (n int, err error) {
	return ring.internalBuffer.Read(p)
}

func NewRingBuffer(capacity int) *RingBuffer {
	internalBuffer := new(bytes.Buffer)
	internalBuffer.Grow(capacity)

	return &RingBuffer{
		internalBuffer: internalBuffer,
	}
}
