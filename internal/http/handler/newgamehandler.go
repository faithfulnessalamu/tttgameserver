package handler

import (
	"log"
	"net/http"

	"github.com/thealamu/tttgameserver/internal/engine"
)

//NewGameHandler handles creating a new game
func NewGameHandler() http.HandlerFunc {
	gE := engine.New()
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := gE.StartNewGame()
		log.Printf("handler.NewGame gameID is %s", gameID)
	}
}