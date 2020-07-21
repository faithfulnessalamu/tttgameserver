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
	gE.saveGame(game.id, game)
	return game.id
}

func (gE GameEngine) getGame(id string) *game {
	gInterface, found := gE.db.Get(id)
	game := gInterface.(*game)
	return game
}

func (gE GameEngine) AttachListener(id string, c chan engine.GameState) {

}

func (gE GameEngine) saveGame(id string, g game) {
	gE.db.Set(id, &g, cache.DefaultExpiration)
}
