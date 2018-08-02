package main

import (
    "net/http"
    "log"
    "time"

    "github.com/gorilla/mux"
)

type Route struct {
  Name        string
  Method      string
  Path        string
  HandlerFunc http.HandlerFunc
}

var routes = [...] Route {
  Route { "Index", "GET", "/", index},
  Route { "GetFile", "GET", "/get", getfile},
  Route { "GetFile", "GET", "/get/{id}", getfile},
  Route { "PostResults", "POST", "/post", postresult},
  }

func Logger(myhandler http.Handler, name string) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    myhandler.ServeHTTP(w, r)

    log.Printf( "%s\t%s\t%s\t%s", r.Method, r.RequestURI, name,
      time.Since(start))
    })
}

func NewRouter() *mux.Router {
router := mux.NewRouter().StrictSlash(true)
for _, route := range routes {
  var handler http.Handler

  handler = route.HandlerFunc
  handler = Logger(handler, route.Name)

  router.Methods(route.Method).Path(route.Path).Name(route.Name).
    Handler(route.HandlerFunc)
  }
return router
}

