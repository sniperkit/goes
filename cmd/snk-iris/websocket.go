package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"

	// external
	"github.com/googollee/go-socket.io"
	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/websocket"
	xwebsocket "golang.org/x/net/websocket"
)

var (
	wsc = websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ws  *websocket.Server
	xws *xwebsocket.Conn
)

func setupSocketIO(app *iris.Application) {
	var transportNames []string
	server, err := socketio.NewServer(transportNames)
	if err != nil {
		app.Logger().Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) {
		app.Logger().Infof("on connection with socketio")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			app.Logger().Infof("emit: %v", so.Emit("chat message", msg))
			so.BroadcastTo("chat", "chat message", msg)
		})
		so.On("disconnection", func() {
			app.Logger().Infof("on disconnect")
		})

		// The return type may vary depending on whether you will return
		// In golang implementation of socket.io don't used callbacks for acknowledgement,
		// but used return value, which wrapped into ack package and returned to the client's callback in JavaScript
		so.On("snk:hello", func(msg string) string {
			app.Logger().Infof("emit: %v", so.Emit("snk says Hello !", msg))
			// so.BroadcastTo("snk:hello", "snk message", msg)
			return msg //Sending ack with data in msg back to client, using "return statement"
		})
		/*
			// You can use Emit or BroadcastTo with last parameter as callback for handling ack from client
			// Sending packet to room "room_name" and event "some:event"
			so.BroadcastTo("room_name", "some:event", dataForClient, func(so socketio.Socket, data string) {
				log.Println("Client ACK with data: ", data)
			})

			// Or
			so.Emit("some:event", dataForClient, func(so socketio.Socket, data string) {
				log.Println("Client ACK with data: ", data)
			})
		*/

	})

	server.On("error", func(so socketio.Socket, err error) {
		app.Logger().Errorf("error: %v", err)
	})

	// serve the socket.io endpoint.
	app.Any("/socket.io/{p:path}", iris.FromStd(server))

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"err":  model.NotFound,
			"msg":  "Not Found",
			"data": iris.Map{},
		})
	})

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"err":     model.ERROR,
			"message": "error",
			"data":    iris.Map{},
		})
	})
}

func setupWebsocket(app *iris.Application) {
	if config.Global.Websocket != nil {
		wsc = websocket.Config{
			ReadBufferSize:    config.Global.Websocket.ReadBufferSize,
			WriteBufferSize:   config.Global.Websocket.WriteBufferSize,
			BinaryMessages:    config.Global.Websocket.BinaryMessages,
			EnableCompression: config.Global.Websocket.EnableCompression,
			// Subprotocols:      config.Global.Websocket.Subprotocols,
			HandshakeTimeout: config.Global.Websocket.HandshakeTimeout,
			WriteTimeout:     config.Global.Websocket.WriteTimeout,
			ReadTimeout:      config.Global.Websocket.ReadTimeout,
			PongTimeout:      config.Global.Websocket.PongTimeout,
			PingPeriod:       config.Global.Websocket.PingPeriod,
		}
		ws = websocket.New(wsc)
	}
	ws.OnConnection(handleConnection)

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html, this endpoint is used to connect to the server.
	app.Get("/echo", ws.Handler())

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html, this endpoint is used to connect to the server.
	app.Get("/ws", ws.Handler())

	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})
}

func handleConnection(c websocket.Connection) {
	// Read events from browser
	c.On("snk", func(msg string) {
		// Print the message to the console, c.Context() is the iris's http context.
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner:
		// c.Emit("snk", msg)
		c.To(websocket.Broadcast).Emit("snk", msg)
	})
	// Read events from browser
	c.On("chat", func(msg string) {
		// Print the message to the console, c.Context() is the iris's http context.
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		// Write message back to the client message owner with:
		// c.Emit("chat", msg)
		// Write message to all except this client with:
		c.To(websocket.Broadcast).Emit("chat", msg)
	})
}

