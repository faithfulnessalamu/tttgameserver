package engine

import (
	"testing"
	"time"

	"github.com/patrickmn/go-cache"
)

func TestCheckGameWon(t *testing.T) {
	testGameID := "ABCDE"
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	testGame := newgame()
	testGame.id = testGameID

	player := newPlayer()
	player.Avatar = "x"
	testGame.state.Data.Players = append(testGame.state.Data.Players, player)

	gE.wonRound(&testGame, "x")
	gE.wonRound(&testGame, "x")
	gE.wonRound(&testGame, "x")

	gE.checkGameWon(&testGame, "x")
	if testGame.state.Win != true {
		t.Error("checkGameWon, expected game win state set to true, got false")
	}
}

func TestHasWonRound(t *testing.T) {
	testGameID := "ABCDE"
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	testGame := newgame()
	testGame.id = testGameID

	player := newPlayer()
	player.Avatar = "x"
	testGame.state.Data.Players = append(testGame.state.Data.Players, player)

	gE.wonRound(&testGame, "x")
	newScore := testGame.state.Data.Players[0].Score
	if newScore != 1 {
		t.Errorf("wonRound does not increase player's score by 1, expected %d, got %d", 1, newScore)
	}
}

func TestRemovePlayer(t *testing.T) {
	testGameID := "ABCDE"
	testChan := make(chan GameState)
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	testGame := newgame()
	testGame.id = testGameID
	gE.saveGame(testGameID, testGame)

	//new player
	//try to remove player from non existent game
	fakeGameID := "FAKE"
	err := gE.RemovePlayer(fakeGameID, Player{}, testChan)
	if err == nil {
		t.Errorf("RemovePlayer expected game not exist error, got nil")
	}

	//add player
	p, _ := gE.NewPlayer(testGameID, testChan)
	//remove player
	err = gE.RemovePlayer(testGameID, p, testChan)
	if err != nil {
		t.Errorf("RemovePlayer expected nil error, got %s", err)
	}

	//validate
	g, _ := gE.getGame(testGameID)
	if len(g.state.Data.Players) != 0 {
		t.Errorf("RemovePlayer does not remove a player, expected length %d, got %d", 0, len(g.state.Data.Players))
	}
}

func TestNewPlayer(t *testing.T) {
	testGameID := "ABCDE"
	testChan := make(chan GameState)
	testDb := cache.New(10*time.Second, 15*time.Second)
	gE := New(testDb)
	testGame := newgame()
	testGame.id = testGameID
	gE.saveGame(testGameID, testGame)

	//new player
	//try to add player to non existent game
	fakeGameID := "FAKE"
	p, err := gE.NewPlayer(fakeGameID, testChan)
	if err == nil {
		t.Errorf("NewPlayer expected game not exist error, got nil")
	}

	//add player
	p, err = gE.NewPlayer(testGameID, testChan)
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if !p.Active {
		t.Errorf("NewPlayer does not set player as active")
	}

	//validate
	g, _ := gE.getGame(testGameID)
	if len(g.state.Data.Players) != 1 {
		t.Errorf("NewPlayer does not insert a new player, expected length %d, got %d", 1, len(g.state.Data.Players))
	}

	//Test adding more than two players
	anotherChan := make(chan GameState)
	yetAnotherChan := make(chan GameState)
	_, err = gE.NewPlayer(testGameID, anotherChan)
	_, err = gE.NewPlayer(testGameID, yetAnotherChan)
	if err == nil {
		t.Errorf("NewPlayer expected error no more players allowed, got nil")
	}
}
