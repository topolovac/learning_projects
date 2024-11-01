package services

import (
	"time"

	"github.com/google/uuid"
)

type Ratings map[int]int

type Sample struct {
	Id          uuid.UUID
	Name        string
	Description string
	Filename    string
	Ratings     Ratings
	Created     time.Time
}

type SampleService struct {
	samples []Sample
}

func (s *SampleService) CreateSample(name, description, filename string) error {
	ratings := make(map[int]int)
	ratings[1] = 0
	ratings[2] = 0
	ratings[3] = 0
	ratings[4] = 0
	ratings[5] = 0

	s.samples = append(s.samples, Sample{uuid.New(), name, description, filename, ratings, time.Now()})
	return nil
}

func (s *SampleService) GetSamples() []Sample {
	return s.samples
}
