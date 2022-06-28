package limiter

import (
	"log"
	"time"
)

type Limiter struct {
	rate   int64
	leased int64
	ts     int64
}

func New(rate int64) *Limiter {
	return &Limiter{
		rate: rate,
	}
}

func (l *Limiter) Lease(size int64) {
	now := time.Now().UnixNano()
	if l.leased+size > l.rate {
		deadline := now - l.ts
		if deadline > time.Second.Nanoseconds() {
			l.leased = size
			l.ts = now
			return
		}
		deadline = time.Second.Nanoseconds() - deadline
		log.Printf("limiter: waiting: %d ns\n", deadline)
		<-time.After(time.Duration(deadline))
		l.leased = size
		l.ts = now
		return
	}

	if now-l.ts > time.Second.Nanoseconds() {
		l.leased = size
		l.ts = now
		return
	}

	l.leased += size
	log.Printf("leased: %d", l.leased)
}
