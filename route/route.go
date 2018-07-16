package route

import (
	// external
	"github.com/kataras/iris"

	// internal
	"github.com/sniperkit/goes/controller/admin"
	"github.com/sniperkit/goes/controller/category"
)

const Prefix = "goes"

// Route 路由
func Route(app *iris.Application) {
	apiPrefix := Prefix

	router := app.Party(apiPrefix)
	{
		router.Get("/categories", nil)
	}

	adminRouter := app.Party(apiPrefix+"/admin", admin.Authentication)
	{
		adminRouter.Post("/category/create", category.Create)
		adminRouter.Post("/category/update", category.Update)
	}
}
