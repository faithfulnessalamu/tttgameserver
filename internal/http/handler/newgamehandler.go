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
		//register this player in the game engine
		c := make(chan engine.GameState)
		done := read(conn)
		for {
			select {
			case gameState := <-c:
				conn.WriteJSON(gameState)
			case <-done:
				log.Println("Client disconnected")
				return
			}
		}
	}

}

func read(conn *websocket.Conn) chan int {
	done := make(chan int)
	go func() {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				//connection closed
				done <- 1
				break
			}
			fmt.Println(strings.TrimSpace(string(p)))
		}
	}()
	return done //return done while the goroutine above is running
}
