package events

import (
	log "github.com/sirupsen/logrus"
	"github.com/yaacov/observer/observer"
)

type EventType string

var (
	EventBlockIndexed EventType = "EventBlockIndexed"
)

type Event struct {
	key   string
	value string
}

type Events struct {
	observer observer.Observer
}

func New() *Events {
	o := observer.Observer{}
	if err := o.Open(); err != nil {
		log.WithError(err).Fatal("Could not create new Observer")
	}

	return &Events{observer: o}
}

func (e *Events) Fire(event EventType, value string) {
	e.observer.Emit(Event{
		key:   string(EventBlockIndexed),
		value: value,
	})

	log.WithFields(log.Fields{"event": event, "value": value}).Info("Event")
}
