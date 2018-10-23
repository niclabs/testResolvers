package dnsserver

import (
  "time"
  "sync"
  "strconv"
  "strings"
  "log"
  "math/rand"
  "os"
  "github.com/miekg/dns"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
  letterIdxBits = 6                    // 6 bits to represent a letter index
  letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
  letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randString(n int) string {
var src = rand.NewSource(time.Now().UnixNano())
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

type qidport struct {
  qid string 
  port string
}

var clients map[string][]qidport
var clientid map[string]int
var lock = sync.RWMutex{}
var zone string 


// invalid query generates SERVFAIL
func nxDomain(msg *dns.Msg, domain string, address string) {
  (*msg).Answer = append (msg.Answer,
  &dns.TXT{
    Hdr: dns.RR_Header{
      Name:domain,
      Rrtype:dns.RcodeNameError,
      Class: dns.ClassINET,
      Ttl: 0,
    },
  })
  log.Printf("domain: %s Qtype: %d from: %s\n", domain,msg.Question[0].Qtype, address)
}

func getidport(port string,msg dns.Msg) qidport {
  return qidport {
    qid : strconv.Itoa(int(msg.MsgHdr.Id)),
    port : port,
  }
}

func packidport (k string) string {
  s := ""
  for _, client := range clients[k] {
    s = s + client.qid + "," + client.port + ","
  }
  return s
}

type handler struct{}
func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
  msg := dns.Msg{}
  msg.SetReply(r)
  msg.Authoritative = true
  domain := msg.Question[0].Name
  address := w.RemoteAddr().String()

  defer w.WriteMsg(&msg)

  /*
  Checking for invalid queries:
    is the query domain smaller than the zone? weird
    is the last part of the query domain not my zone, weird again
    is somebody asking for a keymaped domain? nope
  */

  if len(domain) < len (zone) {nxDomain(&msg,domain,address);return}
  
  lendiff := len(domain) - len (zone)
  if domain[lendiff:] != zone {nxDomain(&msg,domain,address);return}
  if lendiff > 0 && domain[lendiff-1:lendiff] != "." {
    nxDomain(&msg,domain,address)
    return
  }
  s := strings.Split(domain,".")
  key := s[1]
  if _,ok := clientid[key]; !ok && zone != domain {
    nxDomain(&msg,domain,address)
    return
  }

  switch r.Question[0].Qtype {
  case dns.TypeTXT:
    if (domain == zone) {
      randstring := randString(16)
      cname := "0." + randstring + "." + zone
      clientid[randstring] = 0
      msg.Answer = append (msg.Answer, 
      &dns.CNAME{
        Hdr: dns.RR_Header{
          Name:domain, 
          Rrtype:dns.TypeCNAME, 
          Class: dns.ClassINET, 
          Ttl: 0,
        },
        Target: cname,
      })
    } else {
      if (clientid[key] < 10) {
        add := strings.Split(address,":")
        clients[key] = append(clients[key],getidport(add[1],msg))
        clientid[key] = clientid[key] + 1
        cname := strconv.Itoa(clientid[key]) + "." + key + "." + zone
        msg.Answer = append (msg.Answer,
        &dns.CNAME{
          Hdr: dns.RR_Header{
            Name:domain,
            Rrtype:dns.TypeCNAME,
            Class: dns.ClassINET,
            Ttl: 0,
          },
          Target: cname,
        })
      } else {
        msg.Answer = append (msg.Answer, 
          &dns.TXT{
            Hdr: dns.RR_Header{
              Name:domain, 
              Rrtype:dns.TypeTXT, 
              Class: dns.ClassINET, 
              Ttl: 0,
            },
            Txt: []string{packidport(key)},
          })
        delete(clientid,key)
        delete(clients,key)
      }  
    }
  return
  }
}

func main() {
  zone = os.Args[1]
  clients = make(map[string][]qidport)
  clientid = make(map[string]int)
  ip := os.Args[2]

  srv := &dns.Server{Addr: ip + ":" + strconv.Itoa(53), Net: "udp"}
  log.Printf("Starting at %s:%d\n", ip, 53)
  log.Println(zone)

  srv.Handler = &handler{}
  if err := srv.ListenAndServe(); err != nil {
    log.Fatalf("Failed to set udp listener %s\n", err.Error())
  }
}
