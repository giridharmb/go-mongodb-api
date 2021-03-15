package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

func getStringFromMapInterface(data map[string]interface{}) (string, error) {
	marshalledData, err := json.Marshal(data)
	stringToReturn := ""
	if err != nil {
		return stringToReturn, err
	}
	stringToReturn = string(marshalledData)
	return stringToReturn, nil
}

func getMapInterfaceFromString(jsonStr string) (map[string]interface{}, error) {
	returnData := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &returnData)
	if err != nil {
		return returnData, err
	}
	return returnData, nil
}

func getBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
