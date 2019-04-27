package types

type RequestPriority struct {
	Name             string
	HandSize         uint
	Queues           []*Queue
	QueueLengthLimit uint
}
