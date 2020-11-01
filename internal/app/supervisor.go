package app

import (
	"github.com/Michaellqa/iot/internal/aggregation"
	"github.com/Michaellqa/iot/internal/generation"
	"github.com/Michaellqa/iot/internal/messaging"
	"github.com/Michaellqa/iot/internal/storage"
	"log"
	"sync"
	"time"
)

type Supervisor struct {
	aggregators []*aggregation.Aggregator
	generators  []*generation.Generator
	broker      messaging.Broker
}

func (s *Supervisor) Init(cfg Config) {
	// create message broker
	s.broker = messaging.NewBroker(cfg.Queue.Size)

	// create pool of generators
	s.generators = make([]*generation.Generator, len(cfg.Generators))

	for i, gen := range cfg.Generators {
		dataSources := make([]generation.Source, len(gen.DataSources))
		for i, ds := range gen.DataSources {
			dataSources[i] = generation.NewRandomSource(ds.Id, ds.InitValue, ds.MaxChangeStep)
		}

		s.generators[i] = generation.New(s.broker,
			dataSources,
			time.Duration(gen.SendPeriodSec)*time.Second,
			time.Duration(gen.TimeoutSec)*time.Second)
	}

	// create storage
	var store aggregation.Storage
	switch cfg.Storage.Type {
	case 1:
		var filename string
		if cfg.Storage.Options != nil {
			filename = cfg.Storage.Options.Filename
		}
		store = storage.NewFileStorage(filename)
	default:
		store = storage.Console{}
	}

	// create aggregators
	s.aggregators = make([]*aggregation.Aggregator, len(cfg.Aggregators))

	for i, agg := range cfg.Aggregators {
		fifo := &aggregation.ListFifo{}
		asyncStore := aggregation.NewAsyncStorage(fifo, store)
		aggPeriod := time.Duration(agg.AggregationPeriodSec) * time.Second
		s.aggregators[i] = aggregation.NewAggregator(s.broker, asyncStore, aggPeriod, agg.SubIds)
	}
}

func (s *Supervisor) Start() {
	for _, agg := range s.aggregators {
		go agg.Start()
	}

	wg := sync.WaitGroup{}
	wg.Add(len(s.generators))
	for _, gen := range s.generators {
		go func() {
			done := gen.Start()
			<-done
			wg.Done()
		}()
	}
	wg.Wait()

	log.Println("SUPERVISOR: generators done")

	s.broker.Close()

	wg.Add(len(s.aggregators))
	for _, agg := range s.aggregators {
		agg.Wait()
		wg.Done()
	}
	wg.Wait()

	log.Println("SUPERVISOR: aggregators done")
}

func (s *Supervisor) Shutdown() {
	// stop all generators
	for _, gen := range s.generators {
		gen.Stop()
	}

	log.Println("SUPERVISOR: generators stopped")

	s.broker.Close()

	// wait until all aggregators saved their buffers
	wg := sync.WaitGroup{}
	wg.Add(len(s.aggregators))
	for _, agg := range s.aggregators {
		go func() {
			agg.Stop()
			wg.Done()
		}()
	}
	wg.Wait()
	log.Println("SUPERVISOR: aggregators stopped")
}
