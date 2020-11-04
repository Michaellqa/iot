package storage

import (
	"github.com/Michaellqa/iot/aggregation"
	"log"
)

type Console struct{}

func (s Console) Write(r aggregation.Record) error {
	log.Printf("CONSOLE: %s: %v\n", r.Id, r.Value)
	return nil
}
