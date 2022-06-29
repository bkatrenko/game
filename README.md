## What is so strange here?..
Here is a small and very simple concept of online [Air hockey](https://en.wikipedia.org/wiki/Air_hockey) game. 
It is use [Ebiten](https://ebiten.org/) - a dead simple 2D game engine written in Go,
And HTTP/UDP protocols to transfer data between clients and servers.

![](https://github.com/bkatrenko/game/blob/main/game.gif)

## Motivation
I like dead simple and beautiful applications that just do one nice thing - for example, allow two people to play online very ugly 90s like Air-hockey.
It is also:
- Show the concept of how online application could work
- Explain how to use Ebiten itself (could be used as an example)
- Show the simple HTTP/UDP server
- Show simple k8s ingress and application deployment
- Contains unit tests (partially in progress)

## How it works:
In the first version of it, we just moved all calculations to the server-side. 
- Every 'N' period of time or when the Ebiten framework sends an update to the callback via UDP we send a small compressed message to the server-side, that contains the full world state, updatable by both clients. Update to the server-side contains the current vector and speed of the player
- Client receives a response, that contains another player/ball vector
- All calculations handled by server-side: collisions/goal detections, simple physics, and servings player's scores
- We use UDP to transfer real-time data, and HTTP to transfer other kinds of data (as a Game Join request)

### How to:
The most simple way to try it out is just a build golang client/server applications and run them. Please, note that it is tested on macOS only!
The location of the desktop client is in the root of the repository, and the location of the server application is in ./srv, both of them could be assembled with 
```
go build
```
command.
The server application should be started firstly:
```
UDP_ADDRESS=localhost:8081 HTTP_ADDRESS=http://localhost:8080 ./srv
```
The server needs only two environment variables to start: host:port where to listen for TCP/UDP traffic. After server application started, clients are ready to join with the following command:
```
UDP_SERVER_HOST_PORT=localhost:8081 HTTP_SERVER_HOST_PORT=http://localhost:8080 PLAYER_ID=johny GAME_ID=1 PLAYER_NUMBER=1 ./game
```
Where: 
| Variable              | Purpose                                                                     |
|-----------------------|-----------------------------------------------------------------------------|
| UDP_SERVER_HOST_PORT  | Should contain host:port of the UDP server                                  |
| HTTP_SERVER_HOST_PORT | Should contain host:port of the HTTP server (will be used to join the game) |
| PLAYER_ID             | Any styring of any length (such a vulnerability!), not necessary unique     |
| GAME_ID               | ID of the game player wants to join, should be unique                        |
| PLAYER_NUMBER         | Number of players you wanna play for (could be 0 or 1)                       |

Additional README about server side is on ./srv folder.
Have fun!! (^_^)