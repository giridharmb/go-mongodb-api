package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"testing"
)

/*
# go test -v

=== RUN   TestInsertRecord
2021/03/14 19:01:28 TestInsertRecord...
--- PASS: TestInsertRecord (2.88s)
=== RUN   TestFetchRecord
--- PASS: TestFetchRecord (2.46s)
=== RUN   TestListRecords
--- PASS: TestListRecords (3.72s)
=== RUN   TestDeleteRecords
--- PASS: TestDeleteRecords (2.47s)
PASS
ok  	mongoDBAPI	11.834s
*/

/*
# go test -v -bench=BenchmarkCrud

=== RUN   TestInsertRecord
2021/03/14 19:08:55 TestInsertRecord...
--- PASS: TestInsertRecord (2.55s)
=== RUN   TestFetchRecord
--- PASS: TestFetchRecord (2.35s)
=== RUN   TestListRecords
--- PASS: TestListRecords (3.59s)
=== RUN   TestDeleteRecords
--- PASS: TestDeleteRecords (2.38s)
goos: darwin
goarch: amd64
pkg: mongoDBAPI
BenchmarkCrud
BenchmarkCrud-16    	       1	4958402574 ns/op
PASS
ok  	mongoDBAPI	16.104s
*/

func TestInsertRecord(t *testing.T) {

	log.Printf("TestInsertRecord...")

	initialize()

	setDBandCollection("testingDB", "testCollection")

	testData1 := make(map[string]string)
	testData1["animal"] = "cow"
	testData1["sound"] = "mooooo"

	key := "key1"
	myData := DBData{
		ID:    primitive.ObjectID{},
		Key:   key,
		Value: testData1,
	}

	status, err := mongoCreateOrUpdateRecord(key, myData)
	if err != nil {
		t.Errorf("could not create/update record : %v", err.Error())
		return
	}
	log.Printf("status : %v", status)

	record, err := mongoFetchRecord(key)
	if err != nil {
		t.Errorf("could not fetch record with key (%v)  %v", key, err.Error())
		return
	}
	if record.Key != key {
		t.Errorf("could not fetch record with key (%v)", key)
	}
}

func TestFetchRecord(t *testing.T) {

	log.Printf("TestFetchRecord...")

	initialize()

	setDBandCollection("testingDB", "testCollection")

	testData1 := make(map[string]string)
	testData1["animal"] = "dog"
	testData1["sound"] = "bow"

	key := "key2"

	myData := DBData{
		ID:    primitive.ObjectID{},
		Key:   key,
		Value: testData1,
	}

	status, err := mongoCreateOrUpdateRecord(key, myData)
	if err != nil {
		t.Errorf("could not create/update record : %v", err.Error())
		return
	}
	log.Printf("status : %v", status)

	record, err := mongoFetchRecord(key)
	if err != nil {
		t.Errorf("could not fetch record with key (%v)  %v", key, err.Error())
		return
	}
	if record.Key != key {
		t.Errorf("could not fetch record with key (%v)", key)
	}
}

func TestListRecords(t *testing.T) {

	log.Printf("TestListRecords...")

	initialize()

	setDBandCollection("testingDB", "testCollection2")

	testData1 := make(map[string]string)
	testData1["animal"] = "cow"
	testData1["sound"] = "mooo"

	testData2 := make(map[string]string)
	testData2["animal"] = "dog"
	testData2["sound"] = "bow"

	key1 := "record1"

	myData1 := DBData{
		ID:    primitive.ObjectID{},
		Key:   key1,
		Value: testData1,
	}

	key2 := "record2"

	myData2 := DBData{
		ID:    primitive.ObjectID{},
		Key:   key2,
		Value: testData2,
	}

	//////////////////////////////////////////////////////////////////////

	status1, err := mongoCreateOrUpdateRecord(key1, myData1)
	if err != nil {
		t.Errorf("could not create/update record : %v", err.Error())
		return
	}
	log.Printf("status : %v", status1)

	//////////////////////////////////////////////////////////////////////

	status2, err := mongoCreateOrUpdateRecord(key2, myData2)
	if err != nil {
		t.Errorf("could not create/update record : %v", err.Error())
		return
	}
	log.Printf("status : %v", status2)

	//////////////////////////////////////////////////////////////////////

	records, err := mongoListRecords()
	if err != nil {
		t.Errorf("could not fetch records : %v", err.Error())
		return
	}
	log.Printf("found (%v) records", len(records))

	record1found := false
	record2found := false

	for _, record := range records {
		if record.Key == key1 {
			record1found = true
		}
		if record.Key == key2 {
			record2found = true
		}
	}

	if record1found == false {
		t.Errorf("could not find record with key : %v", key1)
	}

	if record2found == false {
		t.Errorf("could not find record with key : %v", key2)
	}

}

func TestDeleteRecords(t *testing.T) {

	log.Printf("TestDeleteRecords...")

	initialize()

	setDBandCollection("testingDB", "testCollection2")

	testData1 := make(map[string]string)
	testData1["animal"] = "crow"
	testData1["sound"] = "kookoo"

	key1 := "record1"

	myData1 := DBData{
		ID:    primitive.ObjectID{},
		Key:   key1,
		Value: testData1,
	}

	//////////////////////////////////////////////////////////////////////

	status1, err := mongoCreateOrUpdateRecord(key1, myData1)
	if err != nil {
		t.Errorf("could not create/update record : %v", err.Error())
		return
	}
	log.Printf("status : %v", status1)

	//////////////////////////////////////////////////////////////////////

	status, err := mongoDeleteRecord(key1)
	if err != nil {
		t.Errorf("could not delete record : %v", err.Error())
	}
	if status == false {
		t.Errorf("could not delete record with key : %v", key1)
	}
}

func BenchmarkCrud(b *testing.B) {
	initialize()
	setDBandCollection("testingDB", "testCollection3")

	testData1 := make(map[string]string)
	testData1["animal"] = "dog"
	testData1["sound"] = "bow"

	key := "key2"

	myData := DBData{
		ID:    primitive.ObjectID{},
		Key:   key,
		Value: testData1,
	}

	for i := 0; i < b.N; i++ {

		status, err := mongoCreateOrUpdateRecord(key, myData)
		if err != nil {
			b.Errorf("could not create/update record : %v", err.Error())
			return
		}
		log.Printf("status : %v", status)

		record, err := mongoFetchRecord(key)
		if err != nil {
			b.Errorf("could not fetch record with key (%v)  %v", key, err.Error())
			return
		}
		if record.Key != key {
			b.Errorf("could not fetch record with key (%v)", key)
		}

		status, err = mongoDeleteRecord(key)
		if err != nil {
			b.Errorf("could not delete record : %v", err.Error())
		}
		if status == false {
			b.Errorf("could not delete record with key : %v", key)
		}
	}

}
