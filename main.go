package main

import (
	"context"
	"dingdong/dingding"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/antlabs/pcurl"
)

var curlFile = flag.String("curl_file", "raw_request.txt", "请提供CURL请求文件地址")
var webhook = flag.String("webhook", "", "钉钉群webhook")

func main() {
	flag.Parse()

	if len(*curlFile) == 0 {
		log.Fatal("请提供CURL请求文件地址")
	}

	body, err := ioutil.ReadFile(*curlFile)
	if err != nil {
		log.Fatal(err)
	}

	configString := string(body)
	configString = strings.Replace(configString, "--data-binary", "--data", 1)

	var robot = dingding.NewRobot("信号", *webhook)

	var client = &http.Client{
		Timeout: time.Second * 5,
	}

	log.Println("正在监控叮咚运力。。。")

	for {
		time.Sleep(1500 * time.Millisecond)

		baseReq, err := pcurl.ParseAndRequest(configString)

		if err != nil {
			log.Fatalf("curl文件解析失败: %v", err)
		}

		baseReq.Header.Del("Accept-Encoding")
		baseReq.Header.Del("user-agent")
		// baseReq.Header.Set("user-agent", "Mozilla/5.0 ( NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/53.0.2785.143 Safari/537.36 MicroMessenger/7.0.9.501 NetType/WIFI MiniProgramEnv/ Wechat")

		resp, err := client.Do(baseReq)
		if err != nil {
			log.Printf("无法请求: %v", err)
			continue
		}

		defer resp.Body.Close()

		var response = new(Response)
		var timeDatas []*MultiReserveTimeResponse
		response.Data = &timeDatas

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("无法读取返回: %v", err)
			continue
		}

		err = json.Unmarshal(b, &response)
		if err != nil {
			log.Printf("无法解析返回: %v", err)
			log.Println(string(b))
			continue
		}

		if len(timeDatas) == 0 {
			log.Println("请求异常")
			log.Println(string(b))
		}
		for _, timeData := range timeDatas {
			if len(timeData.Time) == 0 {
				continue
			}
			for _, time := range timeData.Time {
				for _, period := range time.Times {
					if !period.FullFlag {
						msg := fmt.Sprintf("注意！%s - %s 可约", period.StartTime, period.EndTime)
						robot.Send(context.Background(), msg)
						log.Println(msg)
					}
				}
			}

		}

	}

}

type Response struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

type MultiReserveTimeResponse struct {
	Time []struct {
		Times []struct {
			FullFlag  bool   `json:"fullFlag"`
			StartTime string `json:"start_time"`
			EndTime   string `json:"end_time"`
		} `json:"times"`
	} `json:"time"`
}
