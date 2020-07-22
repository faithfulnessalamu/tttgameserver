package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/thealamu/tttgameserver/internal/engine"
)

//NewGameHandler handles creating a new game
type NewGameHandler struct {
	gE       engine.GameEngine
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	gameID   string
	avatar   string
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
		log.Println("handler.NewGame ", err)
		return
	}
	nh.conn = conn
	defer nh.conn.Close()

	//register this player in the game engine
	c := make(chan engine.GameState) //game engine should use this channel to send updates
	//new player
	player, err := nh.gE.NewPlayer(gameID, c)
	if err != nil {
		//There should not be any error, this is the game creator
		log.Printf("handler.NewGame INVALID STATE: This is the game creator but %s", err)
		return
	}
	defer nh.gE.RemovePlayer(gameID, player, c) //unregister when player disconnects

	//TODO: deliver the gameID and Player Avatar
	nh.gameID = gameID
	nh.avatar = player.Avatar
	nh.writeString(gameID)

	done := nh.readMoves() //handle player actions
	for {                  //listen for dispatch or client disconnection
		select {
		case gameState := <-c:
			conn.WriteJSON(gameState) //handle dispatch
		case <-done:
			log.Println("handler.NewGame player disconnected")
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
			//handle move
			var m engine.Move
			err = json.Unmarshal(p, &m)
			if err != nil {
				log.Printf("handler.readMoves %s", err)
				continue
			}
			//try to do move
			err = nh.gE.MakeMove(nh.gameID, nh.avatar, m)
			if err != nil {
				log.Printf("handler.readMoves %s", err)
				nh.writeString(err.Error())
			}
		}
	}()
	return done //return done while the goroutine above is running
}

func (nh NewGameHandler) writeString(msg string) {
	nh.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
