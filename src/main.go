package main

import (
    "resolvertests"
    "readConfig"
    "bufio"
    "encoding/csv"
    "fmt"
    "io"
    "log"
    "os"
)

func main() {
var config Configuration

if len(os.Args) < 2 {
  log.Fatal("[error] use: " + os.Args[0] + " csv_filename");
  os.Exit(-1)
  }

file, err := os.Open("config.json") 
if err != nil {  
  log.Fatal("[error] loading config.json")
  os.Exit(-2)
  }  
decoder := json.NewDecoder(file) 
err = decoder.Decode(&config) 
if err != nil {  
  log.Fatal("[error] decoding config.json")
  os.Exit(-3)
  }
defer file.Close()

ips := make (chan string, 20000)
res := make (chan resolvertests.Response, 20000)

for w:= 1; w <= 1; w++ {
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
