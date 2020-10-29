package eagle

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/cloudsark/eagle-go/logger"

	_ "github.com/go-sql-driver/mysql"
)

// Ping struct for ping table (mysql)
type Ping struct {
	ID        int
	Domain    string
	Status    string
	Timestamp string
	Flag      int
}

// Certificate struct for certificate (ssl) table (mysql)
type Certificate struct {
	ID         int
	Hostname   string
	RemainDays int
	Flag       int
}

// Port struct for ports table (mysql)
type Port struct {
	ID        int
	Hostname  string
	Port      string
	Status    string
	Timestamp string
	Flag      int
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")"+"/"+dbName)
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}
	return db
}

// PingInsert inserts ping results in database
func PingInsert(domain string, status string, flag int) {
	ts := time.Now()
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO Ping(domain, status, ts, flag) VALUES (?,?,?,?)")
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	insert, err := stmt.Exec(domain, status, ts, flag)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	id, err := insert.LastInsertId()
	ID := strconv.FormatInt(id, 10)
	logger.GeneralLogger.Println("Insert record " + ID + " to ping")

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	defer db.Close()
}

// PingQuery returns the last known flag and timestamp for a spesific domain from mysql
/*
Input: domain name
Outputs:
  1. Flag  0|1
  2. Timestamp in "2020-09-11 00:56:26" format
*/
func PingQuery(domain string) (int, string) {
	var p Ping
	db := dbConn()
	query, err := db.Query("SELECT * FROM Ping where domain =? ORDER BY id DESC LIMIT 1;", domain)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	for query.Next() {
		err = query.Scan(&p.ID, &p.Domain, &p.Status, &p.Timestamp, &p.Flag)
		if err != nil {
			logger.ErrorLogger.Fatalln(err.Error())
		}
	}
	return p.Flag, p.Timestamp
}

// SslInsert inserts ssl results in database
func SslInsert(hostname string, RemainingDays int, flag int) {
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO certificate(hostname, remaining_days, flag) VALUES (?,?,?)")
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	insert, err := stmt.Exec(hostname, RemainingDays, flag)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	id, err := insert.LastInsertId()
	ID := strconv.FormatInt(id, 10)
	logger.GeneralLogger.Println("Insert record " + ID + " to certificate")

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	defer db.Close()
}

// SslQuery returns the last known flag for a spesific hostname from mysql certificate table
/*
Input: domain name
Outputs:
  1. Flag  0|1
*/
func SslQuery(hostname string) int {
	var c Certificate
	db := dbConn()
	query, err := db.Query("SELECT * FROM certificate where hostname =? ORDER BY id DESC LIMIT 1;", hostname)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	for query.Next() {
		err = query.Scan(&c.ID, &c.Hostname, &c.RemainDays, &c.Flag)
		if err != nil {
			logger.ErrorLogger.Fatalln(err.Error())
		}
	}
	return c.Flag
}

// PortInsert inserts ping results in database
func PortInsert(hostname string, port string, status string, flag int) {
	ts := time.Now()
	db := dbConn()
	stmt, err := db.Prepare("INSERT INTO ports(hostname, port, status, ts, flag) VALUES (?,?,?,?,?)")
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	insert, err := stmt.Exec(hostname, port, status, ts, flag)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	id, err := insert.LastInsertId()
	ID := strconv.FormatInt(id, 10)
	logger.GeneralLogger.Println("Insert record " + ID + " to port")
	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	defer db.Close()
}

// PortQuery returns the last known flag for a spesific hostname from mysql certificate table
/*
Input:  hostname
Outputs:
  1. Flag  0|1
*/
func PortQuery(hostname, port string) int {
	var p Port
	db := dbConn()
	query, err := db.Query("SELECT * FROM ports where hostname =? AND port=? ORDER BY id DESC LIMIT 1;", hostname, port)

	if err != nil {
		logger.ErrorLogger.Fatalln(err.Error())
	}

	for query.Next() {
		err = query.Scan(&p.ID, &p.Hostname, &p.Port, &p.Status, &p.Timestamp, &p.Flag)
		if err != nil {
			logger.ErrorLogger.Fatalln(err.Error())
		}
	}
	return p.Flag
}
