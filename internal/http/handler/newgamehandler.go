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

//NewGameHandler handles creating a new game
type NewGameHandler struct {
	gE       engine.GameEngine
	upgrader websocket.Upgrader
	conn     *websocket.Conn
	logger   *zap.Logger
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
	//init logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Can't initialize zap logger: %v", err)
		return
	}
	nh.logger = logger
	defer nh.logger.Sync()

	gameID := nh.gE.StartNewGame()
	nh.logger.Info("Started a new game",
		zap.String("gameID", gameID),
	)

	//Upgrade the connection
	conn, err := nh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		nh.logger.Error("Upgrading to websocket failed",
			zap.String("gameID", gameID),
			zap.Error(err),
		)
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
			nh.logger.Info("Player disconnected",
				zap.String("gameID", nh.gameID),
				zap.String("player", nh.avatar),
			)
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
				nh.logger.Info("Error reading move",
					zap.String("gameID", nh.gameID),
					zap.String("player", nh.avatar),
					zap.String("err", err.Error()),
				)
				continue
			}
			//log the move
			nh.logger.Info("Received move",
				zap.String("gameID", nh.gameID),
				zap.String("player", nh.avatar),
				zap.Int("row", m.Row),
				zap.Int("col", m.Col),
			)
			//try to do move
			err = nh.gE.MakeMove(nh.gameID, nh.avatar, m)
			if err != nil {
				nh.logger.Info("Error making move",
					zap.String("gameID", nh.gameID),
					zap.String("player", nh.avatar),
					zap.String("err", err.Error()),
				)
				nh.writeString(err.Error())
			}
		}
	}()
	return done //return done while the goroutine above is running
}

func (nh NewGameHandler) writeString(msg string) {
	nh.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}
