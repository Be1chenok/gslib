package gslib

import "github.com/go-gl/gl/v4.6-core/gl"

type VAOID struct {
	id uint32
}

func GenVertexArray(n int32) *VAOID {
	var VAO uint32
	gl.GenVertexArrays(n, &VAO)
	return &VAOID{VAO}
}

func (vao *VAOID) BindVertexArray() {
	gl.BindVertexArray(vao.id)
}

func (vao *VAOID) UnBindVertexArray() {
	gl.BindVertexArray(0)
}
