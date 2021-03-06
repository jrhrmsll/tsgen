package model

import (
	"net/http"
)

type Fault struct {
	Path       string  `json:"path"`
	Code       int     `json:"code"`
	StatusText string  `json:"status_text"`
	Rate       float32 `json:"rate"`
}

type Faults []Fault

func NewFault(path string, code int, rate float32) (Fault, error) {
	fault := Fault{
		Path:       path,
		Code:       code,
		StatusText: http.StatusText(code),
		Rate:       rate,
	}

	faultSpecification := NewFaultSpecification()
	if ok, err := faultSpecification.IsSatisfyBy(fault); !ok {
		return Fault{}, err
	}

	return fault, nil
}

type Path struct {
	Name   string `json:"name"`
	Faults Faults `json:"faults"`
}
type Paths []Path

func NewPath(name string) (Path, error) {
	return Path{
		Name: name,
	}, nil
}
