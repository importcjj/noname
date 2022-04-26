package serverchan

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ServerChan struct {
	config *Config
	client *http.Client
}

func NewRobot(config Config) *ServerChan {

	robot := &ServerChan{
		config: &config,
	}

	if len(config.Key) > 0 {
		robot.client = http.DefaultClient
	}

	return robot
}

func (r *ServerChan) Send(ctx context.Context, content string) error {
	if r.config.Enable {
		reqUrl := fmt.Sprintf("https://sctapi.ftqq.com/%s.send", r.config.Key)
		values := url.Values{}
		values.Add("title", "抢菜通知")
		values.Add("desp", content)
		resp, err := http.PostForm(reqUrl, values)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

	return nil
}
