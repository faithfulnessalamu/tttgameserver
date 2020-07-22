package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/patrickmn/go-cache"
	"github.com/thealamu/tttgameserver/internal/engine"
	"go.uber.org/zap"
)

//JoinGameHandler handles game joining
type JoinGameHandler struct {
	gE       engine.GameEngine
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	logger   *zap.Logger
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

func (jh JoinGameHandler) doInit() {
	initResponse := struct {
		GameID string `json:"gameID"`
		Avatar string `json:"avatar"`
	}{
		GameID: jh.gameID,
		Avatar: jh.avatar,
	}

	jh.conn.WriteJSON(initResponse)
}

func (jh JoinGameHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//create logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Can't initialize zap logger: %v", err)
		return
	}
	jh.logger = logger
	defer jh.logger.Sync()

	//get gameID from url
	gameID := r.URL.Query().Get("gameid")
	jh.logger.Info("Asking to join a new game",
		zap.String("gameID", gameID),
	)

	conn, err := jh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		jh.logger.Error("Upgrading to websocket failed",
			zap.String("gameID", gameID),
			zap.Error(err),
		)
	}
	jh.conn = conn
	defer jh.conn.Close()

	//register this player in the game engine
	c := make(chan engine.GameState) //game engine should use this channel to send updates
	player, err := jh.gE.NewPlayer(gameID, c)
	if err != nil {
		jh.logger.Info("Error joining game",
			zap.String("gameID", jh.gameID),
			zap.String("err", err.Error()),
		)
		jh.writeString(err.Error()) //write error string to conn
		return                      //we can't go on
	}
	jh.logger.Info("Player joined game",
		zap.String("gameID", gameID),
		zap.String("player", player.Avatar),
	)
	defer jh.gE.RemovePlayer(gameID, player, c) //unregister when player disconnects

	//deliver the gameID and Player Avatar
	jh.gameID = gameID
	jh.avatar = player.Avatar
	jh.doInit()

	done := jh.readMoves() //handle player actions
	for {                  //listen for dispatch or client disconnection
		select {
		case gameState := <-c:
			conn.WriteJSON(gameState) //handle dispatch
		case <-done:
			jh.logger.Info("Player disconnected",
				zap.String("gameID", jh.gameID),
				zap.String("player", jh.avatar),
			)
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
				jh.logger.Info("Error reading move",
					zap.String("gameID", jh.gameID),
					zap.String("player", jh.avatar),
					zap.String("err", err.Error()),
				)
				continue
			}
			//try to do move
			err = jh.gE.MakeMove(jh.gameID, jh.avatar, m)
			if err != nil {
				jh.logger.Info("Error making move",
					zap.String("gameID", jh.gameID),
					zap.String("player", jh.avatar),
					zap.String("err", err.Error()),
				)
				jh.writeString(err.Error())
				continue
			}
			//log the move
			jh.logger.Info("Made move",
				zap.String("gameID", jh.gameID),
				zap.String("player", jh.avatar),
				zap.Int("row", m.Row),
				zap.Int("col", m.Col),
			)
		}
	}()
	return done //return done while the goroutine above is running
}

func (jh JoinGameHandler) writeString(msg string) {
	jh.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
