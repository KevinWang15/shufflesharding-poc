package types

import "hash/crc64"

type Request struct {
	Content string

	FlowSchema        *FlowSchema
	FlowDistinguisher string
}

func (r *Request) GetFlowIdentifierPairHash() uint64 {
	crc64Table := crc64.MakeTable(crc64.ISO)
	return crc64.Checksum([]byte(r.FlowSchema.Name), crc64Table) * crc64.Checksum([]byte(r.FlowDistinguisher), crc64Table)
}

func (r *Request) FindFlowSchema(context *Context) *FlowSchema {
	var bestFlowSchema *FlowSchema
	maxMatchingPriority := ^uint(0)

	for _, schema := range context.FlowSchemas {
		if schema.MatchesRequest(r) {
			// a numerically lower number indicates a logically higher priority
			if schema.MatchingPriority < maxMatchingPriority {
				bestFlowSchema = schema
				maxMatchingPriority = schema.MatchingPriority
			}
		}
	}

	return bestFlowSchema
}
