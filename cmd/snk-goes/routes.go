package main

import (
	"errors"
	"sort"

	// external
	"github.com/gorilla/mux"

	// internal
	"github.com/sniperkit/snk.golang.vuejs-multi-backend/config"
)

// router type is for generating routes from the configuration.
type router struct {
	c    *config.Config
	h    *mux.Router
	errs []error
}

// NewRouter create a new router from the configuration provided.
func newRouter(c *config.Config) *router {
	return &router{c: c, h: mux.NewRouter()}
}

// generate routes for all configuration entries.
func (r *router) generateRoutes() (*mux.Router, []error) {
	rootPath := r.c.Api.Path

	//  atleast one resource or url should present.
	if len(r.c.Api.Resources) == 0 && len(r.c.Api.URLs) == 0 {
		r.errs = append(r.errs, errors.New("Please provide atleast one resource or url"))
		return nil, r.errs
	}

	endpoints := make([]*config.Endpoint, 0, len(r.c.Api.URLs)+(len(r.c.Api.Resources)*7))

	/*
	   if r.c.Api.JWT != nil {
	       e, err := r.c.Api.JWT.GetEndPoint(rootPath)
	       if err != nil {
	           r.errs = append(r.errs, err)
	           return nil, r.errs
	       }

	       endpoints = append(endpoints, e)
	   }
	*/

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

	/*
	   // static files -> to fix asap
	   if r.c.Api.Static != nil {
	       e, err := r.c.Api.Static.GetEndPoint(rootPath)
	       if err != nil {
	           r.errs = append(r.errs, err)
	       }

	       r.h.PathPrefix(r.c.Api.Static.URL).Handler(e.Handler)
	   }
	*/

	return r.h, r.errs
}
