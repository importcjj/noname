package dingding

import (
	"context"

	"github.com/royeo/dingrobot"
)

type DingdingRobot struct {
	robot  dingrobot.Roboter
	config *Config
}

func NewRobot(config Config) *DingdingRobot {

	robot := &DingdingRobot{
		config: &config,
	}

	if len(config.Hook) > 0 {
		robot.robot = dingrobot.NewRobot(config.Hook)
	}

	return robot
}

func (r *DingdingRobot) Send(ctx context.Context, content string) error {
	if r.config.Enable {
		return r.robot.SendText(r.config.Keyword+"\n"+content, nil, false)
	}

	return nil
}
