package opphav

import "strings"

type Signal[S Sign, D Dimension] interface {
	Observable

	Sign() S
	Dimension() D
	IsZero() bool
}

type signal[S Sign, D Dimension] struct {
	sign      S
	dimension D
}

func (signal[S, D]) isObservable() {}

func (s signal[S, D]) Sign() S {
	return s.sign
}

func (s signal[S, D]) Dimension() D {
	return s.dimension
}

func (s signal[S, D]) IsZero() bool {
	return s.sign.IsZero() || s.dimension.IsZero()
}

func (s signal[S, D]) Observation() Observation {
	return NewObservation(s.sign.Name(), s.dimension.Name())
}

func NewSignal[S Sign, D Dimension](sign S, dimension D) Signal[S, D] {
	if sign.IsZero() || dimension.IsZero() {
		panic("opphav: unminted sign or dimension in signal")
	}

	signDomain, dimensionDomain := domainOf(sign.Name()), domainOf(dimension.Name())

	if signDomain != dimensionDomain {
		panic("opphav: signal crosses vocabularies: " + sign.Name() + " with " + dimension.Name())
	}

	return signal[S, D]{
		sign:      sign,
		dimension: dimension,
	}
}

func domainOf(qualified string) string {
	domain, _, _ := strings.Cut(qualified, ".")

	return domain
}
