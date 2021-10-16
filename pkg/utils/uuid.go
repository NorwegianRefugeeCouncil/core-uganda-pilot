package utils

import "github.com/satori/go.uuid"

type UIDGenerator interface {
	GenUID() string
}

type uidGenerator struct {
}

func (uidGenerator) GenUID() string {
	return uuid.NewV4().String()
}

func NewUIDGenerator() UIDGenerator {
	return uidGenerator{}
}

type DelegateUidGenerator struct {
	gen func() string
}

func (g *DelegateUidGenerator) GenUID() string {
	return g.gen()
}

func (g *DelegateUidGenerator) SetGeneratorFn(fn func() string) {
	g.gen = fn
}

func NewDelegateUIDGenerator(fn func() string) *DelegateUidGenerator {
	return &DelegateUidGenerator{
		gen: fn,
	}
}
