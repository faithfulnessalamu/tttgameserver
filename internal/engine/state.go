package engine

//GameState holds the state of an ongoing game
type GameState struct {
	Data   Data         `json:"data"`
	Board  [3][3]string `json:"board"`
	Win    bool         `json:"win"`
	Winner string       `json:"winner"`
}

//Data holds current game data
type Data struct {
	MaxScore int    `json:"maxScore`
	Playerx  Player `json:"playerx"`
	Playero  Player `json:"playero"`
}

//Player represents a connected player
type Player struct {
	ID     int  `json:"id"`
	Score  int  `json:"score"`
	Active bool `json:"active"`
}
