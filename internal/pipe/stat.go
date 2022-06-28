package pipe

import (
	"fmt"
	"time"
)

type Stat struct {
	StartedAt  time.Time
	FinishedAt time.Time
	BytesPiped int64
}

func (s *Stat) Start() {
	s.StartedAt = time.Now()
}

func (s *Stat) Stop() {
	s.FinishedAt = time.Now()
}

func (s *Stat) String() string {
	duration := s.FinishedAt.Sub(s.StartedAt)
	return fmt.Sprintf("duration: %d ns\nbytes piped: %d\nrate: %f bytes/sec\n",
		duration.Nanoseconds(),
		s.BytesPiped,
		float64(s.BytesPiped)/duration.Seconds(),
	)
}
