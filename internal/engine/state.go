package engine

//GameState holds the state of an ongoing game
type GameState struct {
	Data   Data         `json:"data"`
	Board  [3][3]string `json:"board"`
	Win    bool         `json:"win"`
	Winner string       `json:"winner"`
	Turn   string       `json:"turn"` //whose turn is it?
}

//Data holds current game data
type Data struct {
	MaxScore int      `json:"maxScore"`
	Players  []Player `json:"players"`
}
