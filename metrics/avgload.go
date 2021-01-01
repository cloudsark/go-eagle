package metrics

import (
	"fmt"

	"github.com/cloudsark/go-eagle/alerts"
	"github.com/cloudsark/go-eagle/client"
	"github.com/cloudsark/go-eagle/config"
	c "github.com/cloudsark/go-eagle/constants"
	"github.com/cloudsark/go-eagle/database"
)

func loadAvgFlag(ip string) (string, float64, int) {
	getavg := client.GetCPULoadAvg(ip)
	var flag int
	if getavg.Loadavg5 >= 10 {
		flag = 1
	}
	if getavg.Loadavg5 <= 10 {
		flag = 0
	}
	return getavg.HostName, getavg.Loadavg5, flag
}

// LoadAvgAlert sends avg cpu alerts
func LoadAvgAlert() {
	loadAvgSlice := config.Config("Monitor.Metrics")

	for _, host := range loadAvgSlice {
		getavg := client.GetCPULoadAvg(host)

		load := getavg.Loadavg5
		sfloat := fmt.Sprintf("%.2f", getavg.Loadavg5)
		query := database.SortAvgLoad(c.OSEnv("MONGO_DB"),
			"cpu", host)
		if len(query) == 0 {
			database.InsertAvgLoad(host, getavg.Loadavg1, load,
				getavg.Loadavg15, 1)
		}
		if len(query) != 0 {
			flag := query["flag"].(int32)
			if load < 10 {
				if flag == 0 {
					alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
						c.OSEnv("SLACK_CHANNEL"), host,
						alerts.LoadAvgMsg1+host+alerts.LoadAvgMsg2+sfloat, "AvgLoadNormal")
				}
			}
			database.InsertAvgLoad(host, getavg.Loadavg1, load,
				getavg.Loadavg15, 1)
			if load >= 10 {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), host,
					alerts.LoadAvgMsg1+host+alerts.LoadAvgMsg2+sfloat, "AvgLoadHigh")
				database.InsertAvgLoad(host, getavg.Loadavg1, load,
					getavg.Loadavg15, 0)
			}

		}
	}
}
