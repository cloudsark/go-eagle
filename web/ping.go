package web

import (
	"net"
	"time"

	c "github.com/cloudsark/go-eagle/constants"

	"github.com/cloudsark/go-eagle/alerts"
	"github.com/cloudsark/go-eagle/config"
	"github.com/cloudsark/go-eagle/database"
	"github.com/cloudsark/go-eagle/logger"
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

	// append unresolved domains to unresolved slice
	// append resolved domains to resolved slice
	for key, value := range mPing {
		if value == "" {
			unresolved = append(unresolved, key)
		} else {
			resolved = append(resolved, key)
		}
	}

	for _, domain := range resolved {
		query := database.SortPing(c.OSEnv("MONGO_DB"),
			"ping", domain)
		if len(query) == 0 {
			database.InsertPing(domain, "up", 1)
		}
		if len(query) != 0 {
			flag := query["flag"].(int32)
			//ToDo: calculate for how long the domain was down
			if flag == 0 {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), domain,
					domain+alerts.PingUp, "PingUp")
			}
			database.InsertPing(domain, "up", 1)
		}
	}

	for _, domain := range unresolved {
		alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
			c.OSEnv("SLACK_CHANNEL"), domain,
			domain+alerts.PingDown, "PingDown")
		database.InsertPing(domain, "down", 0)
	}
}
