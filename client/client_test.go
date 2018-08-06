package main

import (
  "net"
  "net/http"
  "net/http/httptest"
  "net/url"
  "testing"
  "encoding/json"

  "github.com/niclabs/testResolvers/config"
)

func TestClientNoEncryption(t *testing.T) {
  handler := func (rw http.ResponseWriter, req *http.Request) {

    if req.URL.String() == "/get" {
      rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
      rw.WriteHeader(http.StatusOK)

      err := json.NewEncoder(rw).Encode([]string {"1.2.3.4"})
      if err != nil {
        t.Errorf("Error marshaling list")
      }
    }
  }
  // local http server 
  server := httptest.NewServer(http.HandlerFunc(handler))
  defer server.Close()

  // Now test the client, who receives a "configuration" file

  var cfg config.Configuration = config.Configuration{}
  
  u,err := url.Parse(server.URL)
  cfg.Server, cfg.Port, _ = net.SplitHostPort(u.Host)

  var client Communication = REST {}
  list, err  := client.Get(cfg)

  if err != nil {
    t.Errorf("Error in client.Get %s",err.Error())
  }
  if len(list) != 1 || list[0] != "1.2.3.4" {
    t.Errorf("Message incorrect, got: %s, want: %s.",list,"1.2.3.4")
  }
}
