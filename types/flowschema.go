package types

import (
	"math/rand"
	"shufflesharding/utils"
)

type FlowSchema struct {
	Name             string
	MatchingPriority uint
	RequestPriority  *RequestPriority
}

func (f *FlowSchema) MatchesRequest(request *Request) bool {
	return rand.Intn(100) == 0
}

func (f *FlowSchema) GenFlowDistinguisher(request *Request) string {
	return utils.RandomString(10)
}
