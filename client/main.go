package main

import (
	"github.com/aita/engi"
)

type Game struct {
	*engi.Game
	bot   engi.Drawable
	batch *engi.Batch
	font  *engi.Font
}

func (game *Game) Preload() {
	engi.Files.Add("gopher", "data/gopher_s.png")
	engi.Files.Add("font", "data/font.png")
	game.batch = engi.NewBatch(engi.Width(), engi.Height())
}

func (game *Game) Setup() {
	engi.SetBg(0x2d3739)
	game.bot = engi.Files.Image("gopher")
	game.font = engi.NewGridFont(engi.Files.Image("font"), 20, 20)
}

func (game *Game) Render() {
	game.batch.Begin()
	game.font.Print(game.batch, "GOPHER", 460, 190, 0xffffff)
	game.batch.Draw(game.bot, 512, 320, 0.5, 0.5, 1, 1, 0, 0xffffff, 1)
	game.batch.End()
}

func main() {
	engi.Open("hello", 800, 600, false, &Game{})
}
