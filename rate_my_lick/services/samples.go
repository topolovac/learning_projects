package services

import (
	"errors"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Ratings map[int][]uuid.UUID

func (r Ratings) GetTotal() int {
	total := 0
	for i, rr := range r {
		total += len(rr) * i
	}
	return total
}

type Sample struct {
	Id          uuid.UUID
	Name        string
	Description string
	Filename    string
	Ratings     Ratings
	Created     time.Time
	UserId      uuid.UUID
}

type SampleService struct {
	samples []Sample
}

func (s *SampleService) CreateSample(name, description, filename string, userId uuid.UUID) (*Sample, error) {
	ratings := make(map[int][]uuid.UUID)

	sample := &Sample{uuid.New(), name, description, filename, ratings, time.Now(), userId}
	s.samples = append(s.samples, *sample)
	return sample, nil
}

func (s *SampleService) GetSamplesByRating() []Sample {
	sort.Slice(s.samples, func(i, j int) bool {
		totalA := s.samples[i].Ratings.GetTotal()
		totalB := s.samples[j].Ratings.GetTotal()
		return totalB < totalA
	})
	return s.samples
}

func (s *SampleService) GetSamplesOrderByLatest() []Sample {
	sort.Slice(s.samples, func(i, j int) bool {
		return s.samples[i].Created.After(s.samples[j].Created)
	})
	return s.samples
}

func (s *SampleService) GetUserSamples(userId uuid.UUID) (samples []Sample) {
	allSamples := s.GetSamplesOrderByLatest()
	for _, sample := range allSamples {
		if sample.UserId == userId {
			samples = append(samples, sample)
		}
	}
	return
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

	sample.Ratings = removeExistingRatings(sample.Ratings, guest_user_id)

	sample.Ratings[rate] = append(sample.Ratings[rate], guest_user_id)
	return sample, nil
}

func removeExistingRatings(ratings Ratings, id uuid.UUID) Ratings {
	for i := 1; i <= 5; i++ {
		for ii, u := range ratings[i] {
			if u == id {
				ratings[i] = append(ratings[i][:ii], ratings[i][ii+1:]...)
				return ratings
			}
		}
	}
	return ratings
}
