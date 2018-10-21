package benchmark

import "time"

type timer struct {
	starttime time.Time
}

func NewTimer() *timer {
	return &timer{}
}

func (t *timer) Start() *timer {
	t.starttime = time.Now()
	return t
}

func (t *timer) CostMillisecond() int64 {
	return time.Since(t.starttime).Nanoseconds() / (1000 * 1000)
}

func (t *timer) CostSecond() int64 {
	return int64(time.Since(t.starttime).Seconds())
}
