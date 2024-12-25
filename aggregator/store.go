package main

import (
	"fmt"

	"github.com/okusarobert/toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() Storage {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (s *MemoryStore) Insert(d types.Distance) error {
	s.data[d.OBUID] += d.Value
	return nil
}

func (s *MemoryStore) Get(id int) (float64, error) {
	dist, ok := s.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for obu id %d", id)
	}
	return dist, nil
}
