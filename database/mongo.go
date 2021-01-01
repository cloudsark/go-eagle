package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudsark/go-eagle/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Ping struct
type Ping struct {
	Domain    string
	Status    string
	Timestamp time.Time
	Flag      int
}

// Port struct
type Port struct {
	Hostname  string
	Port      string
	Status    string
	Timestamp time.Time
	Flag      int
}

// Ssl struct
type Ssl struct {
	Hostname   string
	RemainDays int
	Flag       int
}

// DiskStats struct
type DiskStats struct {
	HostName string
	Name     string
	Path     string
	FsType   string
	Total    string
	Free     string
	Used     string
	Percent  float64
	Flag     int
}

type Avgload struct {
	Hostname  string
	Loadavg1  float64
	Loadavg5  float64
	Loadavg15 float64
	Flag      int
}

var monDB = os.Getenv("MONGO_DB")
var monUser = os.Getenv("MONGO_USER")
var monPass = os.Getenv("MONGO_PASSWORD")
var monURL = os.Getenv("MONGO_URL")
var monParm = os.Getenv("MONGO_PARM")
var timestamp = time.Now()

func mongodbConn() *mongo.Client {
	URL := fmt.Sprintf("mongodb+srv://%s:%s@%v/%v",
		monUser, monPass, monURL, monParm)
	conn, err := mongo.NewClient(options.Client().ApplyURI(URL))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = conn.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

// InsertDiskStats inserts disk stats into db
func InsertDiskStats(hostname, name, path, fstype, total,
	free, used string, percent float64, flag int) {
	diskStat := DiskStats{hostname, name, path, fstype,
		total, free, used, percent, flag}
	c := mongodbConn()
	collection := c.Database("eagle").Collection("disks")

	insertResult, err := collection.InsertOne(context.TODO(), diskStat)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	logger.GeneralLogger.Println("Inserted disk stats with ID:",
		insertResult.InsertedID)

}

func SortDiskStat(db, coll, hostname,
	diskPath string) map[string]interface{} {
	conn := mongodbConn()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := conn.Database(db).Collection(coll)

	//https://stackoverflow.com/questions/51179588/how-to-sort-and-limit-results-in-mongodb
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	options.SetSkip(0)
	//cursor, err := collection.Find(context.Background(), bson.D{}, options)
	// find by `hostname` field
	cursor, err := collection.Find(context.Background(),
		bson.M{"hostname": hostname, "path": diskPath},
		options)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	var empty map[string]interface{}

	if len(result) != 0 {
		host, _ := result[0]["hostname"].(string)
		path, _ := result[0]["path"]
		percent, _ := result[0]["percent"].(float64)
		flag, _ := result[0]["flag"]
		//fmt.Println(flag)
		//fmt.Println(result[0]["flag"])
		sortMap := map[string]interface{}{
			"hostname": host,
			"path":     path,
			"percent":  percent,
			"flag":     flag,
		}
		return sortMap
	}

	return empty
}

// InsertPing inserts ping result into db
func InsertPing(domain, status string, flag int) {
	p := Ping{domain, status, timestamp, flag}
	c := mongodbConn()
	collection := c.Database("eagle").Collection("ping")

	insertResult, err := collection.InsertOne(context.TODO(), p)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	logger.GeneralLogger.Println("Inserted ping with ID:",
		insertResult.InsertedID)

}

func SortPing(db, coll, hostname string) map[string]interface{} {
	conn := mongodbConn()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := conn.Database(db).Collection(coll)
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	options.SetSkip(0)
	cursor, err := collection.Find(context.Background(),
		bson.M{"domain": hostname},
		options)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	var empty map[string]interface{}

	if len(result) != 0 {
		domain, _ := result[0]["domain"].(string)
		status, _ := result[0]["status"]
		flag, _ := result[0]["flag"]

		sortMap := map[string]interface{}{
			"hostname": domain,
			"status":   status,
			"flag":     flag,
		}
		return sortMap
	}

	return empty
}

// InsertPort inserts port monitoring result into db
func InsertPort(hostname, port, status string, flag int) {
	p := Port{hostname, port, status, timestamp, flag}
	c := mongodbConn()
	collection := c.Database("eagle").Collection("port")

	insertResult, err := collection.InsertOne(context.TODO(), p)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	logger.GeneralLogger.Println("Inserted port with ID:",
		insertResult.InsertedID)
}

func SortPort(db, coll, host, port string) map[string]interface{} {
	conn := mongodbConn()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := conn.Database(db).Collection(coll)
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	options.SetSkip(0)
	cursor, err := collection.Find(context.Background(),
		bson.M{"hostname": host, "port": port},
		options)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	var empty map[string]interface{}

	if len(result) != 0 {
		hostName, _ := result[0]["hostname"].(string)
		port, _ := result[0]["port"]
		status, _ := result[0]["status"]
		flag, _ := result[0]["flag"]

		sortMap := map[string]interface{}{
			"hostname": hostName,
			"port":     port,
			"status":   status,
			"flag":     flag,
		}
		return sortMap
	}

	return empty
}

// InsertSsl inserts certificate monitoring result into db
func InsertSsl(hostname string, days, flag int) {
	p := Ssl{hostname, days, flag}
	c := mongodbConn()
	collection := c.Database("eagle").Collection("ssl")

	insertResult, err := collection.InsertOne(context.TODO(), p)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	logger.GeneralLogger.Println("Inserted ssl with ID:",
		insertResult.InsertedID)
}

func SortSsl(db, coll, host string) map[string]interface{} {
	conn := mongodbConn()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := conn.Database(db).Collection(coll)
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	options.SetSkip(0)
	cursor, err := collection.Find(context.Background(),
		bson.M{"hostname": host},
		options)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	var empty map[string]interface{}

	if len(result) != 0 {
		hostName, _ := result[0]["hostname"].(string)
		days, _ := result[0]["days"]
		flag, _ := result[0]["flag"]

		sortMap := map[string]interface{}{
			"hostname": hostName,
			"days":     days,
			"flag":     flag,
		}
		return sortMap
	}

	return empty
}

// InsertAvgLoad inserts cpu avarage load result into db
func InsertAvgLoad(hostname string, loadavg1, loadavg5,
	loadavg15 float64, flag int) {
	l := Avgload{hostname, loadavg1,
		loadavg5, loadavg15, flag}
	c := mongodbConn()
	collection := c.Database("eagle").Collection("cpu")

	insertResult, err := collection.InsertOne(context.TODO(), l)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	logger.GeneralLogger.Println("Inserted cpu avg load with ID:",
		insertResult.InsertedID)
}

func SortAvgLoad(db, coll, host string) map[string]interface{} {
	conn := mongodbConn()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	collection := conn.Database(db).Collection(coll)
	options := options.Find()
	// Sort by `_id` field descending
	options.SetSort(bson.D{{"_id", -1}})
	options.SetLimit(1)
	options.SetSkip(0)
	cursor, err := collection.Find(context.Background(),
		bson.M{"hostname": host},
		options)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	var result []bson.M
	if err = cursor.All(ctx, &result); err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	var empty map[string]interface{}

	if len(result) != 0 {
		hostName, _ := result[0]["hostname"].(string)
		avg5, _ := result[0]["loadavg5"]
		flag, _ := result[0]["flag"]

		sortMap := map[string]interface{}{
			"hostname": hostName,
			"days":     avg5,
			"flag":     flag,
		}
		return sortMap
	}

	return empty
}
