package gslib

import (
	"github.com/go-gl/gl/v4.6-core/gl"
)

type BufferID struct {
	id     uint32
	target uint32
}

type Data interface {
	uint32 | float32
}

func GenBuffer(target uint32, n int32) *BufferID {
	var buffer uint32
	gl.GenBuffers(n, &buffer)
	return &BufferID{buffer, target}
}

func (buffer *BufferID) BindBuffer() {
	gl.BindBuffer(buffer.target, buffer.id)
}

func BufferData[T Data](target uint32, data []T, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}

func (buffer *BufferID) UnBindBuffer() {
	gl.BindBuffer(buffer.target, 0)
}
