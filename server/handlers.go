package main

import (
    "fmt"
    "encoding/json"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "This is an API (wow!)")
}

func getfile(w http.ResponseWriter, r *http.Request) {

w.Header().Set("Content-Type", "application/json; charset=UTF-8")
w.WriteHeader(http.StatusOK)
if err := json.NewEncoder(w).Encode(IPlist); err != nil {
  panic(err)
  }
}

func postresult(w http.ResponseWriter, r *http.Request) {

}

