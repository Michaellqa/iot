package generation

import (
	"github.com/Michaellqa/iot/internal/messaging"
	"log"
	"time"
)

type Metrics struct {
	Id    string
	Value int
}

type Source interface {
	ReadValue() Metrics
}

func New(
	broker messaging.Broker,
	sources []Source,
	sendPeriod time.Duration,
	timeoutPeriod time.Duration,
) *Generator {
	return &Generator{
		broker:        broker,
		sources:       sources,
		sendPeriod:    sendPeriod,
		timeoutPeriod: timeoutPeriod,
		stop:          make(chan struct{}, 1),
	}
}

type Generator struct {
	broker  messaging.Broker
	sources []Source

	sendPeriod    time.Duration
	timeoutPeriod time.Duration
	stop          chan struct{}
}

// Start starts the process of periodically asking its data sources for values and
// sends them to the receiver... ?
// Returns a channel that will be closed when the generator stops.
func (g *Generator) Start() <-chan struct{} {
	done := make(chan struct{}, 1)

	go func() {
		timeout := time.NewTimer(g.timeoutPeriod)
		ticker := time.NewTicker(g.sendPeriod)
		defer func() {
			ticker.Stop()
			if timeout.Stop() {
				// todo: wtf
				<-timeout.C
			}
			close(done)
		}()

		for {
			select {
			case <-ticker.C:
				g.generate()

			case <-timeout.C:
				log.Println("GENERATOR: timed out")
				return
			case <-g.stop:
				log.Println("GENERATOR: stopped")
				return
			}
		}
	}()

	return done
}

// Stop must not be called multiple times.
func (g *Generator) Stop() {
	g.stop <- struct{}{}
}

// generate sends metrics from all the data sources.
func (g *Generator) generate() {
	for _, s := range g.sources {
		val := s.ReadValue()
		g.broker.Publish(val)
	}
}
