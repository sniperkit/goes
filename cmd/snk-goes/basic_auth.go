package main

import (
	"time"

	"github.com/sniperkit/iris"
	"github.com/sniperkit/iris/middleware/basicauth"
)

func newBasicAuth(app *iris.Application) *iris.Application {
	authConfig := basicauth.Config{
		Users: map[string]string{
			"myusername":       "mypassword",
			"mySecondusername": "mySecondpassword",
		},
		Realm:   "Authorization Required", // defaults to "Authorization Required"
		Expires: time.Duration(30) * time.Minute,
	}
	authentication := basicauth.New(authConfig)

	// to global app.Use(authentication) (or app.UseGlobal before the .Run)
	// to routes
	/*
	   app.Get("/mysecret", authentication, h)
	*/
	app.Get("/", func(ctx iris.Context) { ctx.Redirect("/admin") })

	// to party
	needAuth := app.Party("/admin", authentication)
	{
		//http://localhost:8080/admin
		needAuth.Get("/", h)
		// http://localhost:8080/admin/profile
		needAuth.Get("/profile", h)
		// http://localhost:8080/admin/settings
		needAuth.Get("/settings", h)
	}

	return app
}

func h(ctx iris.Context) {
	username, password, _ := ctx.Request().BasicAuth()
	// third parameter it will be always true because the middleware
	// makes sure for that, otherwise this handler will not be executed.
	ctx.Writef("%s %s:%s", ctx.Path(), username, password)
}
