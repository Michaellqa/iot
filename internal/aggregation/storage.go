package aggregation

import (
	"log"
)

type Storage interface {
	Write(record Record) error
}

func newAsyncStorage(store Storage) *asyncStorage {
	s := &asyncStorage{
		// todo: check capacities
		adding:  make(chan Record, 1),
		closing: make(chan struct{}, 1),
		waiting: make(chan struct{}),
		done:    make(chan struct{}),
		store:   store,
	}
	go s.process()
	return s
}

type asyncStorage struct {
	store Storage
	q     fifo

	adding  chan Record
	closing chan struct{}
	waiting chan struct{}
	done    chan struct{}
}

// Add adds the record to the queue.
func (s *asyncStorage) Add(r Record) {
	select {
	case <-s.done:
		return
	default:
	}
	s.adding <- r
}

// Wait stops receiving new records and blocks until all the entries in the queue are saved.
func (s *asyncStorage) Wait() {
	select {
	case <-s.done:
		return
	default:
	}
	s.waiting <- struct{}{}
	<-s.done
}

// Close interrupts processing immediately if not in the middle of saving a record.
func (s *asyncStorage) Close() {
	select {
	case <-s.done:
		return
	default:
	}
	s.closing <- struct{}{}
	<-s.done
}

func (s *asyncStorage) process() {
	var draining, saveDone chan struct{}

	for {
		// gives priority to brake the cycle over the draining channel case.
		select {
		case <-s.closing:
			close(s.done)
			return
		case <-s.waiting:
			if s.q.len == 0 {
				close(s.done)
				return
			}
			// disable accepting new records
			s.adding = nil
		default:
		}

		select {
		case record := <-s.adding:
			s.q.add(record)
			if draining == nil {
				// enable draining the queue
				draining = make(chan struct{})
				close(draining)
			}

		case <-draining:
			val := s.q.get()
			if val == nil {
				select {
				case <-s.done:
					return
				default:
				}
				// queue is empty, disable storing
				draining = nil
				continue
			}
			draining = nil
			saveDone = make(chan struct{})
			go func() {
				record := val.(Record)
				if err := s.store.Write(record); err != nil {
					log.Println(err)
				}
			}()

		case <-saveDone:
			draining = nil
			saveDone = make(chan struct{})

		case <-s.waiting:
			if s.q.len == 0 {
				close(s.done)
				return
			}
			// disable accepting new records
			s.adding = nil
		case <-s.closing:
			close(s.done)
			return
		}
	}
}
