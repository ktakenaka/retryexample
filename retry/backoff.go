package retry

import (
	"math"
	"math/rand"
	"time"
)

var stop time.Duration = -1

type backoff struct {
	retryTimes        uint
	maxRetryTimes     uint
	interval          time.Duration
	multiplier        float64
	maxJitterInterval time.Duration
	checkRetryable    checkRetryable
}

func defaultBackoff() *backoff {
	return &backoff{
		retryTimes:        0,
		maxRetryTimes:     5,
		interval:          100.0 * time.Millisecond,
		multiplier:        2.0,
		maxJitterInterval: 30.0 * time.Millisecond,
		checkRetryable:    func(err error) bool { return false },
	}
}

func (b *backoff) Next() time.Duration {
	if b.retryTimes >= b.maxRetryTimes {
		return stop
	}

	b.retryTimes++
	i := b.interval * time.Duration(math.Pow(b.multiplier, float64(b.retryTimes)))
	return b.randomize(i)
}

func (b *backoff) randomize(i time.Duration) time.Duration {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	min := float64(i) - float64(b.maxJitterInterval)
	max := float64(i) + float64(b.maxJitterInterval)
	return time.Duration(min + ((max - min) * r.Float64()))
}

func MaxRetryTimes(mt uint) option {
	return func(eb *backoff) {
		eb.maxRetryTimes = mt
	}
}

func Interval(i time.Duration) option {
	return func(eb *backoff) {
		eb.interval = i
	}
}

func Multiplier(m float64) option {
	return func(eb *backoff) {
		eb.multiplier = m
	}
}
func MaxJitterInterval(ji time.Duration) option {
	return func(eb *backoff) {
		eb.maxJitterInterval = ji
	}
}

func CheckRetryable(cr checkRetryable) option {
	return func(eb *backoff) {
		eb.checkRetryable = cr
	}
}
