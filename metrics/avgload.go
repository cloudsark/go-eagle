package metrics

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/cloudsark/eagle-go/alerts"
	"github.com/cloudsark/eagle-go/config"
	c "github.com/cloudsark/eagle-go/constants"
	"github.com/cloudsark/eagle-go/database"
	"github.com/cloudsark/eagle-go/logger"
)

type Load struct {
	HostName  string  `json:"Hostname"`
	Loadavg1  float64 `json:"Loadavg1"`
	Loadavg5  float64 `json:"Loadavg5"`
	Loadavg15 float64 `json:"Loadavg15"`
}

var clientUsername = os.Getenv("CLIENT_USERNAME")
var clientPassword = os.Getenv("CLIENT_PASSWORD")

func getLoadAvg(ip string) Load {
	Client := &http.Client{}

	req, err := http.NewRequest("GET", "http://"+ip+":10052/api/v1/cpu/load/avg", nil)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	req.Header.Set("Authorization", "Basic "+basicAuth(clientUsername, clientPassword))

	resp, err := Client.Do(req)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}
	var l Load

	err3 := json.Unmarshal(body, &l)
	if err3 != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	return l
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func loadAvgFlag(ip string) (string, float64, int) {
	getavg := getLoadAvg(ip)
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
		getavg := getLoadAvg(host)
		load := getavg.Loadavg5
		sfloat := fmt.Sprintf("%.2f", getavg.Loadavg5)
		flag := database.AvgLoadQuery(host)
		if load < 10 {
			if flag == 0 {
				alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
					c.OSEnv("SLACK_CHANNEL"), host,
					alerts.LoadAvgMsg1+host+alerts.LoadAvgMsg2+sfloat, "AvgLoadNormal")
			}
		}
		database.AvgLoadInsert(host, getavg.Loadavg1, load,
			getavg.Loadavg15, 1)
		if load >= 10 {
			alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
				c.OSEnv("SLACK_CHANNEL"), host,
				alerts.LoadAvgMsg1+host+alerts.LoadAvgMsg2+sfloat, "AvgLoadHigh")
			database.AvgLoadInsert(host, getavg.Loadavg1, load,
				getavg.Loadavg15, 0)
		}

	}
}
