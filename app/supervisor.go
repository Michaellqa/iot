package app

import (
	"context"
	"github.com/Michaellqa/iot/aggregation"
	"github.com/Michaellqa/iot/generation"
	"github.com/Michaellqa/iot/messaging"
	"github.com/Michaellqa/iot/storage"
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
	switch cfg.StorageType {
	case 1:
		store = storage.NewFileStorage()
	case 2:
		store = storage.SlowStorage{}
	default:
		store = storage.Console{}
	}
	fifo := &aggregation.ListFifo{}
	asyncStore := aggregation.NewAsyncStorage(fifo, store)

	// create aggregators
	s.aggregators = make([]*aggregation.Aggregator, len(cfg.Aggregators))

	for i, agg := range cfg.Aggregators {
		aggPeriod := time.Duration(agg.AggregationPeriodSec) * time.Second
		s.aggregators[i] = aggregation.NewAggregator(s.broker, asyncStore, aggPeriod, agg.SubIds)
	}
}

func (s *Supervisor) Start(ctx context.Context) {
	done := make(chan struct{}, 1)

	go func() {
		for _, agg := range s.aggregators {
			go agg.Start()
		}

		wg := sync.WaitGroup{}
		wg.Add(len(s.generators))
		for _, gen := range s.generators {
			go func() {
				genDone := gen.Start()
				<-genDone
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

		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		s.shutdown()
	}
}

func (s *Supervisor) shutdown() {
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
