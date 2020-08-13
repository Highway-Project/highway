package service

const (
	Available = iota
	UnAvailable
	Disrupted
)

type Backend struct {
	Name   string
	Addr   string
	Weight int8
	Status int
}
