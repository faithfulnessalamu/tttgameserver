# tttgameserver
A tic-tac-toe multiplayer game server

![Go](https://github.com/thealamu/tttgameserver/workflows/Go/badge.svg)

## Build
```shell
cd tttgameserver
make
```
To start the server on a port (e.g. 1025) other than the default (8080):
```shell
tttgameserver --port 1025
```
Run --help to see all server options
```shell
tttgameserver --help
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
On connection, the game server returns the game ID in the game state. Share this game ID with an opponent.
### Join Game
```
/ws/joingame?gameid={{GAME_ID}}
```
An opponent can join the game using this endpoint or at the terminal using [websocat](https://github.com/vi/websocat)
```
websocat ws://localhost:8080/ws/joingame?gameid={{GAME_ID}}
```

## Game State
The game state is the absolute source of truth for all clients.
The game state is returned in JSON first after a successful connection and everytime there is a change; For example, when a player makes a move or when a player disconnects or when a player wins.

Here is a sample of a game state:
```JSON
{
    data: {
        maxScore: 3,
        players: [
            {
                avatar: "o",
                score: 2,
                active: true
            }, 
            {
                avatar: "x",
                score: 3,
                active: true
            }
        ],
    },
    board: [["o","","x"],["x","o",""],["","","o"]],
    win: true,
    winner: "o",
}
```