package web

import (
	"net"
	"time"

	c "github.com/cloudsark/eagle-go/constants"

	"github.com/cloudsark/eagle-go/alerts"
	"github.com/cloudsark/eagle-go/config"
	"github.com/cloudsark/eagle-go/database"
	"github.com/cloudsark/eagle-go/logger"
	"github.com/tatsushid/go-fastping"
)

func pingDomain(domain string) string {
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

// Ping ping given domain and returns an ip or error
func Ping() {
	mPing := make(map[string]string)
	var unresolved []string
	var resolved []string

	pingSlice := config.Config("Monitor.Ping")
	for _, d := range pingSlice {
		mPing[d] = pingDomain(d)
	}

	// append unresoved domains to unresolved slice
	// append resolved domains to resolved slice
	for key, value := range mPing {
		if value == "" {
			unresolved = append(unresolved, key)
		} else {
			resolved = append(resolved, key)
		}
	}

	for _, domain := range resolved {
		flag, _ := database.PingQuery(domain) // it can also return timestamp
		//ToDo: calculate for how long the domain was down
		if flag == 0 {
			alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
				c.OSEnv("SLACK_CHANNEL"), domain,
				domain+alerts.PingUp, "PingUp")
		}
		database.PingInsert(domain, "up", 1)
	}

	for _, domain := range unresolved {
		alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
			c.OSEnv("SLACK_CHANNEL"), domain,
			domain+alerts.PingDown, "PingDown")
		database.PingInsert(domain, "down", 0)
	}

}
