package resolvertests

import (
	"fmt"
	"github.com/miekg/dns"
	"math"
	"math/bits"
	"math/rand"
	"net"
	"strings"
	"time"
)

type Response struct {
	Ip            string
	IsAlive       int
	HasRecursion  int
	HasDNSSEC     int
	HasDNSSECfail int
	QidRatio      int
	PortRatio     int
	Txt           string
}

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

// 1 is random
func testRandomness(shorts []uint16) int {
	var pos, neg, n int = 0, 0, 0
	for _, ui := range shorts {
		pos += bits.OnesCount16(ui)
		neg += 16 - bits.OnesCount16(ui)
		n += 16
	}
	s := math.Abs(float64(pos-neg)) / math.Sqrt(float64(n))
	if math.Erfc(s) < 0.01 {
		return 0
	}
	return 1
}

func Chaosquery(ip string) *dns.Msg {

	c := new(dns.Client)
	m := new(dns.Msg)
	//m.Question[0] = dns.Question{"version.bind.", dns.TypeTXT, dns.ClassCHAOS}

	m.SetQuestion("version.bind.", dns.TypeTXT)
	m.Question[0].Qclass = dns.ClassCHAOS
	msg, _, err := c.Exchange(m, ip+":53")
	if err != nil {
		fmt.Println("error en chaosquery")
	} /*
		for _, ans := range msg.Answer {
			if mx,ok := ans.(*dns.TXT); ok {
				fmt.Println("%s\n", mx.String())
			}
		}*/
	//fmt.Println("%s", msg.Answer)
	return msg
}

func Reverselookup(ip string) []string {

	msg, err := net.LookupAddr(ip)
	if err != nil {
		fmt.Println("error en reverse")
	}
	return msg

}

func checkRandomness(ip string) (int, int) {
	var ids []uint16
	for i := 0; i < 10; i++ {
		line := "bip" + randString(16) + ".niclabs.cl"
		c := new(dns.Client)
		m := new(dns.Msg)
		m.SetQuestion(dns.Fqdn(line), dns.TypeA)
		msg, _, err := c.Exchange(m, ip+":53")
		if err == nil {
			ids = append(ids, msg.MsgHdr.Id)
		}
	}
	return testRandomness(ids), 0
}

func checkAuthority(ip string) string {
	b := strings.Split(ip, ".")
	arp := b[3] + "." + b[2] + "." + b[1] + "." + b[0] + ".in-addr.arpa"
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(arp), dns.TypePTR)
	msg, _, err := c.Exchange(m, ip+":53")

	if err != nil {
		return ""
	}

	if len(msg.Answer) < 1 || len(msg.Ns) < 1 {
		return ""
	}
	return fmt.Sprintf("%s", msg.Ns)
}

func checkDNSSECok(ip string) int {
	line := "sigok.verteiltesysteme.net"
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(line), dns.TypeA)
	m.SetEdns0(4096, true)

	msg, _, err := c.Exchange(m, ip+":53")

	if err != nil {
		return -1 // other error, typically i/o timeout
	}

	return msg.MsgHdr.Rcode
}

func checkDNSSECnook(ip string) int {
	line := "sigfail.verteiltesysteme.net"
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(line), dns.TypeA)
	m.SetEdns0(4096, true)

	msg, _, err := c.Exchange(m, ip+":53")

	if err != nil {
		return -1 // other error
	}
	return msg.MsgHdr.Rcode
}

func CheckDNS(id int, ips <-chan string, results chan<- Response) {
	line := "www.google.com"
	for ip := range ips {
		ip := strings.TrimSpace(ip)
		c := new(dns.Client)
		m := new(dns.Msg)
		r := Response{Ip: ip, IsAlive: 1}

		m.SetQuestion(dns.Fqdn(line), dns.TypeA)
		m.RecursionDesired = true
		m.CheckingDisabled = false
		msg, _, err := c.Exchange(m, ip+":53")

		if err != nil {
			r.IsAlive = 0
		} else {
			if msg != nil {
				if msg.Rcode != dns.RcodeRefused && msg.RecursionAvailable {
					r.HasRecursion = 1
					r.HasDNSSEC = checkDNSSECok(ip)
					r.HasDNSSECfail = checkDNSSECnook(ip)
					r.Txt = checkAuthority(ip)
					//        r.qidRatio,_ = checkRandomness(ip);
				}
			}
		}
		results <- r
	}
}
