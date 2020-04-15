/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-15
* Time: 13:44
 */

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"go-stress-testing/model"
	"go-stress-testing/server"
	"net/http"
	"runtime"
	"strings"
	"time"
)

// go 实现的压测工具
//
// 编译可执行文件
//go:generate go build main.go
func main() {

	runtime.GOMAXPROCS(1)

	var (
		concurrency uint64 // 并发数
		totalNumber uint64 // 请求总数(单个并发)
		timeout     int    // 超时秒数
		debugStr    string // 是否是debug
		requestUrl  string // 压测的url 目前支持，http/https ws/wss
		path        string // curl文件路径 http接口压测，自定义参数设置
		verify      string // verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json
		paths       string // curl文件路径 http接口压测，自定义参数设置，用于存放多条请求
	)

	flag.Uint64Var(&concurrency, "c", 1, "并发数")
	flag.Uint64Var(&totalNumber, "n", 1, "请求总数")
	flag.IntVar(&timeout, "t", 0, "超时秒数")
	flag.StringVar(&debugStr, "d", "false", "调试模式")
	flag.StringVar(&requestUrl, "u", "", "压测地址")
	flag.StringVar(&path, "p", "", "curl文件路径")
	flag.StringVar(&verify, "v", "", "验证方法 http 支持:statusCode、json webSocket支持:json")
	flag.StringVar(&paths, "ps", "", "批量处理curl文件请求")

	// 解析参数
	flag.Parse()
	if concurrency == 0 || (totalNumber == 0 && paths == "") || (requestUrl == "" && path == "" && paths == "") {
		fmt.Printf("示例: go run main.go -c 1 -n 1 -u https://www.baidu.com/ \n")
		fmt.Printf("压测地址或curl路径必填 \n")
		fmt.Printf("当前请求参数: -c %d -n %d -d %v -u %s \n", concurrency, totalNumber, debugStr, requestUrl)

		flag.Usage()

		return
	}

	debug := strings.ToLower(debugStr) == "true"

	// setup http transport
	// Customize the Transport to have larger connection pool
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPointer, ok := defaultRoundTripper.(*http.Transport)
	if !ok {
		panic(fmt.Sprintf("defaultRoundTripper not an *http.Transport"))
	}
	defaultTransport := *defaultTransportPointer // dereference it to get a copy of the struct that the pointer points to
	defaultTransport.MaxIdleConns = int(concurrency)
	defaultTransport.MaxIdleConnsPerHost = int(concurrency)
	defaultTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if paths == "" {
		request, err := model.NewRequest(requestUrl, verify, time.Second*time.Duration(timeout), debug, path)
		if err != nil {
			fmt.Printf("参数不合法 %v \n", err)

			return
		}

		fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, totalNumber)
		request.Print()

		// 开始处理
		server.Dispose(concurrency, totalNumber, request)
	} else {
		requests, err := model.NewRequestMulti(paths, "ants", time.Second*time.Duration(timeout), debug)
		if err != nil {
			fmt.Printf("参数不合法 %v \n", err)

			return
		}
		fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", concurrency, len(requests))
		server.DisposeMulti(concurrency, requests)
	}

	return
}
