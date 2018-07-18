package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
	. "github.com/sniperkit/snk.golang.vuejs-multi-backend/logger"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"
	// "github.com/sniperkit/snk.golang.vuejs-multi-backend/route"

	// external
	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/middleware/logger"
	"github.com/sniperkit/iris/middleware/recover"
	"github.com/sniperkit/iris/websocket"
	// "github.com/sniperkit/iris/context"
	// "github.com/sniperkit/iris/mvc"
	// "github.com/sniperkit/iris/sessions"

	"github.com/dgrijalva/jwt-go"
	corsmiddleware "github.com/sniperkit/iris-contrib-middleware/cors"
	jwtmiddleware "github.com/sniperkit/iris-contrib-middleware/jwt"

	// debug
	"github.com/k0kubun/pp"
)

const (
	appName         = "bindata"
	appVersionMajor = 1
	appVersionMinor = 3
	VERSION         = "1.3.0" // update later, incrmeent version with Makefile, pass value with -ldflags
)

var (
	// AppVersionRev part of the program version.
	// This will be set automatically at build time like so:
	//     go build -ldflags "-X main.AppVersionRev `date -u +%s`" (go version < 1.5)
	//     go build -ldflags "-X main.AppVersionRev=`date -u +%s`" (go version >= 1.5)
	appVersionRev     string
	currentWorkDir, _ = os.Getwd()
	jwtMode           = flag.Bool("with-jwt", false, "JWT Authenticate mode")
	debugMode         = flag.Bool("debug", false, "Debug mode")
	testServers       = flag.Bool("test-servers", false, "Test servers")
	configPrefixPath  = flag.String("config-dir", currentWorkDir, "Config prefix path")
	configFiles       = []string{"application.yml", "api.yml", "server.yml", "websocket.yml", "database.yml"}
	resDefaultDir     = filepath.Join(currentWorkDir, "data")
	resPrefixPath     = flag.String("resource-dir", resDefaultDir, "Resources prefix path")
)

var (
	app *iris.Application
	cfg *config.Config
)

func main() {

	fmt.Printf("SNK-API - fake rest api server (%s) \n", VERSION)
	fmt.Printf("SNK-API - currentWorkDir: %s\n", currentWorkDir)

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

	if *debugMode {
		pp.Println("config.App", config.Global.App)
		pp.Println("config.Api", config.Global.Api)
		pp.Println("config.Websocket", config.Global.Websocket)
		pp.Println("config.Store", config.Global.Store)
		pp.Println("config.Server", config.Global.Server)
	}

	var wsCfg websocket.Config
	if config.Global.Websocket != nil {
		wsCfg = websocket.Config{
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
	}

	ws := websocket.New(wsCfg)
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

	// set in server.yml file
	crs := corsmiddleware.New(corsmiddleware.Options{
		// AllowedOrigins is a list of origins a cross-domain request can be executed from.
		// If the special "*" value is present in the list, all origins will be allowed.
		// An origin may contain a wildcard (*) to replace 0 or more characters
		// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
		// Only one wildcard can be used per origin.
		// Default value is ["*"]
		AllowedOrigins: []string{
			"http://localhost:9200",
			"http://localhost:7474",
			"http://localhost:8080",
			"http://localhost:3000",
			"http://localhost:9528",
		}, // allows everything, use that to change the hosts.

		// AllowedMethods is a list of methods the client is allowed to use with
		// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
		AllowedMethods: []string{"HEAD", "GET", "POST", "OPTIONS"},

		// AllowedHeaders is list of non simple headers the client is allowed to use with
		// cross-domain requests.
		// If the special "*" value is present in the list, all headers will be allowed.
		// Default value is [] but "Origin" is always appended to the list.
		AllowedHeaders: []string{"Access-Control-Allow-Origin", "X-Auth-Token", "X-Token"},

		// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
		// API specification
		// ExposedHeaders []string{},

		// MaxAge indicates how long (in seconds) the results of a preflight request
		// can be cached
		MaxAge: 3600,

		// AllowCredentials indicates whether the request can include user credentials like
		// cookies, HTTP authentication or client side SSL certificates.
		AllowCredentials: true,

		// OptionsPassthrough instructs preflight to let other potential next handlers to
		// process the OPTIONS method. Turn this on if your application handles OPTIONS.
		OptionsPassthrough: false,

		// Debugging flag adds additional output to debug server side CORS issues
		Debug: true,
	})

	// pp.Println("cors: ", crs)

	if *jwtMode {
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
	}

	app.Use(logger.New())
	app.Use(recover.New())

	// irisMiddleware := iris.FromStd(nativeTestMiddleware)
	// app.Use(irisMiddleware)

	// route.Route(app)

	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions, iris.MethodPost, iris.MethodGet) // .AllowAll() // .AllowMethods(iris.MethodOptions) // <- important for the preflight.
	{
		if errs := generateRoutesParty(v1, iris.MethodOptions); len(errs) != 0 {
			fmt.Printf("%d Error(s) in config: \n", len(errs))
			for i, err := range errs {
				fmt.Printf(" %d: %s\n", i+1, err.Error())
			}
			os.Exit(1)
		}
	}

	/*
		v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodGet) // <- important for the preflight.
		{
			if errs := generateRoutesParty(v1, iris.MethodGet); len(errs) != 0 {
				fmt.Printf("%d Error(s) in config: \n", len(errs))
				for i, err := range errs {
					fmt.Printf(" %d: %s\n", i+1, err.Error())
				}
				os.Exit(1)
			}
		}

		v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodPost) // <- important for the preflight.
		{
			if errs := generateRoutesParty(v1, iris.MethodPost); len(errs) != 0 {
				fmt.Printf("%d Error(s) in config: \n", len(errs))
				for i, err := range errs {
					fmt.Printf(" %d: %s\n", i+1, err.Error())
				}
				os.Exit(1)
			}
		}
	*/

	routes := app.GetRoutes()
	for _, r := range routes {
		pp.Println(r.Name)
	}

	/*
		var errs []error
		app, errs = generateRoutes(app)
		if len(errs) != 0 {
			fmt.Printf("%d Error(s) in config: \n", len(errs))
			for i, err := range errs {
				fmt.Printf(" %d: %s\n", i+1, err.Error())
			}
			os.Exit(1)
		}
	*/

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
