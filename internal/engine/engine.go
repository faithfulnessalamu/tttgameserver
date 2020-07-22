package engine

import (
	"fmt"

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

//NewPlayer tries to attach a new player to the game
func (gE GameEngine) NewPlayer(gameID string, c chan GameState) (Player, error) {
	//get the game for the id
	game, err := gE.getGame(gameID)
	if err != nil {
		return Player{}, err
	}
	//Check  if the slots are filled
	if game.listeners.count == maxListenersCount {
		return Player{}, ErrNoMorePlayers
	}

	//there is at least one slot available, get a new player
	player := newPlayer()
	player.Active = true
	// add player to game
	//get an avatar from the pool
	player.Avatar = game.nextAvatar()
	game.state.Data.Players = append(game.state.Data.Players, player)

	//attach the listener
	gE.attachListener(game, c)

	return player, nil
}

//RemovePlayer removes a player from the game
func (gE GameEngine) RemovePlayer(gameID string, p Player, c chan GameState) error {
	//get the game for the id
	game, err := gE.getGame(gameID)
	if err != nil {
		return err
	}
	//remove player from game
	for i, gamePlayer := range game.state.Data.Players {
		if p.Avatar == gamePlayer.Avatar {
			game.state.Data.Players = append(game.state.Data.Players[:i], game.state.Data.Players[i+1:]...)
		}
	}
	//unregister player
	gE.unregisterListener(game, c)

	return nil
}

func (gE GameEngine) getGame(id string) (*game, error) {
	gInterface, found := gE.db.Get(id)
	if !found {
		return nil, ErrGameNotFound
	}
	game := gInterface.(*game)
	return game, nil
}

//attachListener adds a listener/channel to a game
func (gE GameEngine) attachListener(g *game, c chan GameState) {
	//attach
	g.listeners.channels = append(g.listeners.channels, c)
	g.listeners.count++

	//dispatch state
	go gE.dispatch(g)
}

//unregisterListener removes a listener/channel from a game
func (gE GameEngine) unregisterListener(g *game, c chan GameState) {
	for i := 0; i < g.listeners.count; i++ {
		if g.listeners.channels[i] == c {
			g.listeners.channels = append(g.listeners.channels[:i], g.listeners.channels[i+1:]...)
			g.listeners.count--
			break
		}
	}
}

func (gE GameEngine) saveGame(id string, g game) {
	gE.db.Set(id, &g, cache.DefaultExpiration)
}

//Dispatch sends the game state to all game listeners
func (gE GameEngine) dispatch(g *game) {
	for _, c := range g.listeners.channels {
		c <- g.state
	}
}
