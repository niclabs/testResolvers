package server

import (
  "os"
  "log"
  "net/http"
  "github.com/niclabs/testResolvers/config"
)

func main() {
var cfg config.Configuration

err := config.ReadConfig("./" , &cfg) 
if err > 0 {  
  log.Fatal("Read Config")
  os.Exit(err)
  }  

router := NewRouter()
log.Fatal(http.ListenAndServe(":8080", router))
}


