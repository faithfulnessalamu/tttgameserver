package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/thealamu/tttgameserver/internal/engine"
)

//NewGameHandler handles creating a new game
func NewGameHandler(db *cache.Cache) http.HandlerFunc {
	gE := engine.New(db)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		gameID := gE.StartNewGame()
		log.Printf("handler.NewGame gameID is %s", gameID)

		//Upgrade the connection
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("handler.NewGame ", err)
		}

		//deliver the gameID
		conn.WriteMessage(websocket.TextMessage, []byte(gameID))
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
			}
			fmt.Println(strings.TrimSpace(string(p)))
		}
	}
}
