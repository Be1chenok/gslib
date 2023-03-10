package gslib

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var fps uint32
var lastFpsTick uint32

func VSync(activate bool) {
	switch activate {
	case true:
		err := sdl.GLSetSwapInterval(1)
		if err != nil {
			panic(err)
		}
	case false:
		err := sdl.GLSetSwapInterval(0)
		if err != nil {
			panic(err)
		}
	}
}

func CalculateFps(window *sdl.Window, title string) {
	currentTick := sdl.GetTicks()
	if currentTick > lastFpsTick+1000 { // обновление FPS каждую 1 секунду
		window.SetTitle(fmt.Sprintf("%s (FPS: %d)", title, fps))
		lastFpsTick = currentTick
		fps = 0
	}
	fps++
}

func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

func SetContextVersion(major, minor int) {
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, major)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, minor)
}
