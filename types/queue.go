package types

type Queue struct {
	Index uint
	Requests []*Request
	LengthLimit uint
}
