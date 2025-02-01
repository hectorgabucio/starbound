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

type Stage interface {
	GetID() string
	Update() error
	Draw(screen *ebiten.Image)
	Finished() bool
}

type CursorObject struct {
	*Object
	isHovering bool
	imgNormal  *ebiten.Image
	imgHover   *ebiten.Image
}

func (c *CursorObject) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	img := c.imgNormal
	if c.isHovering {
		img = c.imgHover
	}
	c.img = img
	c.Object.Draw(screen, options)
}

type playStage struct {
	ID         string
	Background *Object
	Object     *Object
	Cursor     *CursorObject
	isFinished bool
}

func (s *playStage) GetID() string {
	return s.ID
}

func (s *playStage) Update() error {
	if s.Finished() {
		return nil
	}

	x, y := ebiten.CursorPosition()
	if x < 0 || y < 0 {
		return nil
	}
	if x >= screenWidth || y >= screenHeight {
		return nil
	}
	s.Cursor.x = float64(x)
	s.Cursor.y = float64(y)

	s.Cursor.isHovering = s.Object.HitBy(s.Cursor.Object)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.Cursor.isHovering && !s.Finished() {
		s.isFinished = true
	}
	return nil
}

func (s *playStage) Finished() bool {
	return s.isFinished
}

func (s *playStage) Draw(screen *ebiten.Image) {
	if s.Finished() {
		return
	}
	scaleGeo := ebiten.GeoM{}
	scaleGeo.Scale(4.2, 4.2)
	opBackground := ebiten.DrawImageOptions{
		GeoM: scaleGeo,
	}
	s.Background.Draw(screen, &opBackground)

	s.Object.Draw(screen, &ebiten.DrawImageOptions{})

	posCursorX := s.Cursor.x - float64(s.Cursor.img.Bounds().Dx())/2
	posCursorY := s.Cursor.y - float64(s.Cursor.img.Bounds().Dy())/2

	s.Cursor.SetPosition(posCursorX, posCursorY)
	s.Cursor.Draw(screen, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("DEBUG MESSAGES: %t, %d, %d", s.Cursor.isHovering, posCursorX, posCursorY))
}

type Game struct {
	CurrentStage Stage
	ApexStage    Stage
	FlorianStage Stage
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
	if g.CurrentStage.Finished() {
		g.CurrentStage = g.FlorianStage
	}
	return g.CurrentStage.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.CurrentStage.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
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

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func NewApexPlayStage(cursorObject *CursorObject) Stage {
	ApexObjectAsset, err := commonAssets.ReadFile("assets/apex/object.png")
	checkErr(err)
	apexStatue := NewObjectFromSprite(ApexObjectAsset, 100, 100)

	ApexBackgroundAsset, err := commonAssets.ReadFile("assets/apex/bg.png")
	checkErr(err)
	apexBg := NewObjectFromSprite(ApexBackgroundAsset, 0, 0)

	return &playStage{ID: "apex", Background: apexBg, Object: apexStatue, Cursor: cursorObject}
}

func NewFlorianPlayStage(cursorObject *CursorObject) Stage {
	FlorianObjectAsset, err := commonAssets.ReadFile("assets/florian/object.png")
	checkErr(err)
	florianObject := NewObjectFromSprite(FlorianObjectAsset, 200, 200)

	FlorianBackgroundAsset, err := commonAssets.ReadFile("assets/florian/bg.png")
	checkErr(err)
	florianBg := NewObjectFromSprite(FlorianBackgroundAsset, 0, 0)

	return &playStage{ID: "florian", Background: florianBg, Object: florianObject, Cursor: cursorObject}
}

func main() {

	cursorAsset, err := commonAssets.ReadFile("assets/cursor.png")
	checkErr(err)
	CursorHoverAsset, err := commonAssets.ReadFile("assets/cursorhover.png")
	checkErr(err)
	cursorObject := NewCursorObject(cursorAsset, CursorHoverAsset)
	apexStage := NewApexPlayStage(cursorObject)
	florianStage := NewFlorianPlayStage(cursorObject)

	g := &Game{ApexStage: apexStage, CurrentStage: apexStage, FlorianStage: florianStage}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Starbound Storyline experience")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
