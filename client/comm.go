package main

import (
    "github.com/niclabs/testResolvers/config"
    "encoding/json"
    "net/http"
    "time"
    "io/ioutil"
    "bytes"
)


type Communication interface {
  Get (cfg config.Configuration) ([]string,error)
  Post (cfg config.Configuration, b *bytes.Buffer) error 
  }

type REST struct {
  url string
  }

func (REST) Get (cfg config.Configuration) ([]string,error) {
  var iplist []string

  url := "http://" + cfg.Server + ":" + cfg.Port  
  cli := http.Client{
    Timeout: time.Second * 10, // Maximum of 2 secs
    }

  req, err := http.NewRequest(http.MethodGet, url + "/get", nil)
  if err != nil {
    return nil, err
    }
  // TODO:  add auth on headers
  req.Header.Set("User-Agent", "testresolver client")
  httpResp, err := cli.Do(req)
  if err != nil {
    return nil, err
    } 

  data, _ := ioutil.ReadAll(httpResp.Body)
  defer httpResp.Body.Close()

  err = json.Unmarshal([]byte(data), &iplist)
  if err != nil {
    return nil, err
    }
  return iplist, nil
  }

func (REST)  Post (cfg config.Configuration, b *bytes.Buffer) error {
  
  url := "http://" + cfg.Server + ":" + cfg.Port  
  _ , err := http.Post(url + "/post", "application/json; charset=utf-8", b)
  return err
  }
