package engine

import (
	"fmt"
	"log"

	"github.com/patrickmn/go-cache"
)

var (
	//ErrGameNotFound is returned when game requested is not in db
	ErrGameNotFound = fmt.Errorf("Game with that ID does not exist")
	//ErrNoMorePlayers is returned when trying to attach after reaching maxListenerCount
	ErrNoMorePlayers = fmt.Errorf("No more players allowed")
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

//AttachListener adds a listener/channel to a game
func (gE GameEngine) AttachListener(id string, c chan GameState) error {
	//get the game for the id
	game, err := gE.getGame(id)
	if err != nil {
		return err
	}

	//get current game listeners count
	count := game.listeners.count
	//check if the slots are filled
	if count == maxListenersCount {
		return ErrNoMorePlayers
	}
	//attach
	game.listeners.channels[count] = c
	game.listeners.count++

	//dispatch state
	go gE.dispatch(id)

	return nil
}

func (gE GameEngine) saveGame(id string, g game) {
	gE.db.Set(id, &g, cache.DefaultExpiration)
}

//Dispatch sends the game state to all game listeners
func (gE GameEngine) dispatch(id string) {
	game, err := gE.getGame(id)
	if err != nil {
		log.Println("gameEngine.dispatch could not dispatch game state", err)
	}
	for i, c := range game.listeners.channels {
		c <- game.state
	}
}
