package handler

import (
	"log"
	"net/http"

	"github.com/patrickmn/go-cache"
	"github.com/thealamu/tttgameserver/internal/engine"
)

//NewGameHandler handles creating a new game
func NewGameHandler(db *cache.Cache) http.HandlerFunc {
	gE := engine.New(db)
	return func(w http.ResponseWriter, r *http.Request) {
		gameID := gE.StartNewGame()
		log.Printf("handler.NewGame gameID is %s", gameID)
	}
}
