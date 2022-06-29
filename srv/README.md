## Motivation
Here is hear of the hockey game - it's server.
It is a simple application that is responsible for storing and serving world state for  a clients, so, it is somewhere in the middle between two client of the game.

## The main idea
Application could work in a two modes:
- Backend-only: means that all calculations will be done exactly on the server side. In this case clients will have more or less the same status (there will be no one who actually "create a game" 'cause all data and calculations is on backend side)
- webRTC mode: in this case this server just play a role of the signal server, and the first connected client will actually create a game and will handle all calculations on his side.

The game works in a very simple way:
1. After client joined, we start a goroutine, that will contain all information about the game, and actually, will be a game itself. This way have a few advantages:
    - There is no race conditions. Goroutine serve the data in a pretty similar way how small web server would do that
    - Game will be easily stopped and all data will be cleaned up after goroutine stopped
    - We could apply "time" changes to the world in the internal goroutine cycle without blocking other parts and workflows of the application
2. After some defined period of time, we move the world according to the speed of every object (speed itself coming from the client side)
3. If we have a regular update from client side, we do an update on server side (speed, position, etc.), do a calculations (check collisions, goals), and change input object vector to notify the client

## Components:
- server_http.go: responsible for receiving and handling HTTP requests (currently it is join game and health check endpoints)
- server_udp.go: responsible for handling UDP data. The one disadvantage is that currently one application server both HTTP & UDP data, what is not really scalable way. Could be improved in the future if this project will grow
- processor.go: component responsible for processing the request (will use necessary components as well and put them all together)
- handler.go: contains the definition of HTTP handlers 
- game_instance.go: contains the game instance definition and its logic, probably of the most important parts of the application
- compress.go: responsible for compressing and decompressing input/output data to improve network performance
- k8s: contains simple k8s local setup, that gives an ability to deploy application itself with the NGINX  ingress controller 