package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/cloudsark/go-eagle/logger"
)

const (
	HTTPTimeOut    = time.Second * 60
	errGetHTTP     = "HTTP GET failed"
	clientProto    = "http://"
	clientPort     = ":10052"
	pathAvgCPULoad = "api/v1/cpu/load/avg"
	pathDiskStat   = "api/v1/disk/usage/stat"
)

var clientUsername = os.Getenv("CLIENT_USERNAME")
var clientPassword = os.Getenv("CLIENT_PASSWORD")

type CpuAvgLoad struct {
	HostName  string  `json:"Hostname"`
	Loadavg1  float64 `json:"Loadavg1"`
	Loadavg5  float64 `json:"Loadavg5"`
	Loadavg15 float64 `json:"Loadavg15"`
}

type DiskStat struct {
	HostName string  `json:"Hostname"`
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	FsType   string  `json:"fstype"`
	Total    string  `json:"total"`
	Free     string  `json:"free"`
	Used     string  `json:"used"`
	Percent  float64 `json:"percent"`
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getHTTP(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+basicAuth(clientUsername,
		clientPassword))

	client := http.Client{Timeout: HTTPTimeOut}

	resp, err := client.Do(req)
	if err != nil {
		logger.ErrorLogger.Printf("%v: %v", errGetHTTP, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.ErrorLogger.Printf("%v: %v", errGetHTTP, err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.ErrorLogger.Printf("%v: unexpected HTTP status code %v",
			errGetHTTP, resp.StatusCode)
	}

	return body, nil
}

// GetCPULoadAvg returns cpu avarage load for a specific hostname
/*
Inputs: apiURL
  example: example.com
Output: json
  example: {server.example.com 2.0 5.0 6.0}
*/
func GetCPULoadAvg(apiURL string) CpuAvgLoad {
	URL := fmt.Sprintf("%s/%s", clientProto+apiURL+clientPort, pathAvgCPULoad)

	body, err := getHTTP(URL)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	var avgLoad CpuAvgLoad

	err = json.Unmarshal(body, &avgLoad)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	return avgLoad
}

// GetDiskStat returns array of Disk Stats for a specific hostname
/*
Inputs: apiURL
  example: example.com
Output: json
  example: [{development dm-0 / xfs 37798.00 30751.00 7046.00 18.64}]
*/
func GetDiskStat(apiURL string) []DiskStat {
	URL := fmt.Sprintf("%s/%s", clientProto+apiURL+clientPort,
		pathDiskStat)
	body, err := getHTTP(URL)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	var diskStat []DiskStat

	err = json.Unmarshal(body, &diskStat)
	if err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

	return diskStat
}
