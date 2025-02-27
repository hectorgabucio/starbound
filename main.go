package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"starbound-story/internal/object"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// cool sounds
// character_joined.ogg
// cygnus-x1.ogg

//go:embed assets
var commonAssets embed.FS

const (
	screenWidth  = 800
	screenHeight = 600
)

func initOGGPlayer(ctx *audio.Context, path string) *audio.Player {
	audioSrc, err := commonAssets.ReadFile(path)
	checkErr(err)
	s, err := vorbis.DecodeF32(bytes.NewReader(audioSrc))
	checkErr(err)

	p, err := ctx.NewPlayerF32(s)
	checkErr(err)
	return p
}

func main() {

	audioContext := audio.NewContext(44100)
	hornPlayer := initOGGPlayer(audioContext, "assets/horn.ogg")
	questFinishedPlayer := initOGGPlayer(audioContext, "assets/quest_finished.ogg")

	cursorAsset, err := commonAssets.ReadFile("assets/cursor.png")
	checkErr(err)
	CursorHoverAsset, err := commonAssets.ReadFile("assets/cursorhover.png")
	checkErr(err)
	cursorObject := object.NewCursorObject(cursorAsset, CursorHoverAsset)
	apexStage := NewApexPlayStage(cursorObject)
	florianStage := NewFlorianPlayStage(cursorObject)

	cursorEndAsset, err := commonAssets.ReadFile("assets/end/reward_bag.png")
	checkErr(err)
	cursorEndObject := object.NewCursorObject(cursorEndAsset, cursorEndAsset)
	endStage := NewEndStage(cursorEndObject, hornPlayer)

	playStages := []Stage{apexStage, florianStage, endStage}

	g := &Game{playStages: playStages, currentStageIdx: 0, questFinishedPlayer: questFinishedPlayer}

	ebiten.SetVsyncEnabled(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Starbound Storyline experience")
	ebiten.SetCursorMode(ebiten.CursorModeHidden)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	currentStageIdx     int
	playStages          []Stage
	questFinishedPlayer *audio.Player
}

func (g *Game) Update() error {
	currentStage := g.playStages[g.currentStageIdx]
	if currentStage.Finished() {
		g.questFinishedPlayer.Rewind()
		g.questFinishedPlayer.Play()

		g.currentStageIdx++
		return nil
	}
	return currentStage.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.playStages[g.currentStageIdx].Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

type Stage interface {
	GetID() string
	Update() error
	Draw(screen *ebiten.Image)
	Finished() bool
}

type playStage struct {
	ID         string
	Background *object.Object
	Object     *object.Object
	Cursor     *object.CursorObject
	isFinished bool
}

func NewApexPlayStage(cursorObject *object.CursorObject) Stage {
	ApexObjectAsset, err := commonAssets.ReadFile("assets/apex/object.png")
	checkErr(err)
	apexStatue := object.NewObjectFromSprite(ApexObjectAsset, 100, 100)

	ApexBackgroundAsset, err := commonAssets.ReadFile("assets/apex/bg.png")
	checkErr(err)
	apexBg := object.NewObjectFromSprite(ApexBackgroundAsset, 0, 0)

	return &playStage{ID: "apex", Background: apexBg, Object: apexStatue, Cursor: cursorObject}
}

func NewFlorianPlayStage(cursorObject *object.CursorObject) Stage {
	FlorianObjectAsset, err := commonAssets.ReadFile("assets/florian/object.png")
	checkErr(err)
	florianObject := object.NewObjectFromSprite(FlorianObjectAsset, 200, 200)

	FlorianBackgroundAsset, err := commonAssets.ReadFile("assets/florian/bg.png")
	checkErr(err)
	florianBg := object.NewObjectFromSprite(FlorianBackgroundAsset, 0, 0)

	return &playStage{ID: "florian", Background: florianBg, Object: florianObject, Cursor: cursorObject}
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

	posCursorX := float64(x) - float64(s.Cursor.Img.Bounds().Dx())/2
	posCursorY := float64(y) - float64(s.Cursor.Img.Bounds().Dy())/2

	s.Cursor.SetPosition(posCursorX, posCursorY)

	s.Cursor.IsHovering = s.Object.HitBy(s.Cursor.Object)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.Cursor.IsHovering && !s.Finished() {
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

	s.Cursor.Draw(screen, &ebiten.DrawImageOptions{})

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("DEBUG MESSAGES: %t, %d, %d", s.Cursor.IsHovering, s.Cursor.X, s.Cursor.Y))
}

type endStage struct {
	coolDownHorn        int
	currentCoolDownHorn int
	ID                  string
	Background          *object.Object
	Cursor              *object.CursorObject
	hornPlayer          *audio.Player
}

func NewEndStage(cursorObject *object.CursorObject, hornPlayer *audio.Player) Stage {
	EndBackgroundAsset, err := commonAssets.ReadFile("assets/end/bg.png")
	checkErr(err)
	endBg := object.NewObjectFromSprite(EndBackgroundAsset, 0, 0)
	return &endStage{ID: "end", Background: endBg, Cursor: cursorObject, hornPlayer: hornPlayer, coolDownHorn: 50, currentCoolDownHorn: 100}
}

func (s *endStage) GetID() string {
	return s.ID
}

func (s *endStage) Update() error {
	s.currentCoolDownHorn--

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && s.currentCoolDownHorn <= 0 {
		s.hornPlayer.Rewind()
		s.hornPlayer.Play()
		s.currentCoolDownHorn = s.coolDownHorn
	}

	x, y := ebiten.CursorPosition()

	posCursorX := float64(x) - float64(s.Cursor.Img.Bounds().Dx())/2
	posCursorY := float64(y) - float64(s.Cursor.Img.Bounds().Dy())/2

	s.Cursor.SetPosition(posCursorX, posCursorY)

	return nil
}

func (s *endStage) Finished() bool {
	return false
}

func (s *endStage) Draw(screen *ebiten.Image) {
	opBackground := ebiten.DrawImageOptions{}
	s.Background.Draw(screen, &opBackground)

	ebitenutil.DebugPrintAt(screen, "Made with <3 by Héctor Gabucio", 310, 500)

	s.Cursor.Draw(screen, &ebiten.DrawImageOptions{})
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
