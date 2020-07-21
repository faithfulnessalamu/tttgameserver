package engine

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

var (
	//ErrGameNotFound is returned when game requested is not in db
	ErrGameNotFound = fmt.Errorf("Game with that ID does not exist")
)

//GameEngine handles game logic
type GameEngine struct {
	db *cache.Cache
}

//New returns a new GameEngine
func New(db *cache.Cache) GameEngine {
	return GameEngine{db: db}
}

//StartNewGame creates a new game and returns a game ID
func (gE GameEngine) StartNewGame() string {
	game := newgame()
	game.id = newGameID() //generate new game id
	//save game
	gE.saveGame(game.id, game)
	return game.id
}

func (gE GameEngine) getGame(id string) (*game, error) {
	gInterface, found := gE.db.Get(id)
	if !found {
		return nil, ErrGameNotFound
	}
	game := gInterface.(*game)
	return game, nil
}

func (gE GameEngine) AttachListener(id string, c chan GameState) {

}

func (gE GameEngine) saveGame(id string, g game) {
	gE.db.Set(id, &g, cache.DefaultExpiration)
}
