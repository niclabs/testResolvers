package config

import (
    "encoding/json"
    "log"
    "os"
)

type Configuration struct {
  User string
  City string
  Pub_key string
  Server string
  Port int
  Home string
  SSHKEYFILE string  // from /
}

func ReadConfig(home string, config *Configuration) int{
file, err := os.Open(home + "config.json") 
if err != nil {  
  log.Fatal("[error] loading config.json")
  return -2
  }  
decoder := json.NewDecoder(file) 
err = decoder.Decode(config) 
if err != nil {  
  log.Fatal("[error] decoding config.json")
  return -3
  }
defer file.Close()
return 0
}

