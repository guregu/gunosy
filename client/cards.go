package main

import (
	"fmt"
	"github.com/aita/engi"
	"github.com/gophergala/gunosy/daihinmin"
)

type CardSprite struct {
	*engi.Sprite
	daihinmin.Card
}

func NewCardSprite(card daihinmin.Card, x, y float32) *CardSprite {
	texture := engi.Files.Image(cardTextureName(card))
	region := engi.NewRegion(texture, 0, 0, int(texture.Width()), int(texture.Height()))
	sprite := &CardSprite{
		engi.NewSprite(region, x, y),
		card,
	}
	return sprite
}

func (card *CardSprite) HitTest(x, y float32) bool {
	// NOTE: Positionしか考慮しない
	if card.Position.X < x && x < card.Position.X {
		if card.Position.Y < y && y < card.Position.Y {
			return true
		}
	}
	return false
}

func cardTextureName(card daihinmin.Card) string {
	return fmt.Sprintf("%c%d", card.Suit, card.Rank)
}
