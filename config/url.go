package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	// external
	"github.com/sniperkit/iris"
)

/*
// URL type represent a url resource in rest api
// Url is mandatory.
// ContentType defaults to application/json if not provided.
// Status code default to 200
// Method defaults to GET if not provided.
type URL struct {
    URL         string            `json:"url"`
    Method      string            `json:"method,omitempty"`
    ContentType string            `json:"content_type,omitempty"`
    File        string            `json:"file,omitempty"`
    StatusCode  int               `json:"status,omitempty"`
    Headers     map[string]string `json:"headers,omitempty"`
}
*/

// Sanitize validate and sanitize the URL resource
// add default values for the optional fields
// check file is valid and return error if not
func (u *URL) sanitize(root string) error {
	if u.Method == "" {
		u.Method = "GET"
	}

	if u.ContentType == "" {
		u.ContentType = "application/json; charset=utf-8"
	}

	if u.StatusCode == 0 {
		u.StatusCode = 200
	}

	if !strings.HasPrefix(u.URL, "/") {
		u.URL = "/" + u.URL
	}

	if u.File != "" {
		if _, err := os.Stat(filepath.Join(root, u.File)); err != nil {
			return fmt.Errorf("Invalid file '%s' in URL '%s'", u.File, u.URL)
		}
	}

	return nil
}

// GetEndPoint send an endpoint for the URL
func (u *URL) GetEndPoint(rootPath string) (*Endpoint, error) {
	if err := u.sanitize(rootPath); err != nil {
		return nil, err
	}
	return &Endpoint{URL: u.URL, Method: u.Method, Handler: u.getHandle(rootPath)}, nil
}

func (u *URL) getHandle(root string) func(ctx iris.Context) {
	if u.File == "" {
		return func(ctx iris.Context) {
			ctx.ContentType(u.ContentType)
			for key, value := range u.Headers {
				if key == "Content-type" {
					continue
				}
				ctx.Header(key, value)
			}
			ctx.StatusCode(u.StatusCode)
			ctx.Write([]byte(""))
		}
	}
	return func(ctx iris.Context) {
		data, err := os.Open(filepath.Join(root, u.File))
		if err != nil {
			ctx.NotFound()
			return
		}
		defer data.Close()

		ctx.ContentType(u.ContentType)
		for key, value := range u.Headers {
			if key == "Content-type" {
				continue
			}
			ctx.Header(key, value)
		}
		ctx.StatusCode(u.StatusCode)

		/*
		   fileInfo, err := data.Stat()
		   if err != nil {
		       fmt.Println(err)
		   }
		*/
		dataBytes, err := ioutil.ReadAll(data)
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println("SIZE1:", fileInfo.Size())
		// fmt.Println("SIZE2:", len(dataBytes))
		// io.Copy(w, data)
		ctx.Write(dataBytes)
	}
}

/*
// getHandle generate a static handle for the URL.
func (u *URL) getHandle(root string) http.Handler {
	if u.File == "" {
		// func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", u.ContentType)
			for key, value := range u.Headers {
				if key == "Content-type" {
					continue
				}
				w.Header().Set(key, value)
			}

			w.WriteHeader(u.StatusCode)
			w.Write([]byte(""))
		})
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := os.Open(filepath.Join(root, u.File))
		if err != nil {
			w.WriteHeader(404)
			return
		}

		defer data.Close()
		w.Header().Set("Content-type", u.ContentType)
		for key, value := range u.Headers {
			if key == "Content-type" {
				continue
			}
			w.Header().Set(key, value)
		}
		w.WriteHeader(u.StatusCode)
		io.Copy(w, data)
	})
}
*/
