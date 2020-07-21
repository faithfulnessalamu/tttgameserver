package engine

import "github.com/dchest/uniuri"

const maxListenersCount = 2

//game holds the state of an ongoing game
type game struct {
	id        string
	state     GameState
	listeners listeners
}

//track the listeners for this game
type listeners struct {
	count    int
	channels [maxListenersCount]chan GameState
}

//NewGame returns a new game
func newgame() game {
	return game{}
}

const idLength = 5

func newGameID() string {
	return uniuri.NewLen(idLength)
}
