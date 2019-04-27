package types

type Context struct {
	FlowSchemas       []*FlowSchema
	RequestPriorities []*RequestPriority
	Queues            []*Queue
}
