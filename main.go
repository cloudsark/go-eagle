package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/cloudsark/eagle-go/alerts"
	"github.com/cloudsark/eagle-go/config"
	eagle "github.com/cloudsark/eagle-go/eagle"
	"github.com/robfig/cron/v3"
)

var slackToken = os.Getenv("SLACK_TOKEN")
var slackChannel = os.Getenv("SLACK_CHANNEL")

func main() {
	c := cron.New()
	c.AddFunc(config.Cron("Intervals.Ping"),
		func() { ping() })
	c.AddFunc(config.Cron("Intervals.Port"),
		func() { checkPort() })
	c.AddFunc(config.Cron("Intervals.Ssl"),
		func() { ssl() })
	c.Start()
	select {}
}

func ping() {
	mPing := make(map[string]string)
	var unresolved []string
	var resolved []string

	pingSlice := config.Config("Monitor.Ping")
	for _, d := range pingSlice {
		mPing[d] = eagle.PingDomain(d)
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
		flag, _ := eagle.PingQuery(domain) // it can also return timestamp
		//ToDo: calculate for how long the domain was down
		if flag == 0 {
			alerts.Alerter(slackToken, slackChannel,
				domain, domain+alerts.PingUp, "PingUp")
		}
		eagle.PingInsert(domain, "up", 1)
	}

	for _, domain := range unresolved {
		alerts.Alerter(slackToken, slackChannel,
			domain, domain+alerts.PingDown, "PingDown")
		eagle.PingInsert(domain, "down", 0)
	}

}

func ssl() {
	mSsl := make(map[string]int)
	sslSlice := config.Config("Monitor.SSL")

	for _, d := range sslSlice {
		mSsl[d] = eagle.VerifySsl(d)
	}

	var normal []string
	var warning []string
	var critical []string

	for hName, rDays := range mSsl {
		norm := append(normal, hName)
		warn := append(warning, hName)
		crit := append(critical, hName)
		switch {
		case rDays >= 30:
			for _, hostname := range norm {
				flag := eagle.SslQuery(hName)
				if flag == 0 {
					alerts.Alerter(slackToken, slackChannel,
						hostname, hostname+
							alerts.ValidSsl, "SslValid")
				}
			}
			eagle.SslInsert(hName, rDays, 1)
		case rDays < 30 && rDays > 20:
			for _, hostname := range warn {
				alerts.Alerter(slackToken, slackChannel, hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpiredDate2,
					"SslNotValidWarn")
			}
			eagle.SslInsert(hName, rDays, 0)
		case rDays <= 20:
			for _, hostname := range crit {
				alerts.Alerter(slackToken, slackChannel, hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpiredDate2,
					"SslNotValidCrit")
			}
			eagle.SslInsert(hName, rDays, 0)
		case rDays <= 0:
			critical = append(crit, hName)
			for _, hostname := range critical {
				alerts.Alerter(slackToken, slackChannel, hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpired,
					"SslNotValidCrit")
			}
			eagle.SslInsert(hName, rDays, 0)
		}
	}
}

func checkPort() {
	portSlice := config.Config("Monitor.Port")
	opened := [][]string{}
	closed := [][]string{}

	for _, host := range portSlice {
		full := strings.Split(host, ":")
		Host := full[0]
		Port := full[1]

		S := [][]string{
			[]string{Host, Port},
		}

		for i := 0; i < len(S); i++ {
			check := eagle.CheckPort(Host, Port)
			if check == 1 {
				opened = append(opened, S[i])
			} else if check == 0 {
				closed = append(closed, S[i])
			}

		}
	}

	for _, hostname := range opened {
		Host := hostname[0]
		Port := hostname[1]
		flag := eagle.PortQuery(Host, Port)

		if flag == 0 {
			alerts.Alerter(slackToken, slackChannel,
				Host, alerts.CheckPort+
					Port+alerts.CheckPortUp+
					Host, "PingUp")
		}
		eagle.PortInsert(Host, Port, "up", 1)
	}

	for _, hostname := range closed {
		Host := hostname[0]
		Port := hostname[1]
		alerts.Alerter(slackToken, slackChannel,
			Host, alerts.CheckPort+
				Port+alerts.CheckPortDown+
				Host, "PingDown")
		eagle.PortInsert(Host, Port, "down", 0)
	}
}
