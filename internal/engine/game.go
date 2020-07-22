package engine

import (
	"github.com/dchest/uniuri"
)

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
	g.avatarPool = append(g.avatarPool, []string{"x", "o"}...)
	g.state.Turn = "o" //o plays first
	return g
}

func (g *game) returnAvatar(avt string) {
	g.avatarPool = append(g.avatarPool, avt)
}

func (g *game) nextAvatar() string {
	l := len(g.avatarPool)
	j := l - 1
	avt := g.avatarPool[j]
	//remove this avatar from the pool
	g.avatarPool = append(g.avatarPool[:j], g.avatarPool[j+1:]...)
	return avt
}

const idLength = 5

func newGameID() string {
	return uniuri.NewLen(idLength)
}
