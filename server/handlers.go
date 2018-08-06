package main

import (
    "fmt"
    "encoding/json"
    "compress/gzip"
    "net/http"
    "log"

    "github.com/niclabs/testResolvers/resolvertests"
)

type Query struct {
  ips []string
  }

type Reply struct {
  Time int64
  Login string
  Location string
  Res [] resolvertests.Response
  }

func getData(r *http.Request) (*Reply,error) {
  var data Reply
  var decoder *json.Decoder
  switch r.Header.Get("Content-Encoding") {
    case "gzip":
      gz, err := gzip.NewReader(r.Body)
      if err != nil {
        return nil, err
        }
      defer gz.Close()
      decoder = json.NewDecoder(gz)
    default:
      decoder = json.NewDecoder(r.Body)
    }
  err := decoder.Decode(&data)
  if err != nil {
    log.Fatal(err.Error())
    return nil, err
    }
  return &data, err
  }

func index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintln(w, "This is an API (wow!)")
  }

func getfile(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  err := json.NewEncoder(w).Encode(IPlist)
  if err != nil {
    log.Fatal("Encoding ip list to json")
    }
  }

func postresult(w http.ResponseWriter, r *http.Request) {

  reply, err := getData(r)
  if err != nil {
    log.Fatal(err.Error())
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(422) // unprocessable entity
    if err := json.NewEncoder(w).Encode(err); err != nil {
      http.Error(w, err.Error(), 500)
      return
      }
    }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)

  log.Println(*reply)
  }

