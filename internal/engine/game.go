package engine

import "github.com/dchest/uniuri"

const maxListenersCount = 2
const defaultMaxScore = 3

//game holds the state of an ongoing game
type game struct {
	id         string
	state      GameState
	listeners  listeners
	avatarPool []string
}

//track the listeners for this game
type listeners struct {
	count    int
	channels []chan GameState
}

//NewGame returns a new game
func newgame() game {
	g := game{}
	//init game state
	g.state.Data.MaxScore = defaultMaxScore
	g.avatarPool = append(g.avatarPool, "x")
	g.avatarPool = append(g.avatarPool, "o")
	return g
}

const idLength = 5

func newGameID() string {
	return uniuri.NewLen(idLength)
}
