package zeromq

import (
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer"
	"github.com/NavExplorer/navexplorer-indexer-go/internal/indexer/IndexOption"
	"github.com/getsentry/raven-go"
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
		raven.CaptureErrorAndWait(err, nil)
		log.WithError(err).Fatal("Failed to open new 0MQ socket")
	}
	defer subscriber.Close()

	if err := subscriber.Connect(s.address); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.WithError(err).Fatalf("Failed to connect to %s", s.address)
	}
	if err := subscriber.SetSubscribe("hashblock"); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.WithError(err).Fatal("Failed to subscribe to 0MQ")
	}

	log.Info("Waiting for ZMQ messages")
	for {
		msg, err := subscriber.Recv(0)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			log.WithError(err).Fatal("Failed to receive message")
			break
		}

		if msg == "hashblock" {
			log.Info("New Block found")
			if err := s.indexer.Index(IndexOption.SingleIndex); err != nil {
				if err.Error() != "-8: Block height out of range" {
					raven.CaptureError(err, nil)
					log.WithError(err).Fatal("Failed to index subscribed block")
				}
			}
		}
	}
}
