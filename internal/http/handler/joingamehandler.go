package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/thealamu/tttgameserver/internal/engine"
)

//JoinGameHandler handles game joining
type JoinGameHandler struct {
	gE       engine.GameEngine
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	gameID   string
	avatar   string
}

//GetJoinGameHandler creates a new JoinGameHandler
func GetJoinGameHandler(db *cache.Cache) JoinGameHandler {
	return JoinGameHandler{
		gE: engine.New(db),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func (jh JoinGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//get gameID from url
	gameID := r.URL.Query().Get("gameid")
	log.Printf("handler.JoinGame gameID is %s", gameID)

	conn, err := jh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("handler.JoinGame ", err)
	}
	jh.conn = conn
	defer jh.conn.Close()

	//register this player in the game engine
	c := make(chan engine.GameState) //game engine should use this channel to send updates
	player, err := jh.gE.NewPlayer(gameID, c)
	if err != nil {
		log.Println(err)
		jh.writeString(err.Error()) //write error string to conn
		return                      //we can't go on
	}
	defer jh.gE.RemovePlayer(gameID, player, c) //unregister when player disconnects

	//TODO: deliver the gameID and Player Avatar
	jh.gameID = gameID
	jh.avatar = player.Avatar
	//nh.writeString(gameID)

	done := jh.readMoves() //handle player actions
	for {                  //listen for dispatch or client disconnection
		select {
		case gameState := <-c:
			conn.WriteJSON(gameState) //handle dispatch
		case <-done:
			log.Println("handler.JoinGame player disconnected")
			return
		}
	}
}

func (jh JoinGameHandler) readMoves() chan int {
	done := make(chan int)
	go func() {
		for {
			_, p, err := jh.conn.ReadMessage()
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
			err = jh.gE.MakeMove(jh.gameID, jh.avatar, m)
			if err != nil {
				log.Printf("handler.readMoves %s", err)
				jh.writeString(err.Error())
			}
		}
	}()
	return done //return done while the goroutine above is running
}

func (jh JoinGameHandler) writeString(msg string) {
	jh.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
