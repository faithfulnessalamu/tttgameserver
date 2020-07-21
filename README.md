# tttgameserver
A tic-tac-toe multiplayer game server

![Go](https://github.com/thealamu/tttgameserver/workflows/Go/badge.svg)

## Build
```shell
cd tttgameserver
make
```

## API
There are two endpoints of interest, the **New Game** and **Join Game** websocket endpoints.
### New Game
```
/ws/newgame
```
You can test this endpoint at the terminal using [websocat](https://github.com/vi/websocat)
```
websocat ws://localhost:8080/ws/newgame
```
On connection, the game server returns the game ID. Share this game ID with an opponent.
### Join Game
```
/ws/joingame?gameid={{GAME_ID}}
```
An opponent can join the game using this endpoint or at the terminal using [websocat](https://github.com/vi/websocat)
```
websocat ws://localhost:8080/ws/joingame?gameid={{GAME_ID}}
```