package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Ratings map[int][]uuid.UUID

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
	ratings := make(map[int][]uuid.UUID)

	s.samples = append(s.samples, Sample{uuid.New(), name, description, filename, ratings, time.Now()})
	return nil
}

func (s *SampleService) GetSamples() []Sample {
	return s.samples
}

func (s *SampleService) GetSampleById(id uuid.UUID) (*Sample, error) {
	for _, sample := range s.samples {
		if sample.Id == id {
			return &sample, nil
		}
	}
	return &Sample{}, errors.New("sample not found")
}

func (s *SampleService) RateSample(id uuid.UUID, rate int, guest_user_id uuid.UUID) (*Sample, error) {
	sample, err := s.GetSampleById(id)
	if err != nil {
		return nil, err
	}

	for _, u := range sample.Ratings[rate] {
		if u == guest_user_id {
			return nil, errors.New("user vote already registered")
		}
	}

	sample.Ratings[rate] = append(sample.Ratings[rate], guest_user_id)
	return sample, nil
}
