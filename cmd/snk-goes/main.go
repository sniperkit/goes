package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
	. "github.com/sniperkit/snk.golang.vuejs-multi-backend/logger"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/route"

	// external
	"github.com/k0kubun/pp"

	"github.com/googollee/go-socket.io"
	xwebsocket "golang.org/x/net/websocket"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
	// _ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/middleware/logger"
	"github.com/sniperkit/iris/middleware/recover"
	"github.com/sniperkit/iris/websocket"
	// "github.com/sniperkit/iris/context"
	// "github.com/sniperkit/iris/mvc"
	// "github.com/sniperkit/iris/sessions"

	"github.com/dgrijalva/jwt-go"
	jwtmiddleware "github.com/sniperkit/iris-contrib-middleware/jwt"
)

const (
	appName         = "bindata"
	appVersionMajor = 1
	appVersionMinor = 3
	VERSION         = "1.3.0" // update later, incrmeent version with Makefile, pass value with -ldflags
)

var (
	//
	// AppVersionRev part of the program version.
	//
	// This will be set automatically at build time like so:
	//
	//     go build -ldflags "-X main.AppVersionRev `date -u +%s`" (go version < 1.5)
	//     go build -ldflags "-X main.AppVersionRev=`date -u +%s`" (go version >= 1.5)
	appVersionRev     string
	currentWorkDir, _ = os.Getwd()
	testServers       = flag.Bool("test-servers", false, "Test servers")

	configPrefixPath = flag.String("config-dir", currentWorkDir, "Config prefix path")
	configFiles      = []string{"config.yml"}
	// configFilename   = flag.String("config-file", "config.yaml", "Config filename")
	resDefaultDir = filepath.Join(currentWorkDir, "data")
	resPrefixPath = flag.String("resource-dir", resDefaultDir, "Resources prefix path")
)

var (
	app *iris.Application
	xws *xwebsocket.Conn
	cfg *config.Config
)

func initDB() {
	db, err := gorm.Open(config.Global.Store.Dialect, config.Global.Store.DSN)
	if err != nil {
		Error("open db connect error.")
		os.Exit(-1)
	}

	db.DB().SetMaxIdleConns(config.Global.Store.MaxIdleConns)
	db.DB().SetMaxOpenConns(config.Global.Store.MaxOpenConns)

	// global
	model.DB = db
}

