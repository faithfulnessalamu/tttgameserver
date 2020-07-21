package engine

//GameEngine handles game logic
type GameEngine struct {}

//New returns a new GameEngine
func New() GameEngine {
	return GameEngine{}
}

//StartNewGame creates a new game and returns a game ID
func (gE GameEngine) StartNewGame() string {
	game := newgame()
	game.id = newGameID()	//generate new game id
	//save game

	return game.id
}