package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

func TestGenerateSprite(t *testing.T) {
	textures := make([]Texture, 0)

	for i := 0; i < 20; i++ {
		size := rand.Intn(32) + 32
		pixelRatio := i%2 + 1
		name := fmt.Sprintf("%d", i)

		img := image.NewRGBA(image.Rect(0, 0, size, size))

		c := color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}

		draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

		addLabel(img, 2, 10, fmt.Sprintf("%s@%dx", name, pixelRatio))
		addLabel(img, 2, 20, fmt.Sprintf("%d,%d", size, size))

		buf := bytes.NewBuffer(nil)
		png.Encode(buf, img)

		textures = append(textures, &TextureTest{
			SymbolName: name,
			Width:      size,
			Height:     size,
			Binary:     buf.Bytes(),
			PixelRatio: pixelRatio,
		})
	}

	{
		_, img, err := GenerateSprite(textures, 2, true)
		require.NoError(t, err)
		SavePNG("./testdata/test@2x.png", img)
	}

	{
		_, img, err := GenerateSprite(textures, 1, true)
		require.NoError(t, err)
		SavePNG("./testdata/test.png", img)
	}
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}

	d.DrawString(label)
}

type TextureTest struct {
	SymbolName string
	Width      int
	Height     int
	Binary     []byte
	PixelRatio int
}

func (m *TextureTest) TexturePixelRatio() int {
	return m.PixelRatio
}

func (m *TextureTest) TextureName() string {
	return m.SymbolName
}

func (m *TextureTest) TextureWidth() int {
	return m.Width
}

func (m *TextureTest) TextureHeight() int {
	return m.Height
}

func (m *TextureTest) TextureImage() image.Image {
	img, _ := png.Decode(bytes.NewBuffer(m.Binary))
	return img
}

func SavePNG(path string, img image.Image) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
