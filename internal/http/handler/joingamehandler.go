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

//JoinGameHandler handles game joining
func JoinGameHandler(db *cache.Cache) http.HandlerFunc {
	gE := engine.New(db)

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	read := func(conn *websocket.Conn) chan int {
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

	writeString := func(conn *websocket.Conn, msg string) {
		conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		//get gameID from url
		gameID := r.URL.Query().Get("gameid")
		log.Printf("Join with gameID %s", gameID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal("handler.JoinGame ", err)
		}

		//register this player in the game engine
		c := make(chan engine.GameState)
		err = gE.AttachListener(gameID, c) //game engine should use this channel to send updates
		if err != nil {                    //no more players are allowed
			//			writeString(err.Error())
			return
		}

	}
}
