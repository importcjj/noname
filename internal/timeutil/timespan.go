package timeutil

import (
	"time"
)

const layout = "15:04:05"

type Span struct {
	start time.Time
	end   time.Time
}

func NewSpan(start, end string) (*Span, error) {
	startTime, err := time.Parse(layout, start)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse(layout, end)
	if err != nil {
		return nil, err
	}

	span := &Span{
		start: startTime,
		end:   endTime,
	}

	return span, nil
}

func (s *Span) Include(check time.Time) bool {
	var (
		start = s.start
		end   = s.end
	)

	t := check.Format(layout)
	check, _ = time.Parse(layout, t)

	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}

	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
