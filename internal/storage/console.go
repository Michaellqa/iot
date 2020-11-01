package storage

import (
	"github.com/Michaellqa/iot/internal/aggregation"
	"log"
	"time"
)

type Console struct{}

func (s Console) Write(r aggregation.Record) error {
	time.Sleep(3 * time.Second)
	log.Printf("CONSOLE: %s: %v\n", r.Id, r.Value)
	return nil
}
