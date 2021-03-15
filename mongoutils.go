package main

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func utilGetCurrentEpochTime() int64 {
	now := time.Now()
	secs := now.Unix()
	return secs
}

// connects to MongoDB
func mongoConnect() *mongo.Client {

	mongoDBConnectionString := os.Getenv(mongoDBConnectionURL)
	if mongoDBConnectionString == "" {
		log.Fatal("missing environment variable: ", mongoDBConnectionURL)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(true)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("unable to create a client %v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("unable to initialize connection %v", err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("unable to mongoconnect %v", err)
	}
	return client
}

func mongoFetchRecord(key string) (DBData, error) {
	var filter interface{}
	var myData DBData
	errorMessage := ""
	//debugMessage := ""

	c := mongoConnect()
	ctx := context.Background()
	defer func() {
		_ = c.Disconnect(ctx)
	}()

	filter = bson.D{{"key", key}}

	//filter = bson.D{}

	issueCollection := c.Database(MetadataDB).Collection(MetadataCollection)

	var result bson.M

	err := issueCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		errorMessage = fmt.Sprintf("no record found in the DB")
		log.Printf(errorMessage)
		return myData, errors.New(errorMessage)
	}

	myData = DBData{
		ID:    primitive.ObjectID{},
		Key:   key,
		Value: result,
	}

	fmt.Printf("myData : %v", myData)

	return myData, nil
}

func mongoListRecords() ([]DBData, error) {

	dataList := make([]DBData, 0)

	var filter interface{}

	c := mongoConnect()
	ctx := context.Background()
	defer func() {
		_ = c.Disconnect(ctx)
	}()

	filter = bson.D{}

	issueCollection := c.Database(MetadataDB).Collection(MetadataCollection)
	rs, err := issueCollection.Find(ctx, filter)
	if err != nil {
		return dataList, errors.New(fmt.Sprintf("failed to fetch all issues : %v", err.Error()))
	}

	err = rs.All(ctx, &dataList)
	if err != nil {
		return dataList, errors.New(fmt.Sprintf("failed to fetch all issues : %v", err.Error()))
	}
	if len(dataList) == 0 {
		emptyListOfIssues := make([]DBData, 0)
		return emptyListOfIssues, nil
	}
	/*
		recordTable := make([][]string, 0)

		for _, data := range dataList {
			s, _ := data.ID.MarshalJSON()
			recordTable = append(recordTable, []string{string(s), data.Key, data.Value})
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(fields)

		for _, v := range recordTable {
			table.Append(v)
		}
		table.Render()
	*/
	return dataList, nil
}

func mongoCreateOrUpdateRecord(key string, newValue DBData) (bool, error) {
	c := mongoConnect()
	ctx := context.Background()
	defer func() {
		_ = c.Disconnect(ctx)
	}()

	issueCollection := c.Database(MetadataDB).Collection(MetadataCollection)
	filter := bson.D{{"key", key}}
	update := bson.D{{"$set", newValue}}

	_, err := issueCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return false, errors.New(fmt.Sprintf("failed to update record with key (%v) : %v", key, err.Error()))
	}
	return true, nil
}

func mongoDeleteRecord(key string) (bool, error) {
	errorMessage := ""
	c := mongoConnect()
	ctx := context.Background()
	defer func() {
		_ = c.Disconnect(ctx)
	}()
	isssueCollection := c.Database(MetadataDB).Collection(MetadataCollection)
	filter := bson.D{{"key", key}}
	deletionResult, err := isssueCollection.DeleteOne(ctx, filter)
	if err != nil {
		errorMessage = fmt.Sprintf("failed to delete record : %v", key)
		log.Printf(errorMessage)
		return false, errors.New(errorMessage)
	}
	numberOfRecordsDeleted := deletionResult.DeletedCount
	if numberOfRecordsDeleted == 0 {
		log.Printf("nothing was deleted, because nothing was found")
		return false, nil
	} else {
		log.Printf("number of records deleted : %v", numberOfRecordsDeleted)
		return true, nil
	}

}
