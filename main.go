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
	_ "embed"
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

//go:embed cursor.png
var Cursor_png []byte

//go:embed apexstatue1.png
var Apexstatue1_png []byte

const (
	screenWidth  = 800
	screenHeight = 600
)

type Game struct {
	x      float64
	y      float64
	Cursor *Object
	Statue *Object
}

type Object struct {
	x, y          float64
	img           *ebiten.Image
	drawDebugRect bool
}

func NewObject(img image.Image, x, y float64) *Object {
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
	fmt.Println(rect, otherRect)
	return rect.Overlaps(otherRect)
}

func (o *Object) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.x, o.y)
	screen.DrawImage(o.img, op)

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

	fmt.Println(g.Statue.HitBy(g.Cursor))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	//g.Statue.drawDebugRect = true
	g.Statue.Draw(screen)

	//g.Cursor.drawDebugRect = true
	g.Cursor.SetPosition(g.x-float64(g.Cursor.img.Bounds().Dx())/2, g.y-float64(g.Cursor.img.Bounds().Dy())/2)
	g.Cursor.Draw(screen)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Move the redddd point by mouse wheeaal\n(%0.2f, %0.2f)", g.x, g.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(Cursor_png))
	if err != nil {
		log.Fatal(err)
	}
	cursor := NewObject(img, 0, 0)

	img, _, err = image.Decode(bytes.NewReader(Apexstatue1_png))
	if err != nil {
		log.Fatal(err)
	}
	statue := NewObject(img, 100, 100)

	g := &Game{x: 0.0, y: 0.0, Statue: statue, Cursor: cursor}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Starbound Storyline experience")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
