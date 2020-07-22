package engine

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestNewGameEngine(t *testing.T) {
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	if gE.db != testDb {
		t.Errorf("GameEngine does not use injected db")
	}
}

func TestStartNewGame(t *testing.T) {
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	gameID := gE.StartNewGame()
	if gameID == "" || len(gameID) != idLength {
		t.Errorf("Invalid gameID: got %s, want a %d length string", gameID, idLength)
	}
}

func TestSaveGame(t *testing.T) {
	testDb := cache.New(1*time.Minute, 2*time.Minute)
	gE := New(testDb)
	testGameID := "ABCDE"
	testGame := newgame()
	testGame.id = testGameID

	gE.saveGame(testGame.id, testGame)

	g, found := testDb.Get(testGame.id)
	if !found {
		t.Errorf("GameEngine's saveGame does not save to injected db")
	}
	retGame, ok := g.(*game)
	if !ok {
		t.Errorf("GameEngine's saveGame does not save a game type")
	} else if retGame.id != testGameID {
		t.Errorf("GameEngine's saveGame saved with wrong ID: expected id: %s, got %s.", testGameID, retGame.id)
	}
}

func TestGetGame(t *testing.T) {
	testDb := cache.New(1*time.Minute, 2*time.Minute)
	gE := New(testDb)
	testGameID := "ABCDE"
	testGame := newgame()
	testGame.id = testGameID
	//save game
	gE.saveGame(testGame.id, testGame)
	//try to get the game back
	retGame, err := gE.getGame(testGameID)
	if err != nil {
		t.Errorf("GameEngine's getGame returns error: %s", err)
		return //the remainder of the tests depend on a proper return value
	}
	if retGame.id != testGameID {
		t.Errorf("GameEngine's getGame returns game with wrong ID: expected id: %s, got %s", testGameID, retGame.id)
	}
}

func TestGetGameNotFound(t *testing.T) {
	testDb := cache.New(1*time.Minute, 2*time.Minute)
	gE := New(testDb)
	fakeGameID := "HELLO FAKER"
	testGameID := "ABCDE"
	testGame := newgame()
	testGame.id = testGameID
	//save game
	gE.saveGame(testGame.id, testGame)
	//try to get the game back
	_, err := gE.getGame(fakeGameID)
	if err == nil {
		t.Error("GameEngine's getGame should return a not found error, got nil")
	}
}

func TestAttachListener(t *testing.T) {
	testDb := cache.New(1*time.Minute, 2*time.Minute)
	gE := New(testDb)
	testGameID := "ABCDE"
	testGame := newgame()
	testGame.id = testGameID
	//save game
	gE.saveGame(testGame.id, testGame)

	testChannel := make(chan GameState)
	gE.attachListener(&testGame, testChannel)

	if testGame.listeners.count != 1 {
		t.Errorf("AttachListener does not attach, expected %d listeners, got %d", 1, testGame.listeners.count)
	}
}

func TestUnregisterListener(t *testing.T) {
	testDb := cache.New(1*time.Minute, 2*time.Minute)
	gE := New(testDb)
	testGameID := "ABCDE"
	testGame := newgame()
	testGame.id = testGameID
	//save game
	gE.saveGame(testGame.id, testGame)

	testChannel1 := make(chan GameState)
	testChannel2 := make(chan GameState)

	//attach listeners
	//note: do not attach more than maxListenersCount
	gE.attachListener(&testGame, testChannel1)
	gE.attachListener(&testGame, testChannel2)

	//remove second listener
	gE.unregisterListener(&testGame, testChannel2)

	if testGame.listeners.count != 1 {
		t.Errorf("UnregisterListener fails, expected %d listeners after unregister, got %d", 1, testGame.listeners.count)
	}
}

func TestUpdateTurn(t *testing.T) {
	testGame := newgame() //first turn is o
	gE := New(cache.New(10*time.Second, 15*time.Second))

	gE.updateTurn(&testGame)
	if testGame.state.Turn != "x" {
		t.Errorf("updateTurn does not update turn in game state")
	}

	gE.updateTurn(&testGame)
	if testGame.state.Turn != "o" {
		t.Errorf("updateTurn does not update turn in game state")
	}
}

func TestNewRound(t *testing.T) {
	testGameID := "ABCDE"
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	testGame := newgame()
	testGame.id = testGameID

	player := newPlayer()
	player.Avatar = "x"
	testGame.state.Data.Players = append(testGame.state.Data.Players, player)

	//play something
	testGame.state.Board[0][0] = "x"

	gE.newRound(&testGame)

	if testGame.state.Board != [3][3]string{} {
		t.Errorf("newRound, expected board to be reset, it is not")
	}

	if testGame.state.Data.Players[0].Score != 0 {
		t.Errorf("newRound, expected player scores to be reset")
	}
}
