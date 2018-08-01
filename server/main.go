package server

import (
  "os"
  "fmt"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "time"
  "github.com/niclabs/testResolvers/config"
)

// ROUTER
type Route struct {
  Name        string
  Method      string
  Path        string
  HandlerFunc http.HandlerFunc
}

var routes = [...] Route {
  Route { "Index", "GET", "/", index},
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

// HANDLERS

func index(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "This is an API (wow!)")
}

func getfile(w http.ResponseWriter, r *http.Request) {
//params := mux.Vars(r)

}

func postresult(w http.ResponseWriter, r *http.Request) {
//params := mux.Vars(r)

}

// MAIN
func main() {
var cfg config.Configuration

err := config.ReadConfig("./" , &cfg) 
if err > 0 {  
  os.Exit(err)
  }  

router := NewRouter()
log.Fatal(http.ListenAndServe(":" + cfg.Port, router))
}


