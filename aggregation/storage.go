package aggregation

import (
	"log"
)

type Storage interface {
	Write(record Record) error
}

type AsyncStorage interface {
	Add(r Record)
	Wait()
	Close()
}

func NewAsyncStorage(fifo Fifo, store Storage) *asyncStorage {
	s := &asyncStorage{
		fifo:  fifo,
		store: store,

		adding:  make(chan Record, 1),
		closing: make(chan struct{}, 1),
		waiting: make(chan struct{}),
		done:    make(chan struct{}, 1),
	}
	go s.process()
	return s
}

type asyncStorage struct {
	fifo  Fifo
	store Storage

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

// Wait blocks until all the entries in the queue are saved and breaks the process loop.
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
	defer close(s.done)

	var draining, saveDone chan struct{}
	quitOnSavedLast := make(chan struct{}, 1)

	for {
		// gives priority to brake the cycle over the draining channel case.
		select {
		case <-s.closing:
			return
		case <-s.waiting:
			if s.fifo.Len() == 0 {
				return
			}
			quitOnSavedLast <- struct{}{}
		default:
		}

		select {
		case record := <-s.adding:
			s.fifo.Add(record)
			if draining == nil && saveDone == nil {
				// enable draining the queue
				draining = make(chan struct{})
				close(draining)
			}

		case <-draining:
			val := s.fifo.Get()
			if val == nil {
				select {
				case <-quitOnSavedLast:
					return
				default:
				}
				// queue is empty, disable storing
				draining = nil
				continue
			}
			draining = nil
			saveDone = make(chan struct{}, 1)
			go func() {
				record := val.(Record)
				if err := s.store.Write(record); err != nil {
					log.Println(err)
				}
				saveDone <- struct{}{}
			}()
		case <-saveDone:
			saveDone = nil
			draining = make(chan struct{})
			close(draining)

		case <-s.waiting:
			if s.fifo.Len() == 0 {
				return
			}
			quitOnSavedLast <- struct{}{}

		case <-s.closing:
			return
		}
	}
}
