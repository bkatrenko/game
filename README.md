## What is so strange here?..
Here is small and very simple concept of online [Air hockey](https://en.wikipedia.org/wiki/Air_hockey) game. 
It is using [Ebiten](https://ebiten.org/) - dead simple 2D game engine written in Go to draw,
And HTTP/UDP to transfer data between clients/server.

![](https://github.com/bkatrenko/game/blob/main/game.gif)

## Motivation
I like dead simple and beautiful applications that just do one nice think - allow two people to play online in very ugly 90s like Air-hockey.
It is also:
- Showing the concept of how online application could work
- Explain how to use Ebiten itself (could be used as an example)
- Show the simple HTTP/UDP server
- Show simple k8s ingress and application deployment
- Contains unit tests

## How it works:
In the first version of it, we just moved all calculations into the server side. 
- Every 'N' period of time or when Ebiten framework send an update to the callback we send a small update to the server side, that contains full world state, updatable by both clients. Update the server side contains current vector of the player
- Client receive a response, that contains another player/ball vector
- All calculations handled by server side: collision/goal detections, simple physics and servings player's scores
- We use UDP to transfer realtime data, and HTTP to transfer another kinds of data (as Game Join request)

### How to:
The most simple way to try it our, is just a build golang client/server and run. Please, note that it is tested on MacOS only!
Location of the desktop client is in the root of the repository, and the location of the server application is in ./srv, both of them build be builded with 
```
go build
```
command, and the output will be just normal binary files.
Server application should be started firstly:
```
UDP_ADDRESS=localhost:8081 HTTP_ADDRESS=http://localhost:8080 ./srv
```
Server needs only two environment variables to start: host:port where to listen for TCP/UDP traffic. After server started, clients are ready to join with the following command:
```
UDP_SERVER_HOST_PORT=localhost:8081 HTTP_SERVER_HOST_PORT=http://localhost:8080 PLAYER_ID=johny GAME_ID=1 PLAYER_NUMBER=1 ./game
```
Where: 
| Variable              | Purpose                                                                     |
|-----------------------|-----------------------------------------------------------------------------|
| UDP_SERVER_HOST_PORT  | Should contain host:port of the UDP server                                  |
| HTTP_SERVER_HOST_PORT | Should contain host:port of the HTTP server (will be used to join the game) |
| PLAYER_ID             | Any styring of any length (such a volnurability!), not necessary unique     |
| GAME_ID               | ID of the game player want to join, should be unique                        |
| PLAYER_NUMBER         | Number of player you wanna play for (could be 0 or 1)                       |

Additional README about server side is on ./srv folder.
Have fun!! (^_^)