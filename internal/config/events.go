package config

type EventType string

var (
	EventBlockIndexed  EventType = "EventBlockIndexed"
	EventSignalIndexed EventType = "EventSignalIndexed"
)
