package main

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 480
	screenHeight = 640
	gravity      = 0.4
	jumpSpeed    = -18
	playerScale  = 0.25
)

type Game struct {
	playerImage      *ebiten.Image
	backgroundImage  *ebiten.Image
	platformImage    *ebiten.Image
	playerX          float64
	playerY          float64
	playerSpeedY     float64
	facingRight      bool
	platforms        []Platform
	highestPlatformY float64
}

type Platform struct {
	X float64
	Y float64
}

func assetPath(filename string) string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to get executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	log.Print(exeDir)

	return filepath.Join(exeDir, "../Resources", filename)
}

func getImage(filepath string) *ebiten.Image {
	image, _, err := ebitenutil.NewImageFromFile(assetPath(filepath))
	if err != nil {
		log.Fatalf("Failed to load background image: %v", err)
	}

	return image
}

func NewGame() *Game {
	backgroundImage := getImage("background.png")
	playerImage := getImage("player2.png")
	platformImage := getImage("platform.png")

	platforms := []Platform{
		{X: screenWidth / 2, Y: 550},
		{X: 0, Y: 500},
		{X: 200, Y: 400},
		{X: 300, Y: 300},
	}

	return &Game{
		playerImage:      playerImage,
		backgroundImage:  backgroundImage,
		platformImage:    platformImage,
		playerX:          screenWidth / 2,
		playerY:          screenHeight - 60,
		playerSpeedY:     0,
		facingRight:      true,
		platforms:        platforms,
		highestPlatformY: 300,
	}
}

func (g *Game) Update() error {
	g.playerY += g.playerSpeedY
	g.playerSpeedY += gravity

	if g.playerY >= screenHeight-60 {
		g.playerY = screenHeight - 60
		g.playerSpeedY = jumpSpeed
	}

	for _, platform := range g.platforms {
		playerWidth := float64(g.playerImage.Bounds().Dx()) * playerScale
		playerHeight := float64(g.playerImage.Bounds().Dy()) * playerScale
		platformWidth := float64(g.platformImage.Bounds().Dx())

		if g.playerSpeedY > 0 &&
			g.playerY+playerHeight > platform.Y &&
			g.playerY+playerHeight-g.playerSpeedY < platform.Y &&
			g.playerX+playerWidth > platform.X &&
			g.playerX < platform.X+platformWidth {
			g.playerY = platform.Y - playerHeight
			g.playerSpeedY = jumpSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		g.playerX -= 5
		g.facingRight = false

		if g.playerX < -40 {
			g.playerX = 480
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		g.playerX += 5
		g.facingRight = true

		if g.playerX > 520 {
			g.playerX = 0
		}
	}

	if g.playerY < screenHeight/2 {
		offset := screenHeight/2 - g.playerY
		g.playerY = screenHeight / 2

		for i := range g.platforms {
			g.platforms[i].Y += offset
		}

		g.highestPlatformY += offset

		for g.highestPlatformY > 0 {
			newPlatform := Platform{
				X: float64(rand.Intn(screenWidth - int(g.platformImage.Bounds().Dx()))),
				Y: g.highestPlatformY - 100,
			}
			g.platforms = append(g.platforms, newPlatform)
			g.highestPlatformY -= 100
		}

		newPlatforms := []Platform{}
		for _, platform := range g.platforms {
			if platform.Y < screenHeight {
				newPlatforms = append(newPlatforms, platform)
			}
		}
		g.platforms = newPlatforms
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.backgroundImage, nil)

	for _, platform := range g.platforms {
		options := &ebiten.DrawImageOptions{}
		options.GeoM.Translate(platform.X, platform.Y)
		screen.DrawImage(g.platformImage, options)
	}

	options := &ebiten.DrawImageOptions{}

	if !g.facingRight {
		options.GeoM.Scale(-1, 1)
		options.GeoM.Translate(float64(g.playerImage.Bounds().Dx()), 0)
	}

	options.GeoM.Scale(playerScale, playerScale)
	options.GeoM.Translate(g.playerX, g.playerY)

	screen.DrawImage(g.playerImage, options)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	rand.Seed(rand.Int63())
	game := NewGame()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Doodle Jump")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
