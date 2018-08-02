package main

import (
  "os"
  "log"
  "net/http"
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
log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}


