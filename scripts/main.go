package main

import (
	"encoding/json"
	"fmt"
	"go-stress-testing/model"
	"os"
	"strconv"
	"strings"
)

const (
	//HOST = "http://tech-ants-charge.test.za-tech.net"
	HOST = "http://localhost:8080"

	URL_REGISTER_USER = "/v1/charge/seckill/register"
	URL_GET_STOCK     = "/v1/charge/seckill/getStock/%s"
	URL_GET_URL       = "/v1/charge/seckill/getUrl/%s"
	URL_ORDER         = "/v1/charge/seckill/order/%s"
)

var baseHeader = map[string]string{
	"Content-Type": "application/json",
	"userId":       "26",
}

var LINES = 20000
var EVENT_CODE = "LBJN202000420"

func generateOrder(f *os.File) {
	curlNum := LINES
	eventCode := EVENT_CODE

	prdId := "9"
	itemId := "32"
	accountType := "1"
	quantity := "1"
	userId := 77
	for i := 0; i < curlNum; i++ {
		baseHeader["userId"] = strconv.Itoa(userId)
		postStr := order(eventCode, prdId, itemId, accountType, quantity)
		f.WriteString(postStr + "\n")
		userId += 1
		if userId >= 20075 {
			break
		}
	}
}

func generateRegister(f *os.File) {
	curlNum := LINES
	eventCode := EVENT_CODE

	userId := 77
	for i := 0; i < curlNum; i++ {
		baseHeader["userId"] = strconv.Itoa(userId)
		getUrlStr := registerUser(eventCode)
		f.WriteString(getUrlStr + "\n")
		userId += 1
		if userId >= 20075 {
			break
		}
	}
}

func gen(filename string, handler func(*os.File)) {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Printf("open file %s failed: %v\n", filename, err)
		return
	}
	defer f.Close()
	handler(f)
}

func main() {
	gen("clst-reg", generateRegister)
	gen("clst-order", generateOrder)
}

func getUrl(eventCode string) string {
	url := fmt.Sprintf(HOST+URL_GET_URL, eventCode)
	return getStr(url, baseHeader, nil)
}

func registerUser(eventCode string) string {
	url := HOST + URL_REGISTER_USER
	data := map[string]string{
		"eventCode": eventCode,
	}
	return postStr(url, baseHeader, data)
}

func order(eventCode, prdId, itemId, accoutType, quantity string) string {
	orderKey := "5NVpDSmrqnSTbx3RQaeACdezjd885Fsv"

	url := fmt.Sprintf(HOST+URL_ORDER, orderKey)
	data := map[string]string{
		"eventCode":       eventCode,
		"productId":       prdId,
		"itemId":          itemId,
		"accountType":     accoutType,
		"rechargeAccount": "3298423",
		"quantity":        quantity,
		"did":             "did",
		"token":           "xx",
		"sid":             "sid",
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
