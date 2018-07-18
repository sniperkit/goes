package route

import (
	// external
	"github.com/sniperkit/iris"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/controller/admin"
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/controller/category"
)

// Route
func Route(apiPrefix string, app *iris.Application) {
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

// Party
func Party(apiPrefix string, route irouter.Party) {
	// router := app.Party(apiPrefix)
	//{
	router.Get("/categories", nil)
	//}

	adminRouter := app.Party(apiPrefix+"/admin", admin.Authentication)
	{
		adminRouter.Post("/category/create", category.Create)
		adminRouter.Post("/category/update", category.Update)
	}
}
