package aggregation

import (
	"github.com/Michaellqa/iot/internal/generation"
	"github.com/Michaellqa/iot/internal/messaging"
	"log"
	"time"
)

type Record struct {
	Id    string
	Value float64
}

func NewAggregator(
	broker messaging.Broker,
	store Storage,
	period time.Duration,
	subIds []string,
) *Aggregator {
	a := &Aggregator{
		buf:    newBuffer(),
		broker: broker,
		period: period,
		store:  newAsyncStorage(store),
		subIds: make(map[string]struct{}),
	}
	for _, id := range subIds {
		a.subIds[id] = struct{}{}
	}
	return a
}

type Aggregator struct {
	buf    *buffer
	broker messaging.Broker
	store  *asyncStorage

	period time.Duration
	subIds map[string]struct{}

	done chan struct{}
}

func (a *Aggregator) Start() {
	msgCh := a.broker.Subscribe()

	ticker := time.NewTicker(a.period)
	defer ticker.Stop()

	// should I add another case to force break the cycle or I could rely on closing of the broker?
	for {
		select {
		case msg, alive := <-msgCh:
			if !alive {
				// flush what is left in the buffer.
				a.aggregate()
				return
			}
			metrics := msg.(generation.Metrics)
			if _, ok := a.subIds[metrics.Id]; ok {
				a.buf.add(metrics)
			}

		case <-ticker.C:
			a.aggregate()
		}
		// todo: add stop
	}

	// log: aggregation has stopped
}

// Wait is safe to be called if already stopped.
func (a *Aggregator) Wait() {
	a.store.Wait()
}

// Stop is safe to be called if already stopped.
func (a *Aggregator) Stop() {
	a.store.Close()
}

func (a *Aggregator) aggregate() {
	records := a.buf.flush()
	log.Println("DEBUG: aggregated data:", records)
	for _, r := range records {
		a.store.Add(r)
	}
}

func newBuffer() *buffer {
	return &buffer{data: make(map[string][]int)}
}

type buffer struct {
	data map[string][]int
}

func (b *buffer) add(msg generation.Metrics) {
	b.data[msg.Id] = append(b.data[msg.Id], msg.Value)
}

// flush aggregates and returns stored data and clears the buffer.
func (b *buffer) flush() []Record {
	results := make([]Record, 0, len(b.data))

	for id, values := range b.data {
		if len(values) == 0 {
			continue
		}

		sum := 0
		for _, v := range values {
			sum += v
		}
		avg := float64(sum) / float64(len(values))
		results = append(results, Record{Id: id, Value: avg})
	}

	// clear the buffer
	// b.data = make(map[string][]int)

	// clear the buffer, reuse memory
	for k := range b.data {
		b.data[k] = b.data[k][:0]
	}

	return results
}
