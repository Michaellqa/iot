package storage

import (
	"github.com/Michaellqa/iot/aggregation"
	"log"
	"time"
)

type SlowStorage struct{}

func (s SlowStorage) Write(r aggregation.Record) error {
	time.Sleep(3 * time.Second)
	log.Printf("CONSOLE: %s: %v\n", r.Id, r.Value)
	return nil
}
