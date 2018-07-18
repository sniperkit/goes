package main

import (
	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/mvc"
	"github.com/sniperkit/iris/websocket"

	controllers "github.com/sniperkit/snk.golang.vuejs-multi-backend/controller/todo"
	models "github.com/sniperkit/snk.golang.vuejs-multi-backend/model/todo"
)

func setupWsMvcTodo(app *iris.Application) {
	// create a sub router and register the client-side library for the iris websockets,
	// you could skip it but iris websockets supports socket.io-like API.
	todosRouter := app.Party("/todos")
	// http://localhost:8080/todos/iris-ws.js
	// serve the javascript client library to communicate with
	// the iris high level websocket event system.
	todosRouter.Any("/ws/iris-ws.js", websocket.ClientHandler())

	// create our mvc application targeted to /todos relative sub path.
	todosApp := mvc.New(todosRouter)

	// any dependencies bindings here...
	todosApp.Register(
		models.NewMemoryService(),
		sess.Start,
		ws.Upgrade,
	)

	// controllers registration here...
	todosApp.Handle(new(controllers.TodoController))
}
