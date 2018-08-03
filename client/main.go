package main

import (
    "github.com/niclabs/testResolvers/resolvertests"
    "github.com/niclabs/testResolvers/config"
    "runtime"
    "log"
)

func run (cfg config.Configuration, c Communication) {

// Get INFO

iplist, err := c.Get(cfg)
if err != nil {
  log.Fatal(err.Error())
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

var reslice []resolvertests.Response
for r:=0 ; r < len(iplist) ; r++ {
  reslice = append(reslice,<-res)
  }

close(ips)
close(res)

// publish them

err = c.Post (cfg, reslice)
if err != nil {
  log.Fatal("Error POST call " + err.Error())
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
