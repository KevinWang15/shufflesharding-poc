package shufflesharding

import (
	"fmt"
	"math/rand"
	"shufflesharding/types"
	"shufflesharding/utils"
	"testing"
	"time"
)

var testContext types.Context

func genFlowSchema(requestPriority *types.RequestPriority) *types.FlowSchema {
	return &types.FlowSchema{
		Name:             utils.RandomString(6),
		MatchingPriority: uint(rand.Intn(10000)),
		RequestPriority:  requestPriority,
	}
}

func genRequestPriority() *types.RequestPriority {
	var queueSize = 128
	var queueLengthLimit = uint((rand.Intn(3) + 1) * 100)

	var requestPriorityQueues []*types.Queue

	for i := 0; i < queueSize; i++ {
		queue := genQueue(uint(i), queueLengthLimit)
		testContext.Queues = append(testContext.Queues, queue)
		requestPriorityQueues = append(requestPriorityQueues, queue)
	}

	return &types.RequestPriority{
		Name:             utils.RandomString(6),
		HandSize:         6,
		Queues:           requestPriorityQueues,
		QueueLengthLimit: queueLengthLimit,
	}
}

func genQueue(index uint, queueLengthLimit uint) *types.Queue {
	return &types.Queue{
		Index:       index,
		LengthLimit: queueLengthLimit,
		Requests:    []*types.Request{},
	}
}

func genRequest() *types.Request {
	return &types.Request{
		Content: utils.RandomString(100),
	}
}

func TestShuffleSharding(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 10; i++ {
		testContext.RequestPriorities = append(testContext.RequestPriorities, genRequestPriority())
	}
	for i := 0; i < 1000; i++ {
		testContext.FlowSchemas = append(testContext.FlowSchemas,
			genFlowSchema(testContext.RequestPriorities[rand.Int()%len(testContext.RequestPriorities)]),
		)
	}
	for i := 0; i < 1000; i++ {
		request := genRequest()
		queue, err := dispatchRequest(request, &testContext)
		if err != nil {
			fmt.Printf("error dispatching request: %v\n", err)
		} else {
			fmt.Printf("dispatched request to queue %d of request priority %s\n", queue.Index, request.FlowSchema.RequestPriority.Name)
		}
	}
}
