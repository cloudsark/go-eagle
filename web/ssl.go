package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cloudsark/eagle-go/alerts"
	"github.com/cloudsark/eagle-go/config"
	c "github.com/cloudsark/eagle-go/constants"
	"github.com/cloudsark/eagle-go/database"
	"github.com/cloudsark/eagle-go/logger"
)

func sslVerify(hostname string) int {
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

func Ssl() {
	mSsl := make(map[string]int)
	sslSlice := config.Config("Monitor.SSL")

	for _, d := range sslSlice {
		mSsl[d] = sslVerify(d)
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
				flag := database.SslQuery(hName)
				if flag == 0 {
					alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
						c.OSEnv("SLACK_CHANNEL"),
						hostname, hostname+
							alerts.ValidSsl, "SslValid")
				}
			}
			database.SslInsert(hName, rDays, 1)
		case rDays < 30 && rDays > 20:
			for _, hostname := range warn {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpiredDate2,
					"SslNotValidWarn")
			}
			database.SslInsert(hName, rDays, 0)
		case rDays <= 20:
			for _, hostname := range crit {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpiredDate2,
					"SslNotValidCrit")
			}
			database.SslInsert(hName, rDays, 0)
		case rDays <= 0:
			critical = append(crit, hName)
			for _, hostname := range critical {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), hostname,
					hostname+alerts.SslExpiredDate1+
						fmt.Sprintf("%d", rDays)+
						alerts.SslExpired,
					"SslNotValidCrit")
			}
			database.SslInsert(hName, rDays, 0)
		}
	}
}
