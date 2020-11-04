package storage

import (
	"fmt"
	"github.com/Michaellqa/iot/aggregation"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const root = "/app/data"

func NewFileStorage() aggregation.Storage {
	return &fileStorage{filename: time.Now().Format("2006-01-02_150405")}
}

type fileStorage struct {
	filename string
	mu       sync.Mutex
}

func (s *fileStorage) Write(record aggregation.Record) error {
	if _, err := os.Stat(root); os.IsNotExist(err) {
		if err = os.Mkdir(root, 0700); err != nil {
			panic(err)
		}
	}

	path := filepath.Join(root, s.filename)

	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s: %s: %v\n", record.Timestamp.Format("15:04:05.000"), record.Id, record.Value)
	_, err = f.WriteString(line)
	return err
}
