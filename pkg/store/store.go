package store

import (
	"fmt"
	"sync"

	"github.com/jrhrmsll/tsgen/pkg/config"
	"github.com/jrhrmsll/tsgen/pkg/model"
	"github.com/jrhrmsll/tsgen/pkg/store/internal/index"
)

func key(s string, i int) string {
	return fmt.Sprintf("%s:%d", s, i)
}

type Store struct {
	paths model.Paths

	pathIndex  index.Index
	faultIndex index.Index

	mu sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		paths: model.Paths{},

		pathIndex:  index.NewIndex(),
		faultIndex: index.NewIndex(),
	}
}

func (s *Store) Init(cfg *config.Config) (*Store, error) {
	for _, cfgPath := range cfg.Paths {
		path, err := model.NewPath(cfgPath.Name, cfgPath.ResponseTime)
		if err != nil {
			return nil, err
		}

		err = s.InsertPath(path)
		if err != nil {
			return nil, err
		}

		for _, cfgFault := range cfgPath.Faults {
			fault, err := model.NewFault(path.Name, cfgFault.Code, cfgFault.Rate)
			if err != nil {
				return nil, err
			}

			err = s.InsertFault(fault)
			if err != nil {
				return nil, err
			}
		}
	}

	return s, nil
}

func (s *Store) Paths() model.Paths {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.paths
}

func (s *Store) Faults() model.Faults {
	s.mu.RLock()
	defer s.mu.RUnlock()

	faults := model.Faults{}
	for _, path := range s.paths {
		faults = append(faults, path.Faults...)
	}

	return faults
}

func (s *Store) InsertPath(path model.Path) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.pathIndex.Has(path.Name) {
		return fmt.Errorf("path '%s' already exist", path.Name)
	}

	s.paths = append(s.paths, path)

	s.pathIndex.Set(path.Name, len(s.paths)-1)

	return nil
}

func (s *Store) InsertFault(fault model.Fault) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.pathIndex.Has(fault.Path) {
		return fmt.Errorf("path '%s' not found", fault.Path)
	}

	faultKey := key(fault.Path, fault.Code)
	if s.faultIndex.Has(faultKey) {
		return fmt.Errorf("fault '%d' already exist for path '%s'", fault.Code, fault.Path)
	}

	pathIndex := s.pathIndex.Get(fault.Path)

	s.paths[pathIndex].Faults = append(s.paths[pathIndex].Faults, fault)
	s.faultIndex.Set(faultKey, len(s.paths[pathIndex].Faults)-1)

	return nil
}

func (s *Store) FindFaultBy(path string, code int) (model.Fault, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.pathIndex.Has(path) {
		return model.Fault{}, fmt.Errorf("path '%s' not found", path)
	}

	faultKey := key(path, code)
	if !s.faultIndex.Has(faultKey) {
		return model.Fault{}, fmt.Errorf("fault '%d' not found for path '%s'", code, path)
	}

	pathIndex := s.pathIndex.Get(path)
	faultIndex := s.faultIndex.Get(faultKey)

	return s.paths[pathIndex].Faults[faultIndex], nil
}

func (s *Store) UpdateFault(fault model.Fault) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !s.pathIndex.Has(fault.Path) {
		return fmt.Errorf("path '%s' not found", fault.Path)
	}

	faultKey := key(fault.Path, fault.Code)
	if !s.faultIndex.Has(faultKey) {
		return fmt.Errorf("fault '%d' not found for path '%s'", fault.Code, fault.Path)
	}

	pathIndex := s.pathIndex.Get(fault.Path)
	faultIndex := s.faultIndex.Get(faultKey)

	s.paths[pathIndex].Faults[faultIndex].Rate = fault.Rate

	return nil
}
