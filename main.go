package main

import (
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mylog "log"
	"log/syslog"
	"os"
)

var (
	syslogger            *syslog.Writer
	MetadataDB           = ""
	MetadataCollection   = ""
	mongoDBConnectionURL = ""
)

type DBData struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Key   string             `bson:"key"`
	Value interface{}        `bson:"value"`
}

var (
//fields = []string{"ID", "Key", "Value"}
)

func initialize() {
	log.SetFormatter(&log.JSONFormatter{})

	MetadataDB = "db1"
	MetadataCollection = "collection1"
	mongoDBConnectionURL = "MONGOURL"

	syslogger, _ = syslog.New(syslog.LOG_DEBUG, "mongoAPI")
	mylog.SetOutput(syslogger)

	mongoDBConnectionString := os.Getenv(mongoDBConnectionURL)
	if mongoDBConnectionString == "" {
		log.Fatal("missing environment variable: ", mongoDBConnectionURL)
	}
	setDBandCollection(MetadataDB, MetadataCollection)
}

func setDBandCollection(dbName string, collectionName string) {
	MetadataDB = dbName
	MetadataCollection = collectionName
}

func main() {
	initialize()
	mylog.Println("mongoAPI : initalized")

}
