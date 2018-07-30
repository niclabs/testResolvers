package main

import (
    "github.com/niclabs/testResolvers/resolvertests"
    "github.com/niclabs/testResolvers/config"
    "bufio"
    "encoding/csv"
    "runtime"
    "io"
    "log"
    "os"
)

func main() {
var cfg config.Configuration

if len(os.Args) < 2 {
  log.Fatal("[error] use: " + os.Args[0] + " csv_filename");
  os.Exit(-1)
  }

// TODO: put a valid directory here
err := config.ReadConfig("./" , &cfg) 
if err > 0 {  
  os.Exit(err)
  }  

ips := make (chan string, 20000)
res := make (chan resolvertests.Response, 20000)

for w:= 1; w <= runtime.NumCPU() * 2; w++ {
  go resolvertests.CheckDNS(w, ips, res)
  }

csvFile, _ := os.Open(os.Args[1])
reader := csv.NewReader(bufio.NewReader(csvFile))
for {
  line, error := reader.Read()
  if error == io.EOF {
    break
    } else if error != nil {
    log.Fatal(error)
    }
  ips <- line[0]
  }
for r :=  range res {
  log.Println(r)
  }
}
