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
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/controller/category"
	. "github.com/sniperkit/snk.golang.vuejs-multi-backend/logger"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/model"

	// external

	// iris - core
	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/context"
	"github.com/sniperkit/iris/middleware/logger"
	"github.com/sniperkit/iris/middleware/recover"
	"github.com/sniperkit/iris/sessions"
	"github.com/sniperkit/iris/sessions/sessiondb/badger"

	// iris - middleware
	corsmiddleware "github.com/sniperkit/iris-contrib-middleware/cors"
	jwtmiddleware "github.com/sniperkit/iris-contrib-middleware/jwt"
	ratemiddleware "github.com/sniperkit/iris-contrib-middleware/tollboothic"

	// external - 3rdparty
	"github.com/sniperkit/yaag/irisyaag"
	"github.com/sniperkit/yaag/yaag"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"

	"github.com/dgrijalva/jwt-go"
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
	cfg *config.Config
	app *iris.Application
	crs context.Handler
	ral *limiter.Limiter
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

	if config.Global.Server.Session != nil {

		db, err := badger.New(config.Global.Server.Session.DataDir)
		if err != nil {
			panic(err)
		}

		// close and unlock the database when control+C/cmd+C pressed
		iris.RegisterOnInterrupt(func() {
			db.Close()
		})

		defer db.Close() // close and unlock the database if application errored.

		sess = sessions.New(sessions.Config{
			Cookie:       config.Global.Server.Session.ID,
			Expires:      config.Global.Server.Session.Expires, // <=0 means unlimited life. Defaults to 0.
			AllowReclaim: config.Global.Server.Session.AllowReclaim,
		})

		// IMPORTANT:
		sess.UseDatabase(db)

	}

	app = iris.New()

	// init Server Application
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

	// init CORS middleware
	if config.Global.Api.Cors != nil {
		// set in server.yml file
		crs = corsmiddleware.New(corsmiddleware.Options{
			// AllowedOrigins is a list of origins a cross-domain request can be executed from.
			// If the special "*" value is present in the list, all origins will be allowed.
			// An origin may contain a wildcard (*) to replace 0 or more characters
			// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
			// Only one wildcard can be used per origin.
			// Default value is ["*"]
			AllowedOrigins: config.Global.Api.Cors.AllowedOrigins,

			// AllowedMethods is a list of methods the client is allowed to use with
			// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
			AllowedMethods: config.Global.Api.Cors.AllowedMethods,

			// AllowedHeaders is list of non simple headers the client is allowed to use with
			// cross-domain requests.
			// If the special "*" value is present in the list, all headers will be allowed.
			// Default value is [] but "Origin" is always appended to the list.
			AllowedHeaders: config.Global.Api.Cors.AllowedHeaders,

			// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
			// API specification
			// ExposedHeaders []string{},

			// MaxAge indicates how long (in seconds) the results of a preflight request
			// can be cached
			MaxAge: config.Global.Api.Cors.MaxAge,

			// AllowCredentials indicates whether the request can include user credentials like
			// cookies, HTTP authentication or client side SSL certificates.
			AllowCredentials: config.Global.Api.Cors.AllowCredentials,

			// OptionsPassthrough instructs preflight to let other potential next handlers to
			// process the OPTIONS method. Turn this on if your application handles OPTIONS.
			OptionsPassthrough: config.Global.Api.Cors.OptionsPassthrough,

			// Debugging flag adds additional output to debug server side CORS issues
			Debug: config.Global.Api.Cors.Debug,
		})
	}

	pp.Println("CORS: ", config.Global.Api.Cors)

	// init JWT middleware
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

	// init logger
	app.Use(logger.New())
	if *debugMode {
		// set logger debug level
		app.Logger().SetLevel("debug")
	}
	// init recover middleware
	app.Use(recover.New())

	// irisMiddleware := iris.FromStd(nativeTestMiddleware)
	// app.Use(irisMiddleware)

	// generate api documentation
	// note: do not use in production
	if config.Global.Api.Docs != nil {
		if config.Global.Api.Docs.Enabled {
			/*
				docBaseUrls := make(map[string]string, len(config.Global.Api.Docs.BaseUrls))
				for k, v := range config.Global.Api.Docs.BaseUrls {
					docBaseUrls[k] = v
				}
			*/
			yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
				On:       config.Global.Api.Docs.Enabled,
				DocTitle: config.Global.Api.Docs.DocTitle,
				DocPath:  "./shared/docs/apidoc.html",
				// DocPath:  filepath.Join(config.Global.Api.Path, config.Global.Api.Docs.DocPath, config.Global.Api.Docs.DocFile),
				BaseUrls: map[string]string{"Production": "", "Staging": ""}, // docBaseUrls,
			})

			app.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.
		}
	}

	// init rate limiter (contrib)
	ral = tollbooth.NewLimiter(5, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})

	app.Get("/rate/limiter", ratemiddleware.LimitHandler(ral), func(ctx iris.Context) {
		ctx.HTML("<b>Hello, rate limiter!</b>")
	})

	// setup routes defined in api.yml
	// note: extend with fake-api, mock-generator, data aggregator features...
	v1 := app.Party("/api/v1", crs).AllowMethods(iris.MethodOptions, iris.MethodPost, iris.MethodGet) // .AllowAll() // .AllowMethods(iris.MethodOptions) // <- important for the preflight.
	{
		// fake-api experimental stuff
		if errs := generateRoutesParty(v1); len(errs) != 0 {
			fmt.Printf("%d Error(s) in config: \n", len(errs))
			for i, err := range errs {
				fmt.Printf(" %d: %s\n", i+1, err.Error())
			}
		}

		// cms related experimental stuff
		// ref: github.com/liunian1004/goes
		v1.Get("/categories", nil)
		v1.Post("/category/create", category.Create)
		v1.Post("/category/update", category.Update)

		// ref: github.com/MuchChaca/Dashpanel
		// !!! to check !!!

		// ref: github.com/minhlucvan/gotodo
		// !!! to check !!!
	}

	// setup websocket/todo app routes (mvc)
	setupWsMvcTodo(app)

	// setup default websocket routes
	setupWebsocket(app)

	// setup socket.io routes
	setupSocketIO(app)

	// setup fake room via websocket
	setupFakeWebsocket(app)

	// create routes for static vue-element-admin
	app.StaticWeb("/admin", "./shared/dist/web/")

	// note: not working, do not redirect index.html to /admin but /admin/index.html works
	// app.StaticEmbeddedGzip("/admin", "shared/dist/web", GzipAsset, GzipAssetNames)

	// start test servers
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

	// print created endpoints
	routes := app.GetRoutes()
	for _, r := range routes {
		pp.Println(r.Name)
	}

	// init web-service
	address := iris.Addr(":" + strconv.Itoa(config.Global.Server.Port))
	if config.Global.Server.Env == model.DevelopmentMode {
		app.Run(address)
	} else {
		// app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
		app.Run(address, iris.WithoutVersionChecker)
	}
}
