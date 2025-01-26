// Copyright 2018 The Ebiten Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed assets
var commonAssets embed.FS

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	x           float64
	y           float64
	Cursor      *Object
	CursorHover *Object
	Statue      *Object
	IsHovering  bool
	Background  *Object
}

type Object struct {
	x, y          float64
	img           *ebiten.Image
	drawDebugRect bool
}

func newObject(img image.Image, x, y float64) *Object {
	return &Object{
		x:   x,
		y:   y,
		img: ebiten.NewImageFromImage(img),
	}
}

func (o *Object) Position() (float64, float64) {
	return o.x, o.y
}

func (o *Object) SetPosition(x, y float64) {
	o.x = x
	o.y = y
}

func (o *Object) HitBy(other *Object) bool {
	rect := image.Rect(int(o.x), int(o.y), int(o.x)+o.img.Bounds().Dx(), int(o.y)+o.img.Bounds().Dy())
	otherRect := image.Rect(int(other.x), int(other.y), int(other.x)+other.img.Bounds().Dx(), int(other.y)+other.img.Bounds().Dy())
	return rect.Overlaps(otherRect)
}

func (o *Object) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	options.GeoM.Translate(o.x, o.y)
	//options.GeoM.Scale(4.2, 4.2)

	screen.DrawImage(o.img, options)

	if o.drawDebugRect {
		ebitenutil.DrawRect(screen, o.x, o.y, float64(o.img.Bounds().Dx()), float64(o.img.Bounds().Dy()), image.White)
	}
}

func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()
	if x < 0 || y < 0 {
		return nil
	}
	if x >= screenWidth || y >= screenHeight {
		return nil
	}
	g.x = float64(x)
	g.y = float64(y)

	g.IsHovering = g.Statue.HitBy(g.Cursor)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	scaleGeo := ebiten.GeoM{}
	scaleGeo.Scale(4.2, 4.2)
	opBackground := ebiten.DrawImageOptions{
		GeoM: scaleGeo,
	}
	g.Background.Draw(screen, &opBackground)

	g.Statue.Draw(screen, &ebiten.DrawImageOptions{})

	posCursorX := g.x - float64(g.Cursor.img.Bounds().Dx())/2
	posCursorY := g.y - float64(g.Cursor.img.Bounds().Dy())/2

	g.Cursor.SetPosition(posCursorX, posCursorY)
	if g.IsHovering {
		g.CursorHover.SetPosition(posCursorX, posCursorY)
		g.CursorHover.Draw(screen, &ebiten.DrawImageOptions{})
	} else {
		g.Cursor.Draw(screen, &ebiten.DrawImageOptions{})
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("DEBUG MESSAGES: %t, %d, %d", g.IsHovering, g.x, g.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewObjectFromSprite(blob []byte, x, y float64) *Object {
	img, _, err := image.Decode(bytes.NewReader(blob))
	if err != nil {
		log.Fatal(err)
	}
	return newObject(img, x, y)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {

	cursorAsset, err := commonAssets.ReadFile("assets/cursor.png")
	checkErr(err)
	cursor := NewObjectFromSprite(cursorAsset, 0, 0)

	CursorHoverAsset, err := commonAssets.ReadFile("assets/cursorhover.png")
	checkErr(err)
	cursorHover := NewObjectFromSprite(CursorHoverAsset, 0, 0)

	ApexObjectAsset, err := commonAssets.ReadFile("assets/apex/object.png")
	checkErr(err)
	statue := NewObjectFromSprite(ApexObjectAsset, 100, 100)

	ApexBackgroundAsset, err := commonAssets.ReadFile("assets/apex/bg.png")
	checkErr(err)
	background := NewObjectFromSprite(ApexBackgroundAsset, 0, 0)

	g := &Game{x: 0.0, y: 0.0, Statue: statue, Cursor: cursor, CursorHover: cursorHover, Background: background}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Starbound Storyline experience")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
