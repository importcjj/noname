package boost

import (
	"time"

	"github.com/importcjj/ddxq/internal/timeutil"
)

type Mode struct {
	config *Config

	warmUpTimeSpan []*timeutil.Span
	timeSpan       []*timeutil.Span

	cartInterval    time.Duration
	reserveInterval time.Duration
	recheckInterval time.Duration
	reorderInterval time.Duration
}

func New(config Config) (*Mode, error) {
	var warmUpTimeSpan []*timeutil.Span
	var timeSpan []*timeutil.Span
	for _, spanVal := range config.WarmUpTimeSpan {
		span, err := timeutil.NewSpan(spanVal.Start, spanVal.End)
		if err != nil {
			return nil, err
		}

		warmUpTimeSpan = append(warmUpTimeSpan, span)
	}

	for _, spanVal := range config.TimeSpan {
		span, err := timeutil.NewSpan(spanVal.Start, spanVal.End)
		if err != nil {
			return nil, err
		}

		timeSpan = append(timeSpan, span)
	}

	cartInterval, err := time.ParseDuration(config.CartInterval)
	if err != nil {
		return nil, err
	}

	reserveInterval, err := time.ParseDuration(config.ReserveInterval)
	if err != nil {
		return nil, err
	}

	recheckInterval, err := time.ParseDuration(config.RecheckInterval)
	if err != nil {
		return nil, err
	}

	reorderInterval, err := time.ParseDuration(config.ReorderInterval)
	if err != nil {
		return nil, err
	}

	mode := &Mode{
		config:         &config,
		warmUpTimeSpan: warmUpTimeSpan,
		timeSpan:       timeSpan,

		cartInterval:    cartInterval,
		reserveInterval: reserveInterval,
		recheckInterval: recheckInterval,
		reorderInterval: reorderInterval,
	}

	return mode, nil
}

func (mode *Mode) Enable() bool {
	return mode.config.Enable
}

func (mode *Mode) WarmUpBoostTime() bool {
	now := time.Now()
	for _, span := range mode.warmUpTimeSpan {
		if span.Include(now) {
			return true
		}
	}

	return false
}

func (mode *Mode) BoostTime() bool {
	now := time.Now()
	for _, span := range mode.timeSpan {
		if span.Include(now) {
			return true
		}
	}

	return false
}

func (mode *Mode) GetCartInterval() time.Duration    { return mode.cartInterval }
func (mode *Mode) GetReserveInterval() time.Duration { return mode.reserveInterval }
func (mode *Mode) GetRecheckInterval() time.Duration { return mode.recheckInterval }
func (mode *Mode) GetReorderInterval() time.Duration { return mode.reorderInterval }
func (mode *Mode) UseBalance() bool                  { return mode.config.UseBalance }
