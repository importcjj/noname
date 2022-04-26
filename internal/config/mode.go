package config

import (
	"time"

	"github.com/importcjj/ddxq/internal/boost"
)

type Mode struct {
	BoostMode       boost.Mode
	cartInterval    time.Duration
	reserveInterval time.Duration
}

func (mode *Mode) CartInterval() time.Duration {
	if mode.BoostMode.Enable() && mode.BoostMode.WarmUpBoostTime() {
		return mode.BoostMode.GetCartInterval()
	}
	return mode.cartInterval
}

func (mode *Mode) ReserveInterval() time.Duration {
	if mode.BoostMode.Enable() && mode.BoostMode.BoostTime() {
		return mode.BoostMode.GetReserveInterval()
	}
	return mode.reserveInterval
}

func (mode *Mode) RecheckInterval() time.Duration {
	return mode.BoostMode.GetRecheckInterval()
}

func (mode *Mode) ReorderInterval() time.Duration {
	return mode.BoostMode.GetReorderInterval()
}
