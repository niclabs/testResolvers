package main

import (
    "github.com/niclabs/testResolvers/resolvertests"
    "github.com/niclabs/testResolvers/config"
    "encoding/json"
    "net/http"
    "runtime"
    "time"
    "log"
    "io/ioutil"
    "os"
)

func main() {
var cfg config.Configuration
var iplist []string

// TODO: put a valid directory here
errno := config.ReadConfig("./" , &cfg) 
if errno > 0 {  
  log.Fatal("Error reading config file")
  os.Exit(errno)
  }  

ips := make (chan string, 20000)
res := make (chan resolvertests.Response, 20000)

for w:= 1; w <= runtime.NumCPU(); w++ {
  go resolvertests.CheckDNS(w, ips, res)
  }

// http request 

cli := http.Client{
  Timeout: time.Second * 10, // Maximum of 2 secs
}

url := "http://" + cfg.Server + ":" + cfg.Port + "/get"

req, err := http.NewRequest(http.MethodGet, url , nil)
if err != nil {
  log.Fatal("Error with request")
  os.Exit(-2)
  }
// TODO:  add auth on headers
req.Header.Set("User-Agent", "testresolver client")
httpResp, err := cli.Do(req)
  if err != nil {
    log.Fatal("Error performing http get")
    os.Exit(-3)
}

// process data
data, _ := ioutil.ReadAll(httpResp.Body)
err = json.Unmarshal([]byte(data), &iplist)
if err != nil {
  log.Fatal("Error Unmarshaling IP list")
  os.Exit(-3)
  }

for _,ip := range iplist {
  ips <- ip
  }
for r :=  range res {
  log.Fatal(r)
  }
}
