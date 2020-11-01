package storage

import (
	"github.com/Michaellqa/iot/internal/aggregation"
	"log"
)

type Console struct{}

func (s Console) Write(r aggregation.Record) error {
	log.Printf("%s: %v\n", r.Id, r.Value)
	return nil
}
