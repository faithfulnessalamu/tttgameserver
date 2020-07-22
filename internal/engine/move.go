package engine

//Move defines a move by a player
type Move struct {
	Row int `json:"row"`
	Col int `json:"col"`
}
