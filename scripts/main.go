package main

import (
	"encoding/json"
	"fmt"
	"go-stress-testing/model"
	"os"
	"strings"
)

const (
	HOST = "http://tech-ants-charge.test.za-tech.net"

	URL_REGISTER_USER = "/v1/charge/seckill/register"
	URL_GET_STOCK = "/v1/charge/seckill/getStock/%s"
	URL_GET_URL = "/v1/charge/seckill/getUrl/%s"
	URL_ORDER = "/v1/charge/seckill/order/%s"
)

var baseHeader = map[string]string{
	"Content-Type": "application/json",
	"userId": "26",
}

func main() {
	fileName := "curlList"
	curlNum := 10000
	eventCode := "LBJN202000420"

	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", fileName, err)
		return
	}
	defer f.Close()

	getUrlStr := registerUser(eventCode)
	for i := 0; i < curlNum; i++ {
		f.WriteString(getUrlStr + "\n")
	}

	//model.ParseTheFileMulti(fileName)
}

func getUrl(eventCode string) string {
	url := fmt.Sprintf(HOST + URL_GET_URL, eventCode)
	return getStr(url, baseHeader, nil)
}

func registerUser(eventCode string) string {
	url := HOST + URL_REGISTER_USER
	data := map[string]string{
		"eventCode": eventCode,
	}
	return postStr(url, baseHeader, data)
}

func getStr(url string, header, data map[string]string) string {
	curl := model.CURLJson{}

	dataKVList := make([]string, 0)
	for k, v := range data {
		kv := k + "=" + v
		dataKVList = append(dataKVList, kv)
	}
	dataStr := strings.Join(dataKVList, "&")

	curl.Method = "GET"
	curl.Header = header
	curl.URL = url + "?" + dataStr
	curl.Data = ""

	curlBytes, _ := json.Marshal(curl)

	return string(curlBytes)
}

func postStr(url string, header, data map[string]string) string {
	curl := model.CURLJson{}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println("marshal json error: ", err)
		return ""
	}
	dataStr := string(dataBytes)

	curl.Method = "POST"
	curl.Header = header
	curl.URL = url
	curl.Data = dataStr

	curlBytes, _ := json.Marshal(curl)

	return string(curlBytes)
}

