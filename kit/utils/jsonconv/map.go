package jsonconv

import (
	"encoding/json"
	"fmt"
)

func AnyToMap(obj interface{}) (map[string]interface{}, error) {
	var result map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func JsonToMap(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func AnyToMapNE(obj interface{}) map[string]interface{} {
	var result map[string]interface{}
	data, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("AnyToMapNE:json convert fail:", err)
		return nil
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("AnyToMapNE:json convert fail:", err)
		return nil
	}
	return result
}

func JsonToMapNE(data []byte) map[string]interface{} {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		fmt.Println("JsonToMapNE:json convert fail:", err)
		return nil
	}
	return result
}
