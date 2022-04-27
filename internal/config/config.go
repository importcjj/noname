package config

import (
	"fmt"
	"time"

	"github.com/importcjj/ddxq/pkg/serverchan"

	"github.com/importcjj/ddxq/internal/boost"
	"github.com/importcjj/ddxq/pkg/api"
	"github.com/importcjj/ddxq/pkg/dingding"
	"github.com/jinzhu/configor"
)

type Config struct {
	API             api.Config `yaml:"api" json:"api"`
	UseBalance      bool       `yaml:"use_balance" json:"use_balance"`
	CartInterval    string     `yaml:"cart_interval" json:"cart_interval" default:"2m"`
	ReserveInterval string     `yaml:"reserve_interval" json:"reserve_interval" default:"2s"`
	HomeInterval    string     `yaml:"home_interval" json:"home_interval" default:"1m"`

	Dingding   dingding.Config   `yaml:"dingding" json:"dingding"`
	ServerChan serverchan.Config `yaml:"serverChan" json:"serverChan"`
	BoostMode  boost.Config      `yaml:"boost_mode" json:"boost_mode"`
}

func Load(filepath string) (Config, error) {
	var config Config
	err := configor.Load(&config, filepath)

	if err != nil {
		return config, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config, nil
}

func (c *Config) NewMode() (*Mode, error) {
	boostMode, err := boost.New(c.BoostMode)
	if err != nil {
		return nil, fmt.Errorf("无法创建boost: %w", err)
	}

	cartInterval, err := time.ParseDuration(c.CartInterval)
	if err != nil {
		return nil, err
	}

	reserveInterval, err := time.ParseDuration(c.ReserveInterval)
	if err != nil {
		return nil, err
	}

	homeInterval, err := time.ParseDuration(c.HomeInterval)
	if err != nil {
		return nil, err
	}

	mode := &Mode{
		BoostMode:       *boostMode,
		useBalance:      c.UseBalance,
		cartInterval:    cartInterval,
		reserveInterval: reserveInterval,
		homeInterval:    homeInterval,
	}

	return mode, nil
}