// ref. https://github.com/kataras/iris/blob/master/_examples/websocket/connectionlist/main.go
func setupFakeWebsocket(app *iris.Application) {
	Conn := make(map[websocket.Connection]bool)
	var userChatRoom = "snk-web"
	var mutex = new(sync.Mutex)

	ws.OnConnection(func(c websocket.Connection) {
		c.Join(userChatRoom)
		mutex.Lock()
		Conn[c] = true
		mutex.Unlock()
		c.On("snk", func(message string) {
			if message == "leave" {
				c.Leave(userChatRoom)
				c.To(userChatRoom).Emit("chat", "Client with ID: "+c.ID()+" left from the room and cannot send or receive message to/from this room.")
				c.Emit("snk", "You have left from the room: "+userChatRoom+" you cannot send or receive any messages from others inside that room.")
				return
			}
		})
		c.OnDisconnect(func() {
			mutex.Lock()
			delete(Conn, c)
			mutex.Unlock()
			fmt.Printf("\nConnection with ID: %s has been disconnected!\n", c.ID())
		})
	})

	var delay = 1 * time.Second
	go func() {
		i := 0
		for {
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()

	go func() {
		i := 0
		for range time.Tick(1 * time.Second) { //another way to get clock signal
			mutex.Lock()
			broadcast(Conn, fmt.Sprintf("aaaa2 %d\n", i))
			mutex.Unlock()
			time.Sleep(delay)
			i++
		}
	}()
}

type clientPage struct {
	Title string
	Host  string
}

func broadcast(Conn map[websocket.Connection]bool, message string) {
	for k := range Conn {
		k.To("snk-web").Emit("snk", message)
	}
}

//ConnectWebSocket connect a websocket to host
func ConnectWebSocket() error {
	var origin = "http://localhost/"
	var url = "ws://localhost:8080/socket"
	var err error
	xws, err = xwebsocket.Dial(url, "", origin)
	return err
}

// CloseWebSocket closes the current websocket connection
func CloseWebSocket() error {
	if xws != nil {
		return xws.Close()
	}
	return nil
}

// SendtBytes broadcast a message to server
func SendtBytes(serverID, to, method string, message []byte) error {
	// look https://github.com/sniperkit/iris/blob/master/websocket/message.go , client.go and client.js
	// to understand the buffer line:
	buffer := []byte(fmt.Sprintf("iris-websocket-message:%v;0;%v;%v;", method, serverID, to))
	buffer = append(buffer, message...)
	_, err := xws.Write(buffer)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////
// server side

// OnConnect handles incoming websocket connection
func OnConnect(c websocket.Connection) {
	fmt.Println("socket.OnConnect()")
	c.On("join", func(message string) { OnJoin(message, c) })
	c.On("objectupdate", func(message string) { OnObjectUpdated(message, c) })
	// ok works too c.EmitMessage([]byte("dsadsa"))
	c.OnDisconnect(func() { OnDisconnect(c) })

}

// OnJoin handles Join broadcast group request
func OnJoin(message string, c websocket.Connection) {
	t := time.Now()
	c.Join("server2")
	fmt.Println("OnJoin() time taken:", time.Since(t))
}

// OnObjectUpdated broadcasts to all client an incoming message
func OnObjectUpdated(message string, c websocket.Connection) {
	t := time.Now()
	s := strings.Split(message, ";")
	if len(s) != 3 {
		fmt.Println("OnObjectUpdated() invalid message format:" + message)
		return
	}
	serverID, _, objectID := s[0], s[1], s[2]
	err := c.To("server"+serverID).Emit("objectupdate", objectID)
	if err != nil {
		fmt.Println(err, "failed to broacast object")
		return
	}
	fmt.Println(fmt.Sprintf("OnObjectUpdated() message:%v, time taken: %v", message, time.Since(t)))
}

// OnDisconnect clean up things when a client is disconnected
func OnDisconnect(c websocket.Connection) {
	c.Leave("server2")
	fmt.Println("OnDisconnect(): client disconnected!")

}
