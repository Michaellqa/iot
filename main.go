package main

import (
	"github.com/Michaellqa/iot/internal/app"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFlags(log.Flags() | log.Llongfile | log.Lmicroseconds)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	cfg := readConfig()

	sv := app.Supervisor{}
	sv.Init(cfg)

	// todo: better sync

	done := make(chan struct{})
	go func() {
		sv.Start()
		done <- struct{}{}
	}()

	go func() {
		<-sigs
		log.Println("APP: shutting down ...")
		sv.Shutdown()
		os.Exit(1)
	}()

	<-done
	log.Println("APP: finished successfully")
}

func readConfig() app.Config {
	return app.Config{
		Generators: []app.GeneratorConfig{
			{
				TimeoutSec:    7,
				SendPeriodSec: 1,
				DataSources: []app.DataSourceConfig{
					{Id: "data_1", InitValue: 50, MaxChangeStep: 5},
				},
			},
		},
		Queue: app.QueueConfig{Size: 50},
		Aggregators: []app.AggregatorConfig{
			{AggregationPeriodSec: 2, SubIds: []string{"data_1"}},
		},
		Storage: app.StorageConfig{Type: 0},
	}
}
