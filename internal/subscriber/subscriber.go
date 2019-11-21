package subscriber

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
	zmq "github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
)

type Subscriber struct {
	address string
	indexer *indexer.Indexer
}

func New(address string, indexer *indexer.Indexer) *Subscriber {
	return &Subscriber{address, indexer}
}

func (s *Subscriber) Subscribe() {
	subscriber, err := zmq.NewSocket(zmq.SUB)
	if err != nil {
		log.WithError(err).Fatal("Failed to connect to 0MQ")
	}
	defer subscriber.Close()

	if err := subscriber.Connect(s.address); err != nil {
		log.WithError(err).Fatal("Failed to connect to 0MQ")
	}
	if err := subscriber.SetSubscribe("hashblock"); err != nil {
		log.WithError(err).Fatal("Failed to subscribe to 0MQ")
	}

	log.Info("Waiting for ZMQ messages")
	for {
		msg, err := subscriber.Recv(0)
		if err != nil {
			log.WithError(err).Fatal("Failed to receive message")
			break
		}

		if msg == "hashblock" {
			log.Info("New Block found")
			if err := s.indexer.Index(indexer.SingleIndex); err != nil {
				if err.Error() != "-8: Block height out of range" {
					log.WithError(err).Fatal("Failed to index subscribed block")
				}
			}
		}
	}
}
