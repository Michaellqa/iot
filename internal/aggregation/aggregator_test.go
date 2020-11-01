package aggregation_test

import (
	"github.com/Michaellqa/iot/internal/aggregation"
	"github.com/Michaellqa/iot/internal/generation"
	"github.com/Michaellqa/iot/internal/messaging"
	"github.com/Michaellqa/iot/internal/storage"
	"log"
	"testing"
	"time"
)

const (
	saveDuration      = 500 * time.Millisecond
	generationPeriod  = 300 * time.Millisecond
	aggregationPeriod = 1000 * time.Millisecond
)

type slowStorage struct {
}

func (s slowStorage) Write(r aggregation.Record) error {
	time.Sleep(saveDuration)
	log.Printf("%s: %v\n", r.Id, r.Value)
	return nil
}

func emulateGeneratorEvents(b messaging.Broker) {
	go func() {
		i := 0
		for {
			time.Sleep(generationPeriod)
			b.Publish(generation.Metrics{Id: "data_1", Value: i})
			i++
		}
	}()
}

func TestAggregator(t *testing.T) {
	log.SetFlags(log.Flags() | log.Lmicroseconds)

	broker := messaging.NewBroker(10)
	emulateGeneratorEvents(broker)

	store := storage.Console{}
	agg := aggregation.NewAggregator(broker, store, aggregationPeriod, []string{"data_1"})

	go agg.Start()

	time.Sleep(3 * time.Second)

	broker.Close()
	agg.Wait()
}
