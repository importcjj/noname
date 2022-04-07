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

	for {
		time.Sleep(1 * time.Second)

		checkTimeRequest, err := pcurl.ParseAndRequest(configString)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(checkTimeRequest)
		if err != nil {
			log.Printf("无法请求: %v", err)
			continue
		}

		defer resp.Body.Close()

		var response = new(Response)
		var timeDatas []*MultiReserveTimeResponse
		response.Data = &timeDatas

		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			log.Printf("无法读取请求题: %v", err)
			continue
		}

		if len(timeDatas) == 0 {
			log.Println("请求异常")
			log.Println(response)
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
