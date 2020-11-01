package aggregation

import (
	"strconv"
	"testing"
	"time"
)

var writeDuration = 100 * time.Millisecond

type mockStorage struct {
	t *testing.T
}

func (s *mockStorage) Write(record Record) error {
	time.Sleep(writeDuration)
	s.t.Log(record)
	return nil
}

func TestAsyncStorage(t *testing.T) {
	store := &mockStorage{t: t}
	as := newAsyncStorage(store)

	now := time.Now()
	lap := func() {
		t.Log(time.Now().Sub(now))
	}

	for i := 0; i < 20; i++ {
		as.Add(Record{Id: strconv.Itoa(i), Value: float64(i)})
		lap()
	}
	go func() {
		time.Sleep(5 * writeDuration)
		as.Close()
		t.Log("close")
	}()

	as.Wait()
}
