package web

import (
	"net"
	"strings"
	"time"

	"github.com/cloudsark/eagle-go/alerts"
	"github.com/cloudsark/eagle-go/config"
	c "github.com/cloudsark/eagle-go/constants"
	"github.com/cloudsark/eagle-go/database"
)

func checkPort(hostname, port string) int8 {
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

func Port() {
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
			check := checkPort(Host, Port)
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
		flag := database.PortQuery(Host, Port)

		if flag == 0 {
			alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
				c.OSEnv("SLACK_CHANNEL"),
				Host, alerts.CheckPort+
					Port+alerts.CheckPortUp+
					Host, "PingUp")
		}
		database.PortInsert(Host, Port, "up", 1)
	}

	for _, hostname := range closed {
		Host := hostname[0]
		Port := hostname[1]
		alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
			c.OSEnv("SLACK_CHANNEL"),
			Host, alerts.CheckPort+
				Port+alerts.CheckPortDown+
				Host, "PingDown")
		database.PortInsert(Host, Port, "down", 0)
	}
}
