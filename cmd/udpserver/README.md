## Motivation
Here is a part of the hockey game - it's a server!
It is a simple application that is responsible for storing and serving world state for clients, so, it is somewhere in the middle between two clients of the game.

## The main idea
The application could work in two modes:
- Backend-only: means that all calculations will be done exactly on the server-side. In this case, clients will have more or less the same status (there will be no one who actually "creates a game" 'cause all data and calculations are on backend side)
- webRTC mode: in this case, this server just plays the role of the signal server, and the first connected client will actually create a game and will handle all calculations on his side.

The game works in a very simple way:
1. After the client joined, we start a goroutine, that will contain all information about the game, and actually, will be the game itself. This way has a few advantages:
    - There are no race conditions. Goroutine serves the data in a pretty similar way to how a small web server would do that
    - Game will be easily stopped and all data will be cleaned up after the goroutine stopped
    - We could apply "time" changes to the world in the internal goroutine cycle without blocking other parts and workflows of the application
2. After some defined period of time, we move the world according to the speed of every object (speed itself coming from the client-side)
3. If we have a regular update from client-side, we do an update on the server-side (speed, position, etc.), do calculations (check collisions, goals), and change the input object vector to notify the client

## Components:
- server_http.go: responsible for receiving and handling HTTP requests (currently it is only "join game" and "health check" endpoints)
- server_udp.go: responsible for handling UDP data. The one disadvantage is that currently one application server both HTTP & UDP data, which is not really scalable way. Could be improved in the future if this project will grow
- processor.go: the component responsible for processing the request (will use necessary components as well and put them all together)
- handler.go: contains the definition of HTTP handlers 
- game_instance.go: contains the game instance definition and its logic, probably of the most important parts of the application
- compress.go: responsible for compressing and decompressing input/output data to improve network performance
- k8s: contains a simple k8s local setup, that gives an ability to deploy the application itself with the NGINX  ingress controller 