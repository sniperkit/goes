package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/sniperkit/iris/websocket"

	"github.com/googollee/go-socket.io"
	xwebsocket "golang.org/x/net/websocket"
)

var (
	xws *xwebsocket.Conn
)

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

func newSocketIO(transportNames []string) (*socketio.Server, error) {
	server, err := socketio.NewServer(transportNames)
	if err != nil {
		return nil, err
	}
	server.On("connection", func(so socketio.Socket) {
		app.Logger().Infof("on connection")
		so.Join("snk-io")
		so.On("snk-io message", func(msg string) {
			app.Logger().Infof("emit: %v", so.Emit("snk-io message", msg))
			so.BroadcastTo("snk-io", "snk-io message", msg)
		})
		so.On("disconnection", func() {
			app.Logger().Infof("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		app.Logger().Errorf("error: %v", err)
	})
	return server, nil
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
}
