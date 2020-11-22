package main

import (
	"github.com/cloudsark/eagle-go/config"
	"github.com/cloudsark/eagle-go/metrics"
	"github.com/cloudsark/eagle-go/web"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New()
	c.AddFunc(config.Cron("Intervals.Ping"),
		func() { web.Ping() })
	c.AddFunc(config.Cron("Intervals.Port"),
		func() { web.Port() })
	c.AddFunc(config.Cron("Intervals.Ssl"),
		func() { web.Ssl() })
	c.AddFunc(config.Cron("Intervals.Metrics"),
		func() { metrics.LoadAvgAlert() })
	c.Start()
	select {}
}
