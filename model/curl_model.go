/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-19
* Time: 09:51
 */

package model

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// curl参数解析
type CURL struct {
	Data map[string][]string
}

// 从文件中解析curl
func ParseTheFile(path string) (curl *CURL, err error) {

	if path == "" {
		err = errors.New("路径不能为空")

		return
	}

	curl = &CURL{
		Data: make(map[string][]string),
	}

	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())

		return
	}

	defer func() {
		file.Close()
	}()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		err = errors.New("读取文件失败:" + err.Error())

		return
	}

	dataStr := string(data)
	curl = newCurlFromStr(dataStr)

	return
}

func ParseTheFileMulti(path string) (curls []*CURLJson, err error) {
	if path == "" {
		err = errors.New("路径不能为空")

		return
	}

	curls = make([]*CURLJson, 0)

	file, err := os.Open(path)
	if err != nil {
		err = errors.New("打开文件失败:" + err.Error())

		return
	}

	defer func() {
		file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		oneLineStr := scanner.Text()
		curl := newCurlJsonFromJsonStr(oneLineStr)
		curls = append(curls, curl)
	}
	return curls, nil
}

func newCurlFromStr(dataStr string) *CURL {
	curl := &CURL{
		Data: make(map[string][]string),
	}

	for true {
		index := strings.Index(dataStr, "'")
		if index <= 0 {
			break
		}

		key := strings.TrimSpace(dataStr[:index])
		key = strings.ReplaceAll(key, "\n", "")

		dataStr = dataStr[index+1:]

		index = strings.Index(dataStr, "'")

		if index <= 0 {
			break
		}

		value := dataStr[:index]

		dataStr = dataStr[index+1:]

		curl.Data[key] = append(curl.Data[key], value)

	}

	return curl
}

type CURLJson struct {
	URL string `json:"url"`
	Method string `json:"method"`
	Header map[string]string `json:"header"`
	Data string `json:"data"`
}

func newCurlJsonFromJsonStr(jsonStr string) *CURLJson {
	var curlJson CURLJson
	err := json.Unmarshal([]byte(jsonStr), &curlJson)
	if err != nil {
		log.Panicf("unmashal json error: %v", err)
	}
	return &curlJson
}

func (c *CURL) String() (url string) {
	curlByte, _ := json.Marshal(c)

	return string(curlByte)
}

// GetUrl
func (c *CURL) GetUrl() (url string) {
	value, ok := c.Data["curl"]
	if !ok {

		return
	}

	if len(value) <= 0 {

		return
	}

	url = value[0]

	return
}

// GetMethod
func (c *CURL) GetMethod() (method string) {
	method = "GET"

	if _, ok := c.Data["--data"]; ok {
		method = "POST"

		return
	}

	// TODO::目前发送不了
	if _, ok := c.Data["--data-binary $"]; ok {
		method = "POST"

		return
	}

	value, ok := c.Data["-X"]
	if !ok {

		return
	}

	if len(value) <= 0 {

		return
	}

	method = strings.ToUpper(value[0])

	return
}

// GetHeaders
func (c *CURL) GetHeaders() (headers map[string]string) {
	headers = make(map[string]string, 0)

	value, ok := c.Data["-H"]
	if !ok {

		return
	}

	for _, v := range value {
		index := strings.Index(v, ":")
		if index < 0 {
			continue
		}

		vIndex := index + 2
		if len(v) >= vIndex {
			headers[v[:index]] = v[vIndex:]
		}
	}

	return
}

// GetHeaders
func (c *CURL) GetHeadersStr() string {
	headers := c.GetHeaders()
	bytes, _ := json.Marshal(&headers)

	return string(bytes)
}

// 获取body
func (c *CURL) GetBody() (body string) {

	value, ok := c.Data["--data"]
	if !ok {
		// data-binary
		value, ok = c.Data["--data-binary $"]
		if !ok {

			return
		}
	}

	if len(value) <= 0 {

		return
	}

	// body = strings.NewReader(value[0])
	body = value[0]

	return
}
