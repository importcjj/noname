package dingding

import (
	"context"

	"github.com/royeo/dingrobot"
)

type DingdingRobot struct {
	robot   dingrobot.Roboter
	keyword string
	enabled bool
}

func NewRobot(keyword, webhook string) *DingdingRobot {
	robot := dingrobot.NewRobot(webhook)
	return &DingdingRobot{
		robot:   robot,
		keyword: keyword,
		enabled: len(webhook) > 0,
	}
}

func (r *DingdingRobot) Send(ctx context.Context, content string) error {
	return r.robot.SendText(r.keyword+"\n"+content, nil, false)
}
