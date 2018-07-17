package main

import (
	// "errors"
	"net/http"
	"path/filepath"
	"sort"

	// external
	"github.com/sniperkit/iris"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
)

// const Prefix = "goes"

/*
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
*/

func nativeTestMiddleware(w http.ResponseWriter, r *http.Request) {
	println("Request path: " + r.URL.Path)
}

// generate routes for all configuration entries.
func generateRoutes(app *iris.Application) []error {
	var errs []error
	rootPath := config.Global.Api.Path

	endpoints := make([]*config.Endpoint, 0, len(config.Global.Api.URLs)+(len(config.Global.Api.Resources)*7))

	// validate and generate urls URLS
	for _, url := range config.Global.Api.URLs {
		u := url
		if u.File != "" {
			u.File = filepath.Join(*resPrefixPath, u.File) // to fix !!!
		}

		e, err := u.GetEndPoint(rootPath)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		endpoints = append(endpoints, e)
	}
	sort.Sort(config.Endpoints(endpoints))

	for _, en := range endpoints {
		app.Handle(en.Method, en.URL, en.Handler)
	}

	return errs
}

/*
// generate routes for all configuration entries.
func generateRoutes2(app *iris.Application) (*iris.Application, []error) {
    var errs []error
    rootPath := config.Global.Api.Path

    //  atleast one resource or url should present.
    if len(config.Global.Api.Resources) == 0 && len(config.Global.Api.URLs) == 0 {
        errs = append(errs, errors.New("Please provide atleast one resource or url"))
        return nil, errs
    }

    endpoints := make([]*config.Endpoint, 0, len(config.Global.Api.URLs)+(len(config.Global.Api.Resources)*7))

    // validate and generate urls URLS
    for _, url := range config.Global.Api.URLs {

        u := url
        e, err := u.GetEndPoint(rootPath)
        if err != nil {
            errs = append(errs, err)
            continue
        }
        endpoints = append(endpoints, e)
    }

    sort.Sort(config.Endpoints(endpoints))

    for _, en := range endpoints {

        // Method:   GET
        // Resource: http://localhost:8080/
        app.Handle(en.Method, en.URL, en.Handler)
        // app.Handle(en.URL, en.Handler).Methods(en.Method)
    }

    return app, errs
}
*/

/*
// router type is for generating routes from the configuration.
type router struct {
    c    *config.Config
    h    *ir.Router
    errs []error
}

// NewRouter create a new router from the configuration provided.
func newRouter(c *config.Config) *router {
    return &router{c: c, h: ir.NewRouter()}
}

// generate routes for all configuration entries.
func (r *router) generateRoutes() (*ir.Router, []error) {
    rootPath := r.c.Api.Path

    //  atleast one resource or url should present.
    if len(r.c.Api.Resources) == 0 && len(r.c.Api.URLs) == 0 {
        r.errs = append(r.errs, errors.New("Please provide atleast one resource or url"))
        return nil, r.errs
    }

    endpoints := make([]*config.Endpoint, 0, len(r.c.Api.URLs)+(len(r.c.Api.Resources)*7))

    // validate and generate urls URLS
    for _, url := range r.c.Api.URLs {
        u := url
        e, err := u.GetEndPoint(rootPath)
        if err != nil {
            r.errs = append(r.errs, err)
            continue
        }

        endpoints = append(endpoints, e)
    }

   if r.c.Api.JWT != nil {
       e, err := r.c.Api.JWT.GetEndPoint(rootPath)
       if err != nil {
           r.errs = append(r.errs, err)
           return nil, r.errs
       }

       endpoints = append(endpoints, e)
   }

   for _, re := range r.c.Api.Resources {
       res := re
       e, err := res.GetEndPoints(rootPath)
       if err != nil {
           r.errs = append(r.errs, err)
           continue
       }
       endpoints = append(endpoints, e...)
   }

    sort.Sort(config.Endpoints(endpoints))

    for _, en := range endpoints {
        r.h.Handle(en.URL, en.Handler).Methods(en.Method)
    }

   // static files -> to fix asap
   if r.c.Api.Static != nil {
       e, err := r.c.Api.Static.GetEndPoint(rootPath)
       if err != nil {
           r.errs = append(r.errs, err)
       }

       r.h.PathPrefix(r.c.Api.Static.URL).Handler(e.Handler)
   }

    return r.h, r.errs
}
*/
