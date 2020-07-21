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
type NewGameHandler struct {
	gE       engine.GameEngine
	upgrader websocket.Upgrader
	conn     *websocket.Conn
}

//GetNewGameHandler creates a new NewGameHandler
func GetNewGameHandler(db *cache.Cache) NewGameHandler {
	return NewGameHandler{
		gE: engine.New(db),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (nh NewGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gameID := nh.gE.StartNewGame()
	log.Printf("handler.NewGame gameID is %s", gameID)

	//Upgrade the connection
	conn, err := nh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("handler.NewGame ", err)
	}
	nh.conn = conn
	defer nh.conn.Close()

	//deliver the gameID
	nh.writeString(gameID)
	//register this player in the game engine
	c := make(chan engine.GameState)
	err = nh.gE.AttachListener(gameID, c) //game engine should use this channel to send updates
	if err != nil {                       //no more players are allowed
		//There should not be any error, this is the game creator
		log.Fatalf("handler.NewGame INVALID STATE: This is the game creator but %s", err)
	}
	defer nh.gE.UnregisterListener(gameID, c) //unregister listener when client disconnects

	done := nh.readMoves() //handle player actions
	for {                  //listen for dispatch or client disconnection
		select {
		case gameState := <-c:
			conn.WriteJSON(gameState) //handle dispatch
		case <-done:
			log.Println("player disconnected")
			return
		}
	}
}

func (nh NewGameHandler) readMoves() chan int {
	done := make(chan int)
	go func() {
		for {
			_, p, err := nh.conn.ReadMessage()
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

func (nh NewGameHandler) writeString(msg string) {
	nh.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
