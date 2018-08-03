package main

import (
    "github.com/niclabs/testResolvers/resolvertests"
    "github.com/niclabs/testResolvers/config"
    "runtime"
    "time"
    "log"
    "bytes"
    "encoding/json"
)

func run (cfg config.Configuration, c Communication) {

// Get INFO

iplist, err := c.Get(cfg)
if err != nil {
  log.Fatal(err)
  }

// Parallel process

ips := make (chan string)
res := make (chan resolvertests.Response)

for w:= 1; w <= runtime.NumCPU(); w++ {
  go resolvertests.CheckDNS(w, ips, res)
  }

// send ips to workers
for _,ip := range iplist {
  ips <- ip
  }

// receive test results

m := struct {
  time int64 
  login string
  location string
  res [] resolvertests.Response
  } {
  time.Now().Unix(),
  cfg.Login,
  cfg.Location,
  make ([]resolvertests.Response,0),
  }

log.Println(iplist)

for r:=0 ; r < len(iplist) ; r++ {
  m.res = append(m.res,<-res)
  }

close(ips)
close(res)

b := new(bytes.Buffer)
log.Println(b)

json.NewEncoder(b).Encode(m)

err = c.Post (cfg, b)
if err != nil {
  log.Fatal("Error Posting IP list ")
  }
}

func main() {
var cfg config.Configuration
var c Communication = REST {}

// TODO: put a valid directory here
errno := config.ReadConfig("./" , &cfg) 
if errno > 0 {  
  log.Fatal("Error reading config file")
  }  

run(cfg,c)

}
