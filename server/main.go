package main

import (
  "os"
  "log"
  "net/http"
  "time"
  "github.com/niclabs/testResolvers/config"
)

var IPlist []string
var cfg config.Configuration

func main() {

err := config.ReadConfig("./" , &cfg) 
if err > 0 {  
  log.Fatal("Read Config")
  os.Exit(err)
  }  

IPlist, err = ReadData()
if err > 0 {
  log.Fatal("Read Data")
  os.Exit(err)
  }

router := NewRouter()

s := &http.Server{
	Addr:           ":" + cfg.Port,
        Handler:        router,
	ReadTimeout:    10 * time.Second,
	WriteTimeout:   10 * time.Second,
	MaxHeaderBytes: 1 << 20,
}
log.Fatal(s.ListenAndServe())
}


