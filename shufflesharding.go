package shufflesharding

import (
	"fmt"
	"shufflesharding/types"
)

func dispatchRequest(request *types.Request, context *types.Context) (*types.Queue, error) {
	fs := request.FindFlowSchema(context)
	if fs == nil {
		return nil, fmt.Errorf("no flow schema found for request %+v", request)
	}

	// bind request to FlowSchema
	request.FlowSchema = fs
	request.FlowDistinguisher = fs.GenFlowDistinguisher(request)

	queue := findQueueForRequest(request)
	if queue == nil {
		return nil, fmt.Errorf("unable to find an available queue for request %+v", request)
	}
	queue.Requests = append(queue.Requests, request)
	return queue, nil
}

func findQueueForRequest(request *types.Request) *types.Queue {
	pairHash := request.GetFlowIdentifierPairHash()

	rp := request.FlowSchema.RequestPriority
	queues := make([]*types.Queue, len(rp.Queues))
	copy(queues, rp.Queues)

	indices := shuffleShard(pairHash, uint(len(queues)), rp.HandSize)

	var bestQueue *types.Queue
	shortestQueueLength := ^uint(0)

	for _, index := range indices {
		index := index % uint(len(queues))
		queue := queues[index]
		// rm used element from array
		queues = append(queues[:index], queues[index+1:]...)

		queueLength := uint(len(queue.Requests))
		if queueLength > rp.QueueLengthLimit {
			continue
		}

		if queueLength < shortestQueueLength {
			shortestQueueLength = queueLength
			bestQueue = queue
		}
	}

	if bestQueue == nil {
		return nil
	}

	return bestQueue
}

func shuffleShard(V uint64, N uint, H uint) []uint {
	findAContext := findAContext{
		closetsADelta: ^uint64(0),
	}
	a := findA([]uint{}, V, N, H, &findAContext)
	if a != nil {
		return *a
	} else {
		return *findAContext.closetsA
	}
}

type findAContext struct {
	closetsA      *[]uint
	closetsADelta uint64
}

// use DFS to find A, where V == sigma(A[i] * AiTimesFF(N,i))
// TODO: such A[i] may not exist.
//  I changed this algorithm so that when such A[i] does not exist, A[i] such that sigma(A[i] * AiTimesFF(N,i)) is the closest to V is returned instead
func findA(A []uint, remainingV uint64, N uint, H uint, context *findAContext) *[]uint {
	if uint(len(A)) == H {
		if remainingV == 0 {
			return &A
		} else {
			if context.closetsADelta > remainingV {
				context.closetsADelta = remainingV
				context.closetsA = &A
			}
			return nil
		}
	}

	index := uint(len(A))
	rangeMin := uint(0)
	rangeMax := N - index

	for AiCandidate := rangeMin; AiCandidate < rangeMax; AiCandidate++ {
		if contains(&A, AiCandidate) {
			continue
		}
		i := uint(len(A))
		result := findA(append(A, AiCandidate), remainingV-AiTimesFF(AiCandidate, N, i), N, H, context)
		if result != nil {
			return result
		}
	}
	return nil
}

func contains(array *[]uint, element uint) bool {
	for _, item := range *array {
		if item == element {
			return true
		}
	}
	return false
}

func AiTimesFF(Ai uint, N uint, M uint) uint64 {
	result := uint64(Ai)
	for i := M + 1; i <= N; i++ {
		result *= uint64(i)
	}
	return result
}
