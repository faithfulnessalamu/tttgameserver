package engine

import "github.com/patrickmn/go-cache"

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
	saveGame(game.id, game)
	return game.id
}

func saveGame(id string, g game) {

}
