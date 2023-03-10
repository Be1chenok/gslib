package gslib

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type programID uint32
type shaderID uint32

type Shader struct {
	programID        programID
	vertexPath       string
	fragmentPath     string
	vertexModified   time.Time
	fragmentModified time.Time
}

func NewShader(vertexPath, fragmentPath string) (*Shader, error) {
	id, err := createProgram(vertexPath, fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create program: %v", err)
	}
	vertexModTime, err := getModifiedTime(vertexPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check vertex shader for changes: %v", err)
	}
	fragmentModTime, err := getModifiedTime(fragmentPath)
	if err != nil {
		return nil, fmt.Errorf("failed to check fragment shader for changes: %v", err)
	}
	result := &Shader{programID(id), vertexPath, fragmentPath, vertexModTime, fragmentModTime}

	return result, nil
}

func createProgram(vertPath, fragPath string) (uint32, error) {
	vert, err := loadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to load vertex shader: %v", err)
	}
	frag, err := loadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, fmt.Errorf("failed to load fragment shader: %v", err)
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link shader program: %v", log)
	}
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))

	return shaderProgram, nil
}

func loadShader(path string, shaderType uint32) (shaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("failed to read shader file: %v", err)
	}
	shaderFileStr := string(shaderFile)
	shaderId, err := createShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, fmt.Errorf("failed to create shader: %v", err)
	}

	return shaderId, nil
}

func createShader(shaderSource string, shaderType uint32) (shaderID, error) {
	shaderId := gl.CreateShader(shaderType)
	shaderSource = shaderSource + "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}
	return shaderID(shaderId), nil
}

func (shader *Shader) Use() {
	useProgram(shader)
}

func useProgram(shader *Shader) {
	gl.UseProgram(uint32(shader.programID))
}

func (shader *Shader) CheckShaderForChanges() error {
	vertexModTime, err := getModifiedTime(shader.vertexPath)
	if err != nil {
		return fmt.Errorf("failed to check vertex shader for changes: %v", err)
	}
	fragmentModTime, err := getModifiedTime(shader.fragmentPath)
	if err != nil {
		return fmt.Errorf("failed to check fragment shader for changes: %v", err)
	}
	if !vertexModTime.Equal(shader.vertexModified) ||
		!fragmentModTime.Equal(shader.fragmentModified) {
		id, err := createProgram(shader.vertexPath, shader.fragmentPath)
		if err != nil {
			fmt.Printf("failed to create shader program: %v", err)
		}
		gl.DeleteShader(uint32(shader.programID))
		shader.programID = programID(id)
	}
	return nil
}

func getModifiedTime(filepath string) (time.Time, error) {
	file, err := os.Stat(filepath)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get the modified time: %v", err)
	}
	return file.ModTime(), nil
}

func (shader *Shader) GetUniformLocation(name string) int32 {
	return gl.GetUniformLocation(uint32(shader.programID), gl.Str(name+"\x00"))
}

func (shader *Shader) SetFloat(name string, f float32) {
	gl.Uniform1f(shader.GetUniformLocation(name), f)
}

func (shader *Shader) SetInt(name string, i int32) {
	gl.Uniform1i(shader.GetUniformLocation(name), i)
}
