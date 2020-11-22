package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	_ "image/png"
	"log"
	"math/rand"
	"time"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 700
)

type Sprite struct {
	pict *ebiten.Image
	xLoc int
	yLoc int
	dx   int
	dy   int
}

type Game struct {
	playerSprite  Player
	coinSprite    Sprite
	drawOps       ebiten.DrawImageOptions
	collectedGold bool
}

type Player struct{
	name string
	startX int
	startY int
	manifest Sprite
	ammo []Sprite
}

func gotGold(player, gold Sprite) bool {
	goldWidth, goldHeight := gold.pict.Size()
	playerWidth, playerHeight := player.pict.Size()
	if player.xLoc < gold.xLoc+goldWidth &&
		player.xLoc+playerWidth > gold.xLoc &&
		player.yLoc < gold.yLoc+goldHeight &&
		player.yLoc+playerHeight > gold.yLoc {
		return true
	}
	return false
}

func (game *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		game.playerSprite.manifest.dx = -3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		game.playerSprite.manifest.dx = 3
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) || inpututil.IsKeyJustReleased(ebiten.KeyD) {
		game.playerSprite.manifest.dx = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		game.playerSprite.manifest.dy = -3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.playerSprite.manifest.dy = 3
	} else if inpututil.IsKeyJustReleased(ebiten.KeyW) || inpututil.IsKeyJustReleased(ebiten.KeyS) {
		game.playerSprite.manifest.dy = 0
	}
	game.playerSprite.manifest.yLoc += game.playerSprite.manifest.dy
	game.playerSprite.manifest.xLoc += game.playerSprite.manifest.dx
	if game.collectedGold == false {
		game.collectedGold = gotGold(game.playerSprite.manifest, game.coinSprite)
	}
	return nil
}

func (game Game) Draw(screen *ebiten.Image) {
	screen.Fill(colornames.Mediumaquamarine)
	game.drawOps.GeoM.Reset()
	game.drawOps.GeoM.Translate(float64(game.playerSprite.manifest.xLoc), float64(game.playerSprite.manifest.yLoc))
	screen.DrawImage(game.playerSprite.manifest.pict, &game.drawOps)
	if !game.collectedGold {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.coinSprite.xLoc), float64(game.coinSprite.yLoc))
		screen.DrawImage(game.coinSprite.pict, &game.drawOps)
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("The Death Road Through DMF")
	gameObject := Game{}
	loadImage(&gameObject)
	gameObject.playerSprite.manifest.yLoc = ScreenHeight / 2 //Setting player start point
	width, height := gameObject.coinSprite.pict.Size()
	rand.Seed(int64(time.Now().Second()))
	gameObject.coinSprite.xLoc = rand.Intn(ScreenWidth - width)
	gameObject.coinSprite.yLoc = rand.Intn(ScreenHeight - height)
	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}

func loadImage(game *Game) {
	pict, _, err := ebitenutil.NewImageFromFile("assets/galleon.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.playerSprite.manifest.pict = pict
	coins, _, err := ebitenutil.NewImageFromFile("assets/gold-coins.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.pict = coins
}
