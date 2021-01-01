package config

import (
	"strings"

	"github.com/cloudsark/go-eagle/logger"
	"github.com/spf13/viper"
)

type Configurations struct {
	Monitor   Monitor   `yaml:"Monitor"`
	Alerts    Alerts    `yaml:"Alerts"`
	Intervals Intervals `yaml:"Triggers"`
}
type Monitor struct {
	SSL     []string `yaml:"SSL"`
	Ping    []string `yaml:"Ping"`
	Port    []string `yaml:"Port"`
	Metrics []string `yaml:"Metrics"`
}
type Alerts struct {
	Slack    string `yaml:"Slack"`
	Telegram string `yaml:"Telegram"`
	Email    string `yaml:"Email"`
}
type Intervals struct {
	Ssl  string `yaml:"Ssl"`
	Ping string `yaml:"Ping"`
	Port string `yaml:"Port"`
}

// Config returns slice from struct
func Config(slice string) []string {
	// Set the file name of the configurations file
	viper.SetConfigName("main")

	// Set the path to look for the configurations file
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	viper.SetConfigType("yml")
	var c Configurations

	if err := viper.ReadInConfig(); err != nil {
		logger.ErrorLogger.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		logger.ErrorLogger.Fatalf("Unable to decode into struct, %v", err)
	}

	return viper.GetStringSlice(slice)
}

// Cron returns interval from Intervals
func Cron(interval string) string {
	Interval := Config(interval)
	cron := strings.Join(Interval, " ")
	return cron
}

// AlertStruct returns the status of alerts integration
// Output: true|false
func AlertStruct(integ string) string {
	alerts := Config("Alerts." + integ)
	return alerts[0]
}