func prettyPrint(msg interface{}) {
	pp.Println(msg)
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

func main() {

	fmt.Printf("SNK-API - fake rest api server (%s) \n", VERSION)

	flag.Var((*AppendSliceValue)(&configFiles), "config-file", "Regex pattern to ignore")
	flag.Parse()

	Log(*configPrefixPath)
	Log(configFiles)

	var configPrefixedFiles []string
	for _, v := range configFiles {
		// todo: check if only filename, either prefix with default prefix path...
		configPrefixedFiles = append(configPrefixedFiles, fmt.Sprintf("%s/%s", *configPrefixPath, v))
	}

	Log(configPrefixedFiles)
	config.Global = config.New(configPrefixedFiles...)

	initDB()

	pp.Println("config: ", config.Global)

	ws := websocket.New(websocket.Config{
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
		BinaryMessages:    false,
		EnableCompression: false,
		Subprotocols:      []string{},
		// IDGenerator: "",
		// CheckOrigin: "",
		// HandshakeTimeout:        time.Duration(),
		// WriteTimeout:        time.Duration(),
		// ReadTimeout:        time.Duration(),
		// PongTimeout:        time.Duration(),
		// PingPeriod:        time.Duration(),
	})

	ws.OnConnection(handleConnection)

	app = iris.New()

	if strings.ToLower(config.Global.Server.Engine) == "iris" && config.Global.Server.Settings.Iris != nil {
		app.Configure(iris.WithConfiguration(
			iris.Configuration{
				DisableStartupLog:                 config.Global.Server.Settings.Iris.DisableStartupLog,
				DisableInterruptHandler:           config.Global.Server.Settings.Iris.DisableInterruptHandler,
				DisablePathCorrection:             config.Global.Server.Settings.Iris.DisablePathCorrection,
				EnablePathEscape:                  config.Global.Server.Settings.Iris.EnablePathEscape,
				FireMethodNotAllowed:              config.Global.Server.Settings.Iris.FireMethodNotAllowed,
				DisableBodyConsumptionOnUnmarshal: config.Global.Server.Settings.Iris.DisableBodyConsumptionOnUnmarshal,
				DisableAutoFireStatusCode:         config.Global.Server.Settings.Iris.DisableAutoFireStatusCode,
				TimeFormat:                        config.Global.Server.Settings.Iris.TimeFormat,
				Charset:                           config.Global.Server.Settings.Iris.Charset,
			}),
		)
	}

	jwtHandler := jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte("My Secret"), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	})

	app.Use(jwtHandler.Serve)

	app.Use(logger.New())
	// use this recover(y) middleware
	app.Use(recover.New())

	route.Route(app)

	// 测试模式
	//if config.Global.Server.Env == model.DevelopmentMode {
	//	app.Adapt(iris.DevLogger())
	//}

	//app.Adapt(sessions.New(sessions.Config{
	//	Cookie: config.Global.Server.SessionID,
	//	Expires: time.Minute * 20,
	//}))

	//app.Adapt(httprouter.New())

	// register the server on an endpoint.
	// see the inline javascript code in the websockets.html, this endpoint is used to connect to the server.
	app.Get("/ws", ws.Handler())

	// serve the javascript built'n client-side library,
	// see websockets.html script tags, this path is used.
	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})

	// http://localhost:8080/todos/iris-ws.js
	// serve the javascript client library to communicate with
	// the iris high level websocket event system.
	// app.Any("/iris-ws.js", websocket.ClientHandler())

	serverSocketIO, err := newSocketIO(nil)
	if err != nil {
		Error("create socket.io error.")
		os.Exit(-1)
	}

	// serve the socket.io endpoint.
	app.Any("/socket.io/{p:path}", iris.FromStd(serverSocketIO))

	// serve the index.html and the javascript libraries at
	// http://localhost:8080
	// app.StaticWeb("/", "./public")

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

	// serve our app in public, public folder
	// contains the client-side vue.js application,
	// no need for any server-side template here,
	// actually if you're going to use vue without any
	// back-end services, you can just stop after this line and start the server.
	app.StaticWeb("/admin", "./shared/dist/web/")
	// app.StaticWeb("/admin/js", "./shared/dist/web/js")

	// $ go get -u github.com/jteeuwen/go-bindata/...
	// go-bindata -pkg main -o ./cmd/snk-goes/bindata.go ./shared/dist/web/...
	// app.StaticWeb("/admin", "./static")
	// app.StaticEmbedded("/admin", "./shared/dist/web/", Asset, AssetNames)

	// $ go get -u github.com/kataras/bindata/cmd/bindata
	// bindata -pkg embedded -o ./embedded/bindata-gz.go ./shared/dist/web/...
	// bindata -pkg main -o ./cmd/snk-goes/bindata.go ./shared/dist/web/...

	/*
		Strange behavior of app.StaticEmbeddedGzip. It cannot detech and render index page automatically
		http://localhost:8080/index.html works
		but
		http://localhost:8080/ not found
	*/

	/*
		I propose another bad trick to write content of index.html when request is /
		Chrome, FireFox render correctly but Safari does not work at all
	*/
	/*
		app.Get("/", func (ctx iris.Context) {
			if data, err := GzipAsset("assets/index.html"); err != nil {
				ctx.StatusCode(http.StatusInternalServerError)
				ctx.WriteString("index.html is not found")
				return
			} else {
				ctx.StatusCode(http.StatusOK)
				ctx.Header("Content-Encoding", "gzip")
				ctx.ContentType("text/html")
				ctx.WriteGzip(data)
			}
		})
	*/
	// app.StaticEmbeddedGzip("/admin", "shared/dist/web", GzipAsset, GzipAssetNames)

	/*
		app.Get("/admin/{p:path}", func(ctx iris.Context) {
			context.AddGzipHeaders(ctx.ResponseWriter())
			ctx.ContentType("text/html")
			ctx.Write(_gzipBindataWebindexhtml)
		})
	*/

	// app.StaticWeb("/admin", "./shared/dist/web/")
	// app.StaticWeb("/admin/js", "./shared/dist/web/js")

	// $ go get -u github.com/shuLhan/go-bindata/...
	// $ go-bindata ./templates/...
	// $ go build
	// $ ./embedding-templates-into-app
	// html files are not used, you can delete the folder and run the example.
	// tmpl.Binary(views.Asset, views.AssetNames) // <-- IMPORTANT

	// app.RegisterView(tmpl)

	/*
		// // DASH
		// create a sub router an register the client-side library for the iris websockets,
		// you could skip it but iris websockets supports socket.io-like API.
		dashRouter := app.Party("/dash")
		// http://localhost:8080/todos/iris-ws.js
		// serve the javascript client library to communicate with
		// the iris high level websocket event system.
		dashRouter.Any("/iris-ws.js", websocket.ClientHandler())

		//create our mvc app targeted to /dash relative sub path.
		dashApp := mvc.New(dashRouter)

		// any dependencies bindings here . . .
		dashApp.Register(
			// dash.NewMemoryService(),
			// sess.Start,
			ws.Upgrade,
		)

		// controllers registration here
		dashApp.Handle(new(controller.ProcessController))

		// app.Party("/process", )

		// PartyFunc is working ! Yay !
		app.PartyFunc("/process", func(process iris.Party) {
			// users.Use(myAuthMiddlewareHandler)

			// http://localhost:8080/users/42/profile
			process.Put("/start/{name:string}", controller.StartProcess)
			// http://localhost:8080/users/messages/1
			// users.Get("/inbox/{id:int}", userMessageHandler)
		})

		// process := app.Party("/process")
		// app.PartyFunc("/process", func (process iris.Party)) {
		//	process.Put("/start/{name:string}", controller.ProcessController.StartProcess())
		// })
	*/

	// create our mvc app targeted to /dash relative sub path.
	// snkApp := mvc.New(dashRouter)

	if *testServers {
		// start a secondary server listening on localhost:9090.
		// use "go" keyword for Listen functions if you need to use more than one server at the same app.
		//
		// http://localhost:9090/
		// http://localhost:9090/mypath
		srv1 := &http.Server{Addr: ":9090", Handler: app}
		go srv1.ListenAndServe()
		println("Start a server listening on http://localhost:9090")

		// start a "second-secondary" server listening on localhost:5050.
		//
		// http://localhost:5050/
		// http://localhost:5050/mypath
		srv2 := &http.Server{Addr: ":5050", Handler: app}
		go srv2.ListenAndServe()
		println("Start a server listening on http://localhost:5050")
	}

	address := iris.Addr(":" + strconv.Itoa(config.Global.Server.Port))

	if config.Global.Server.Env == model.DevelopmentMode {
		app.Run(address)
	} else {
		app.Run(address, iris.WithoutVersionChecker)
	}
}
