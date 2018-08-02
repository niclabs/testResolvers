package server

import (
    "fmt"
    "net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
fmt.Fprintln(w, "This is an API (wow!)")
}

func getfile(w http.ResponseWriter, r *http.Request) {

}

func postresult(w http.ResponseWriter, r *http.Request) {

}

