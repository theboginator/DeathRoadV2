package main

import (
	"fmt"
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
	ScreenWidth   = 800
	ScreenHeight  = 700
	OrdnanceSpeed = 8
)

type Sprite struct {
	pict *ebiten.Image
	xLoc int
	yLoc int
	dx   int
	dy   int
}

type Game struct {
	playerSprite   Player
	playerOrdnance Ordnance
	coinSprite     Sprite
	drawOps        ebiten.DrawImageOptions
	collectedGold  bool
}

type Player struct {
	name     string
	health   int
	startX   int
	startY   int
	firing   bool
	manifest Sprite
	weapon   Sprite
}

type Ordnance struct {
	manifest Sprite
	consumed bool
}

func getPlayerInput(game *Game) { //Handle any movement from the player, and initiate any ordnance the player fires
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		game.playerSprite.firing = true
		game.playerOrdnance.manifest.dx, game.playerOrdnance.manifest.dy = ebiten.CursorPosition() //Get the direction in which to fire ordnance
		game.playerOrdnance.manifest.xLoc = game.playerSprite.manifest.xLoc                        //Set the start point for new ordnance
		game.playerOrdnance.manifest.yLoc = game.playerSprite.manifest.yLoc
		game.playerOrdnance.consumed = false
		//screen.DrawImage(game.playerSprite.manifest.pict, &game.drawOps)
		fmt.Println("Fired ordnance. Coords: ", game.playerOrdnance.manifest.dx, game.playerOrdnance.manifest.dy)
	} else {
		game.playerSprite.firing = false
	}
}

func trackPlayer(game *Game) { //Move the player per keyboard input
	game.playerSprite.manifest.yLoc += game.playerSprite.manifest.dy
	game.playerSprite.manifest.xLoc += game.playerSprite.manifest.dx
}

func trackOrdnance(game *Game) { //Move any ordnance in the direction it was fired in
	game.playerOrdnance.manifest.yLoc += game.playerOrdnance.manifest.dy
	game.playerOrdnance.manifest.xLoc += game.playerOrdnance.manifest.dx
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
	getPlayerInput(game)
	trackPlayer(game)
	//trackOrdnance(game)
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
	if game.playerSprite.firing {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.playerOrdnance.manifest.xLoc), float64(game.playerOrdnance.manifest.yLoc))
		screen.DrawImage(game.playerOrdnance.manifest.pict, &game.drawOps)
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func loadPlayer(game *Game) {
	game.playerSprite.manifest.yLoc = ScreenHeight / 2 //Setting player start point
	pict, _, err := ebitenutil.NewImageFromFile("assets/galleon.png")
	if err != nil {
		log.Fatal("failed to load player image", err)
	}
	game.playerSprite.manifest.pict = pict
	pict, _, err = ebitenutil.NewImageFromFile("assets/player-ammo.png")
	if err != nil {
		log.Fatal("failed to load ammunition image", err)
	}
	game.playerOrdnance.manifest.pict = pict
}

func loadCoinSprite(game *Game) {
	coins, _, err := ebitenutil.NewImageFromFile("assets/gold-coins.png")
	if err != nil {
		log.Fatal("failed to load image", err)
	}
	game.coinSprite.pict = coins
	width, height := game.coinSprite.pict.Size()
	rand.Seed(int64(time.Now().Second()))
	game.coinSprite.xLoc = rand.Intn(ScreenWidth - width)
	game.coinSprite.yLoc = rand.Intn(ScreenHeight - height)
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("The Death Road Through DMF")
	gameObject := Game{}
	loadPlayer(&gameObject)
	loadCoinSprite(&gameObject)
	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}
