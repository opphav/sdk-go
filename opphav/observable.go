package opphav

import "reflect"

type Observable interface {
	Observation() Observation

	isObservable()
}

type Observation struct {
	sign      string
	dimension string
}

func NewObservation(sign, dimension string) Observation {
	return Observation{
		sign:      sign,
		dimension: dimension,
	}
}

func (obs Observation) Sign() string {
	return obs.sign
}

func (obs Observation) Dimension() string {
	return obs.dimension
}

func (obs Observation) IsZero() bool {
	return obs.sign == ""
}

func (e SignEntry) Observation() Observation {
	return Observation{sign: e.name}
}

func (SignEntry) isObservable() {}

func isNil(i any) bool {
	if i == nil {
		return true
	}

	val := reflect.ValueOf(i)
	switch val.Kind() {
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	default:
		return false
	}
}
