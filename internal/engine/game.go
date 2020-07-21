package engine

//game holds the state of an ongoing game
type game struct {
	state gameState
}

//NewGame returns a new game
func NewGame() game {
	return game{}
}
