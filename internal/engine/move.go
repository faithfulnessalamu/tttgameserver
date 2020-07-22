package engine

//Move defines a move by a player
type Move struct {
	row int `json:"row"`
	col int `json:"col"`
}
