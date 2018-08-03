package main

import (
    "fmt"
    "encoding/json"
    "net/http"
    "log"
    "io"
    "io/ioutil"

    "github.com/niclabs/testResolvers/resolvertests"
)

type Query struct {
ips []string
}

type Reply struct {
time string
login string
location string
res [] resolvertests.Response
}

func index(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "This is an API (wow!)")
}

func getfile(w http.ResponseWriter, r *http.Request) {

w.Header().Set("Content-Type", "application/json; charset=UTF-8")
w.WriteHeader(http.StatusOK)
if err := json.NewEncoder(w).Encode(IPlist); err != nil {
  log.Fatal("Encoding ip list to json")
  }
}

func postresult(w http.ResponseWriter, r *http.Request) {
var reply Reply

body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
if err != nil {
  panic(err)
  }
if err := r.Body.Close(); err != nil {
  panic(err)
  }
// if we have an unmarshal error we send the error code to the client
if err := json.Unmarshal(body, &reply); err != nil {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(422) // unprocessable entity
  if err := json.NewEncoder(w).Encode(err); err != nil {
    panic(err)
    }
  }
// panic returns immediately!
// so from here we have to tell the client everything was fine
w.Header().Set("Content-Type", "application/json; charset=UTF-8")
w.WriteHeader(http.StatusCreated)

log.Fatal(reply)

}

