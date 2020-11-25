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
	ScreenWidth  = 1024
	ScreenHeight = 768
	NumEnemies   = 3
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
	enemySprites   [NumEnemies]Enemy
	coinSprite     Sprite
	drawOps        ebiten.DrawImageOptions
	activeOrdnance bool
	collectedGold  bool
}

type Player struct {
	name     string
	health   int
	startX   int
	startY   int
	score    int32
	manifest Sprite
}

type Enemy struct {
	name     string
	lastMove time.Time
	health   int
	defeated bool
	startX   int
	startY   int
	manifest Sprite
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
	} else if inpututil.IsKeyJustReleased(ebiten.KeyA) {
		game.playerSprite.manifest.dx = 0
	} else if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		game.playerSprite.manifest.dx = 0
	} else if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		game.playerSprite.manifest.dy = -3
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		game.playerSprite.manifest.dy = 3
	} else if inpututil.IsKeyJustReleased(ebiten.KeyW) {
		game.playerSprite.manifest.dy = 0
	} else if inpututil.IsKeyJustReleased(ebiten.KeyS) {
		game.playerSprite.manifest.dy = 0
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		game.activeOrdnance = true
		launchPlayerOrdnance(game)
	}
}

func trackPlayer(game *Game) { //Move the player per keyboard input
	game.playerSprite.manifest.yLoc += game.playerSprite.manifest.dy
	game.playerSprite.manifest.xLoc += game.playerSprite.manifest.dx
	for i := range game.enemySprites { //Check for a potential collision with an enemy
		if madeContact(game.playerSprite.manifest, game.enemySprites[i].manifest) {
			fmt.Println("You got n00ned by an enemy :(") //We die
		}
	}

}

func trackOrdnance(game *Game) { //Move any ordnance in the direction it was fired in
	game.playerOrdnance.manifest.yLoc += game.playerOrdnance.manifest.dy
	game.playerOrdnance.manifest.xLoc += game.playerOrdnance.manifest.dx
	if game.playerOrdnance.manifest.xLoc > 800 || game.playerOrdnance.manifest.xLoc < 0 || game.playerOrdnance.manifest.yLoc > 700 || game.playerOrdnance.manifest.yLoc < 0 { //If we've hit a border
		game.activeOrdnance = false
	} else if game.collectedGold == false {
		game.collectedGold = gotGold(game.playerOrdnance.manifest, game.coinSprite) //If ordnance is active, check if it collided with enemy
	}
	for i := range game.enemySprites {
		if !game.enemySprites[i].defeated {
			if madeContact(game.playerOrdnance.manifest, game.enemySprites[i].manifest) { //If the ordnance touched an enemy
				fmt.Println("You n00ned an enemy!")
				game.enemySprites[i].defeated = true
			}
		}
	}
}

func launchPlayerOrdnance(game *Game) { //Initiate the launch of player ordnance
	pict, _, err := ebitenutil.NewImageFromFile("assets/cannonball.png")
	if err != nil {
		log.Fatal("failed to load ammunition image", err)
	}
	game.playerOrdnance.manifest.pict = pict
	if game.playerSprite.manifest.dx == 0 && game.playerSprite.manifest.dy == 0 { //Launch ordnance to the right along the x-axis if the player is stationary.
		game.playerOrdnance.manifest.dx = 6
		game.playerOrdnance.manifest.dy = 0
	} else {
		game.playerOrdnance.manifest.dx = game.playerSprite.manifest.dx * 2
		game.playerOrdnance.manifest.dy = game.playerSprite.manifest.dy * 2 //Set the direction to fire new ordnance
	}
	game.playerOrdnance.manifest.xLoc = game.playerSprite.manifest.xLoc //Set the start point for new ordnance to the player's current position
	game.playerOrdnance.manifest.yLoc = game.playerSprite.manifest.yLoc
	game.playerOrdnance.consumed = false
	fmt.Println("Fired ordnance. Coords: ", game.playerOrdnance.manifest.dx, game.playerOrdnance.manifest.dy)
}

func gotGold(ordnance, gold Sprite) bool {
	goldWidth, goldHeight := gold.pict.Size()
	ordWidth, ordHeight := ordnance.pict.Size()
	if ordnance.xLoc < gold.xLoc+goldWidth &&
		ordnance.xLoc+ordWidth > gold.xLoc &&
		ordnance.yLoc < gold.yLoc+goldHeight &&
		ordnance.yLoc+ordHeight > gold.yLoc {
		return true
	}
	return false
}

func madeContact(manifestA Sprite, manifestB Sprite) bool { //Check if 2 sprite objects are in a collision condition
	aWidth, aHeight := manifestA.pict.Size()
	bWidth, bHeight := manifestB.pict.Size()
	if manifestA.xLoc < manifestB.xLoc+aWidth &&
		manifestA.xLoc+bWidth > manifestB.xLoc &&
		manifestA.yLoc < manifestB.yLoc+aHeight &&
		manifestA.yLoc+bHeight > manifestB.yLoc {
		return true
	}
	return false
}

func (game *Game) Update() error {
	getPlayerInput(game)
	trackPlayer(game)
	if game.activeOrdnance {
		trackOrdnance(game)
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
	if game.activeOrdnance {
		game.drawOps.GeoM.Reset()
		game.drawOps.GeoM.Translate(float64(game.playerOrdnance.manifest.xLoc), float64(game.playerOrdnance.manifest.yLoc))
		screen.DrawImage(game.playerOrdnance.manifest.pict, &game.drawOps)
	}
	for i := range game.enemySprites { //For each enemy in the enemy sprite array
		if !game.enemySprites[i].defeated { //Draw the undefeated ones
			game.drawOps.GeoM.Reset()
			x := float64(game.enemySprites[i].manifest.xLoc)
			y := float64(game.enemySprites[i].manifest.yLoc)
			game.drawOps.GeoM.Translate(x, y)
			screen.DrawImage(game.enemySprites[i].manifest.pict, &game.drawOps)
		}
	}
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func loadPlayer(game *Game) {
	game.playerSprite.manifest.yLoc = ScreenHeight / 2 //Setting player start point
	pict, _, err := ebitenutil.NewImageFromFile("assets/player.png")
	if err != nil { //firing   bool
		log.Fatal("failed to load player image", err)
	}
	game.playerSprite.manifest.pict = pict
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

func loadEnemies(game *Game) {
	for i := range game.enemySprites {
		pict, _, err := ebitenutil.NewImageFromFile("assets/galleon.png")
		if err != nil {
			log.Fatal("Failed to load enemy image", err)
		}
		game.enemySprites[i].manifest.pict = pict
		game.enemySprites[i].health = 2
		game.enemySprites[i].defeated = false
		width, height := game.enemySprites[i].manifest.pict.Size()
		game.enemySprites[i].manifest.xLoc = rand.Intn(ScreenWidth - width)
		game.enemySprites[i].manifest.yLoc = rand.Intn(ScreenHeight - height)
	}
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("The Death Road Through DMF")
	gameObject := Game{}
	loadPlayer(&gameObject)
	loadEnemies(&gameObject)
	loadCoinSprite(&gameObject)
	if err := ebiten.RunGame(&gameObject); err != nil {
		log.Fatal("Oh no! something terrible happened", err)
	}
}
