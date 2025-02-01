package object

import (
	"bytes"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/hajimehoshi/ebiten/v2"
)

type Object struct {
	X, Y          float64
	Img           *ebiten.Image
	drawDebugRect bool
}

type CursorObject struct {
	*Object
	IsHovering bool
	imgNormal  *ebiten.Image
	imgHover   *ebiten.Image
}

func (c *CursorObject) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	img := c.imgNormal
	if c.IsHovering {
		img = c.imgHover
	}
	c.Img = img
	c.Object.Draw(screen, options)
}

func newObject(img image.Image, x, y float64) *Object {
	return &Object{
		X:   x,
		Y:   y,
		Img: ebiten.NewImageFromImage(img),
	}
}

func (o *Object) Position() (float64, float64) {
	return o.X, o.Y
}

func (o *Object) SetPosition(x, y float64) {
	o.X = x
	o.Y = y
}

func (o *Object) HitBy(other *Object) bool {
	rect := image.Rect(int(o.X), int(o.Y), int(o.X)+o.Img.Bounds().Dx(), int(o.Y)+o.Img.Bounds().Dy())
	otherRect := image.Rect(int(other.X), int(other.Y), int(other.X)+other.Img.Bounds().Dx(), int(other.Y)+other.Img.Bounds().Dy())
	return rect.Overlaps(otherRect)
}

func (o *Object) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	options.GeoM.Translate(o.X, o.Y)

	screen.DrawImage(o.Img, options)

	if o.drawDebugRect {
		ebitenutil.DrawRect(screen, o.X, o.Y, float64(o.Img.Bounds().Dx()), float64(o.Img.Bounds().Dy()), image.White)
	}
}

func NewCursorObject(blobNormal, blobHover []byte) *CursorObject {
	imgNormal, _, err := image.Decode(bytes.NewReader(blobNormal))
	if err != nil {
		log.Fatal(err)
	}
	imgHover, _, err := image.Decode(bytes.NewReader(blobHover))
	if err != nil {
		log.Fatal(err)
	}
	cursor := &CursorObject{Object: newObject(imgNormal, 0, 0), imgNormal: ebiten.NewImageFromImage(imgNormal), imgHover: ebiten.NewImageFromImage(imgHover)}
	return cursor
}

func NewObjectFromSprite(blob []byte, x, y float64) *Object {
	img, _, err := image.Decode(bytes.NewReader(blob))
	if err != nil {
		log.Fatal(err)
	}
	return newObject(img, x, y)
}
