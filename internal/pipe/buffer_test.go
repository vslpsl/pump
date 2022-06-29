package pipe

import (
	"github.com/stretchr/testify/require"
	"log"
	"sync"
	"testing"
	"time"
)

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		data := []byte("test1")
		q.Enqueue(data)
		<-time.After(time.Second)
		data = []byte("test2")
		q.Enqueue(data)
		<-time.After(time.Second)
		data = []byte("test3")
		q.Enqueue(data)
		q.Close()
		wg.Done()
	}()

	count := 0
	for {
		data, closed := q.Dequeue()
		if closed {
			break
		}
		log.Println(string(data))
		count++
	}

	wg.Wait()
	require.Equal(t, 3, count)
}
