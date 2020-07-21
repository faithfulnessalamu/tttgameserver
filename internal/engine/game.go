package engine

import "github.com/dchest/uniuri"

//game holds the state of an ongoing game
type game struct {
	id    string
	state GameState
}

//NewGame returns a new game
func newgame() game {
	return game{}
}

const idLength = 5

func newGameID() string {
	return uniuri.NewLen(idLength)
}
