package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	for {
		// 从命令行读取输入
		fmt.Print("Enter the URL: ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		// 发送 GET 请求
		resp, err := http.Get(input)
		if err != nil {
			log.Fatalln(err)
		}

		// 读取响应数据
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		// 打印响应数据
		log.Println(string(body))
	}
}
