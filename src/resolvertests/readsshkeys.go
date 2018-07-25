package resolvertests

// private SSH key

import (
    "encoding/json"
    "crypto/dsa"
    "golang.org/x/crypto/ssh"
    "encoding/pem"
    "os"
    "io/ioutil"
    "log"
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

func ReadConfig(home string) int{
file, err := os.Open("config.json") 
if err != nil {  
  log.Fatal("[error] loading config.json")
  return -2
  }  
decoder := json.NewDecoder(file) 
err = decoder.Decode(&config) 
if err != nil {  
  log.Fatal("[error] decoding config.json")
  return -3
  }
defer file.Close()
}

func getPrivateKey() (interface{}, error) {

file, err := os.Open(SSHKEYFILE)
if err != nil {
  log.Fatal("[error] loading private key file")
  return nil,-4
  }

keyBytes, err := ioutil.ReadAll(file)
if err != nil {
  log.Fatal("[error] reading id.dsa")
  return nil,-5
  }
defer file.Close()

return ParseRawPrivateKey(keyBytes)
}
