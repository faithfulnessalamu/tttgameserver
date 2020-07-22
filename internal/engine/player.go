package engine

//Player represents a connected player
type Player struct {
	Avatar string `json:"avatar"`
	Score  int    `json:"score"`
	Active bool   `json:"active"`
}

func newPlayer() Player {
	return Player{}
}
