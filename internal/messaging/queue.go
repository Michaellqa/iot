package messaging

import "sync"

type Broker interface {
	Publish(msg Message)
	Subscribe() <-chan Message
	Close()
}

type Message interface{}

func NewBroker(queueSize int) Broker {
	return &broker{queueSize: queueSize}
}

type broker struct {
	mu        sync.RWMutex
	queues    []queue
	queueSize int
}

func (b *broker) Publish(msg Message) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, q := range b.queues {
		q.push(msg)
	}
}

func (b *broker) Subscribe() <-chan Message {
	q := newQueue(b.queueSize)
	b.mu.Lock()
	b.queues = append(b.queues, q)
	b.mu.Unlock()
	return q.outCh
}

// Close must be called when all publishing is over to notify subscribers.
func (b *broker) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for _, q := range b.queues {
		q.stop()
	}
	b.queues = nil
}

// queue keeps track of messages that should be delivered to a subscriber. If buffer has reached
// max capacity and a new message is pushed the queue will save it but will remove the oldest
// message.
type queue struct {
	addCh  chan Message
	stopCh chan struct{}
	outCh  chan Message
}

func newQueue(size int) queue {
	q := queue{
		addCh:  make(chan Message, 1),
		stopCh: make(chan struct{}, 1),
		outCh:  make(chan Message, size),
	}
	go q.start()
	return q
}

func (q *queue) start() {
	for {
		select {
		case m := <-q.addCh:
			select {
			case q.outCh <- m:
				continue
			default:
			}
			// output queue is full, remove the oldest value add the new one
			select {
			case q.outCh <- m:
			case <-q.outCh:
				q.outCh <- m
			default:
			}

		case <-q.stopCh:
			// notify subscriber sending is over
			close(q.outCh)
			return
		}
	}
}

// push in not blocking
func (q *queue) push(m Message) {
	q.addCh <- m
}

func (q *queue) stop() {
	q.stopCh <- struct{}{}
}
