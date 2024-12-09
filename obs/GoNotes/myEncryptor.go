package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"reflect"
	"time"
)

type Data struct {
	Value int       `crypto:"true"`
	Date  time.Time `crypto:"true"`
	Alias string    `crypto:"false"`
}

func mapper(v interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	tp := val.Type()
	for i := 0; i < val.NumField(); i++ {
		tag := tp.Field(i).Tag.Get("crypto")
		fieldValue := val.Field(i).Interface()
		if tag == "true" {
			hash := sha256.Sum256([]byte(fmt.Sprintf("%v", fieldValue)))
			hashString := hex.EncodeToString(hash[:])
			result[tp.Field(i).Name] = hashString
		} else {
			result[tp.Field(i).Name] = fieldValue
		}
	}
	return result
}

func main() {
	data := Data{Value: 30123, Date: time.Now(), Alias: "Test"}
	dataMap := mapper(&data)
	fmt.Println(dataMap)
}
