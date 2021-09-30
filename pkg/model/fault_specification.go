package model

import (
	"fmt"
)

const (
	upperBound = 0.999999
)

type FaultSpecification struct {
	upperBound float32
}

func NewFaultSpecification() *FaultSpecification {
	return &FaultSpecification{
		upperBound: upperBound,
	}
}

func (spec *FaultSpecification) IsSatisfyBy(fault *Fault) (bool, error) {
	if fault.Rate > spec.upperBound {
		return false, fmt.Errorf("invalid fault rate '%.3f' for code '%d' and path '%s'", fault.Rate, fault.Code, fault.Path)
	}

	if fault.Code < 400 || fault.Code > 599 || fault.StatusText == "" {
		return false, fmt.Errorf("invalid fault code '%d' for path '%s'", fault.Code, fault.Path)
	}

	return true, nil
}
