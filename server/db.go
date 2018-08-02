package main

import (
    "encoding/csv"
    "os"
    "log"
    "io"
    "bufio"
)

func ReadData() ([]string, int) {
ips := make ([]string,0)
csvFile, _ := os.Open(cfg.CSVfile)
reader := csv.NewReader(bufio.NewReader(csvFile))
for {
  line, error := reader.Read()
  if error == io.EOF {
    break
    } else if error != nil {
    log.Fatal(error)
    return nil,1
    }
  ips = append (ips,line[0])
  }
return ips,0
}  

/*
func StoreData([]Response) int {
}
*/
