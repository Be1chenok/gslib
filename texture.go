package gslib

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Texture struct {
	textureID uint32
	target    uint32
	texUnit   uint32
}

func NewTexture2DFromFile(file string, wrapR, wrapS, n int32) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %v", err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file: %v", err)
	}
	return newTexture2D(img, wrapR, wrapS, n)
}

func newTexture2D(img image.Image, wrapR, wrapS, n int32) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride, only 32-bit colors supported")
	}

	id := genTextures(n)
	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := Texture{
		textureID: id,
		target:    target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	gl.GenerateMipmap(texture.textureID)

	return &texture, nil
}

func newTexture3D(wrapS, wrapT, wrapR, n int32) {

	id := genTextures(n)
	target := uint32(gl.TEXTURE_3D)
	internalFmt := int32(gl.RGBA8)
	format := uint32(gl.RGBA)
	width := int32(64)
	height := int32(64)
	depth := int32(64)
	pixType := uint32(gl.UNSIGNED_BYTE)
	data := make([]byte, width*height*depth*4)
	dataPtr := gl.Ptr(data)

	texture := Texture{
		textureID: id,
		target:    target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_T, wrapT)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage3D(target, 0, internalFmt, width, height, depth, 0, format, pixType, dataPtr)
}

func genTextures(n int32) uint32 {
	var id uint32
	gl.GenTextures(n, &id)
	return id
}

func (tex *Texture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.textureID)
	tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetTexture2D(uniformLocation int32) error {
	if tex.texUnit == 0 {
		return fmt.Errorf("texture not bound")
	}
	gl.Uniform1i(uniformLocation, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}
