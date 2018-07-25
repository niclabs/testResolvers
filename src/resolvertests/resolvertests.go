package resolvertests

import (
    "math"
    "math/bits"
    "fmt"
    "github.com/miekg/dns"
    "time"
    "strings"
    "math/rand"
)

type Response struct {
ip string
isAlive int
hasRecursion int
hasDNSSEC int
hasDNSSECfail int
qidRatio int
portRatio int
txt string
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
    letterIdxBits = 6                    // 6 bits to represent a letter index
    letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
    letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandString(n int) string {
b := make([]byte, n)
// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
  if remain == 0 {
    cache, remain = src.Int63(), letterIdxMax
    }
  if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
    b[i] = letterBytes[idx]
    i--
    }
  cache >>= letterIdxBits
  remain--
  }

return string(b)
}

// 1 is random
func testRandomness(shorts []uint16) int {
var pos,neg,n int = 0,0,0
for _, ui := range shorts {
  pos += bits.OnesCount16(ui)
  neg += 16 - bits.OnesCount16(ui)
  n += 16
  }
s := math.Abs(float64(pos - neg)) / math.Sqrt (float64(n))
if math.Erfc(s) < 0.01 {
  return 0
  }
return 1
}

func checkRandomness(ip string) (int,int) {
var ids []uint16
for i:= 0; i < 10; i++ {
  line := "bip"+RandString(16)+".niclabs.cl"
  c := new(dns.Client)
  m := new(dns.Msg)
  m.SetQuestion(dns.Fqdn(line),dns.TypeA)
  msg ,_ , err := c.Exchange(m , ip + ":53")
  if err == nil {
    ids = append(ids,msg.MsgHdr.Id)
    }
  }
fmt.Println(ids)
return testRandomness(ids),0
}

func checkAuthority(ip string) string {
b := strings.Split(ip , ".")
arp := b[3] + "." + b[2] + "." + b[1] + "." + b[0] + ".in-addr.arpa"
c := new(dns.Client)
m := new(dns.Msg)
m.SetQuestion(dns.Fqdn(arp),dns.TypePTR)
msg ,_ ,err := c.Exchange(m,ip + ":53")

if err != nil {
  return ""
  }

if (len(msg.Answer) < 1 || len(msg.Ns) < 1) {
  return ""
  }
return fmt.Sprintf("%s",msg.Ns)
}

func checkDNSSECok (ip string) int {
line := "sigok.verteiltesysteme.net"
c := new(dns.Client)
m := new(dns.Msg)
m.SetQuestion(dns.Fqdn(line),dns.TypeA)
m.SetEdns0(4096,true)

msg ,_ ,err := c.Exchange(m, ip +":53")

if err != nil {
  return -1  // other error, typically i/o timeout
  }

return  msg.MsgHdr.Rcode
}

func checkDNSSECnook (ip string) int {
line := "sigfail.verteiltesysteme.net"
c := new(dns.Client)
m := new(dns.Msg)
m.SetQuestion(dns.Fqdn(line),dns.TypeA)
m.SetEdns0(4096,true)

msg ,_ ,err := c.Exchange(m, ip + ":53")

if err != nil {
  return -1  // other error
  }
return  msg.MsgHdr.Rcode 
}

func CheckDNS(id int, ips <- chan string, results chan <- Response) {
line := "www.hola.com"
for ip :=  range ips {
  ip := strings.TrimSpace(ip)
  c := new(dns.Client)
  m := new(dns.Msg)
  r := Response{ip : ip, isAlive : 1}

  m.SetQuestion(dns.Fqdn(line), dns.TypeA)
  m.RecursionDesired = true
  m.CheckingDisabled = false
  msg ,_ ,err := c.Exchange(m, ip + ":53")

  if err != nil {
    r.isAlive = 0;
  } else {
    if msg != nil {
      if msg.Rcode != dns.RcodeRefused && msg.RecursionAvailable  {
        r.hasRecursion = 1
        r.hasDNSSEC = checkDNSSECok(ip);
        r.hasDNSSECfail = checkDNSSECnook(ip);
        r.txt = checkAuthority(ip);
        r.qidRatio,_ = checkRandomness(ip);
        }
      }
    }
  time.Sleep(1 * time.Second)
  results <- r
  }
}

