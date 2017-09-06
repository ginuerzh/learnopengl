package texture

import (
	"fmt"
	"image"
	// jpeg format support
	"image/draw"
	_ "image/jpeg"
	// png format support
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture interface {
	SetParameter(uint32, interface{}) error
	Load(string) (*image.RGBA, error)
	Use()
}

type Texture2D struct {
	ID     uint32
	params map[uint32]interface{}
}

func NewTexture2D() Texture {
	var id uint32
	gl.GenTextures(1, &id)
	return &Texture2D{
		ID:     id,
		params: make(map[uint32]interface{}),
	}
}

func (texture *Texture2D) SetParameter(name uint32, param interface{}) error {
	if texture.params == nil {
		texture.params = make(map[uint32]interface{})
	}
	texture.params[name] = param

	switch v := param.(type) {
	case int:
		gl.TexParameteri(gl.TEXTURE_2D, name, int32(v))
	case int32:
		gl.TexParameteri(gl.TEXTURE_2D, name, v)
	case float32:
		gl.TexParameterf(gl.TEXTURE_2D, name, v)
	case float64:
		gl.TexParameterf(gl.TEXTURE_2D, name, float32(v))
	default:
		return fmt.Errorf("unsupported type for %d", name)
	}
	return nil
}

func (texture *Texture2D) Load(textureFile string) (*image.RGBA, error) {
	f, err := os.Open(textureFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	rgba := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride %d", rgba.Stride)
	}
	draw.Draw(rgba, rgba.Bounds(), src, b.Min, draw.Src)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	return rgba, nil
}

func (texture *Texture2D) Use() {
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)
}
