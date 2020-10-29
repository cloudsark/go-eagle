package eagle

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/cloudsark/eagle-go/logger"
	"github.com/tatsushid/go-fastping"
)

// PingDomain ping given domain and returns an ip or error
func PingDomain(domain string) string {
	var st string
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", domain)
	if err != nil {
		return ""
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		st = addr.String()
	}

	err = p.Run()
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
	return st
}

// VerifySsl checks the expired domain date
func VerifySsl(hostname string) int {
	resp, err := http.Head(hostname)
	if err != nil {
		logger.ErrorLogger.Fatalf("Unable to get %q: %s\n", hostname, err)
	}
	resp.Body.Close()
	if resp.TLS == nil {
		logger.ErrorLogger.Fatalf("%q is not HTTPS\n", hostname)
	}

	var days []float64
	for _, cert := range resp.TLS.PeerCertificates {
		for _, name := range cert.DNSNames {
			if !strings.Contains(hostname, name) {
				continue
			}
			dur := cert.NotAfter.Sub(time.Now())
			days = append(days, dur.Hours()/24)
		}
	}
	return int(days[0])
}

// CheckPort checks open ports
func CheckPort(hostname, port string) int8 {
	timeout := time.Second
	var connection int8
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(hostname, port), timeout)

	if err != nil {
		connection = 0
	}

	if conn != nil {
		defer conn.Close()
		connection = 1
	}
	return connection
}
