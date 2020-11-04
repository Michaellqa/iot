package messaging

import (
	"testing"
)

func TestBroker(t *testing.T) {
	b := NewBroker(50)

	startAggregate := func() {
		sub := b.Subscribe()
		total := 0
		for val := range sub {
			total += val.(int)
		}
		t.Log(total)
	}
	// test result are correct
	b.Publish(1)

	for i := 0; i < 3; i++ {
		go startAggregate()
		for j := 0; j < 100000; j++ {
			b.Publish(j)
			//runtime.Gosched()
		}
	}
	b.Close()

	// test no races if published in parallel

}
