package main

import (
    "encoding/json"
    "compress/gzip"
    "net/http"
    "time"
    "io/ioutil"
    "bytes"
    "log"

    "github.com/niclabs/testResolvers/config"
    "github.com/niclabs/testResolvers/resolvertests"
)


type Communication interface {
  Get (config.Configuration) ([]string,error)
  Post (config.Configuration, []resolvertests.Response) error 
  }

type REST struct {
  url string
  }

func (r REST) Get (cfg config.Configuration) ([]string,error) {
  var iplist []string
  r.url = "http://" + cfg.Server + ":" + cfg.Port 

  tr := &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: false,
  }
  cli := http.Client{
    Timeout: time.Second * 10, // Maximum of 2 secs
    Transport: tr,
    }

  req, err := http.NewRequest(http.MethodGet, r.url + "/get", nil)
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

func (r REST) Post(cfg config.Configuration,reslice []resolvertests.Response) error {
  r.url ="http://" + cfg.Server + ":" + cfg.Port

  type message struct {
    Time int64 `json:"Time"`
    Login string `json:"Login"`
    Location string `json:"Location"`
    Res [] resolvertests.Response `json:"Res"`
    } 
  m := message {
    time.Now().Unix(),
    cfg.Login,
    cfg.Location,
    reslice,
    }

  b, err := json.Marshal(m)
  if err != nil {
    log.Fatal("Error Marshaling Message " + err.Error())
    return err
    }

  var buf bytes.Buffer
  gzipped := gzip.NewWriter(&buf)
  _ , err = gzipped.Write(b)
  if err != nil {
    log.Fatal("Error Gzipping Message " + err.Error())
    return err
    }
  gzipped.Close()

  tr := &http.Transport{
        MaxIdleConns:       10,
        IdleConnTimeout:    30 * time.Second,
        DisableCompression: true,
  }
  cli := http.Client{
    Timeout: time.Second * 10, // Maximum of 2 secs
    Transport: tr,
    }

  req, err := http.NewRequest("POST", r.url + "/post" , &buf)
  if err != nil {
    log.Fatal("Error POST request " + err.Error())
    return err
    }

  req.Header.Set("User-Agent", "testresolver client")
  req.Header.Set("Content-Type", "application/json")
  req.Header.Set("Content-Encoding", "gzip")

  resp, err := cli.Do(req)
  if err != nil {
    log.Fatal("Error POSTING " + err.Error())
    return err
    }
  defer resp.Body.Close()

  return nil
  }
