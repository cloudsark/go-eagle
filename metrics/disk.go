package metrics

import (
	"github.com/cloudsark/go-eagle/alerts"
	"github.com/cloudsark/go-eagle/client"
	"github.com/cloudsark/go-eagle/config"
	c "github.com/cloudsark/go-eagle/constants"
	"github.com/cloudsark/go-eagle/database"
	"github.com/cloudsark/go-eagle/utils"
)

// DiskStatAlert sends disk space alerts
func DiskStatAlert() {
	Hosts := config.Config("Monitor.Metrics")
	mDisks := config.Config("Disks.Monitor")

	for _, host := range Hosts {
		stats := client.GetDiskStat(host)

		for _, s := range stats {
			_, found := utils.Find(mDisks, s.Path)
			if found {

				query := database.SortDiskStat(
					c.OSEnv("MONGO_DB"),
					"disks",
					host,
					s.Path)

				if len(query) == 0 {
					database.InsertDiskStats(host, s.Name, s.Path,
						s.FsType, s.Total, s.Free,
						s.Used, s.Percent, 0)

				}

				if len(query) != 0 {
					percent := s.Percent
					if percent < 60 {
						if query["flag"].(int32) == 0 {
							alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
								c.OSEnv("SLACK_CHANNEL"), host,
								s.Path+" "+alerts.DiskNormal+host,
								"DiskNormal")
						}
						database.InsertDiskStats(host, s.Name, s.Path,
							s.FsType, s.Total, s.Free,
							s.Used, s.Percent, 1)
					}

					if percent > 60 {
						alerts.Alerter(c.OSEnv("SLACK_TOKEN"),
							c.OSEnv("SLACK_CHANNEL"), host,
							s.Path+" "+alerts.DiskCritical+host,
							"DiskCritical")
						database.InsertDiskStats(host, s.Name, s.Path,
							s.FsType, s.Total, s.Free,
							s.Used, s.Percent, 0)
					}

				}

			}
		}
	}
}
