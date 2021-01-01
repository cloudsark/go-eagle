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

// DisStats struct
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

var monDB = os.Getenv("MONGO_DB")
var monUser = os.Getenv("MONGO_USER")
var monPass = os.Getenv("MONGO_PASSWORD")
var monURL = os.Getenv("MONGO_URL")
var monParm = os.Getenv("MONGO_PARM")

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
	//defer cli.Disconnect(ctx)
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
