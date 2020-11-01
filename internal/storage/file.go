package storage

import (
	"fmt"
	"github.com/Michaellqa/iot/internal/aggregation"
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewFileStorage(filename string) aggregation.Storage {
	return &fileStorage{filename: filename}
}

type fileStorage struct {
	filename string
}

// todo: relative paths
func (s *fileStorage) Write(record aggregation.Record) error {
	f, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_APPEND, 0700)
	if err != nil {
		return err
	}
	line := fmt.Sprintf("%s: %v", record.Id, record.Value)
	_, err = f.WriteString(line)
	return err
}
