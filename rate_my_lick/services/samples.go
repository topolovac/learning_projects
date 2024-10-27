package services

type Sample struct {
	Name        string
	Description string
	Filename    string
}

type SampleService struct {
	samples []Sample
}

func (s *SampleService) CreateSample(name, description, filename string) error {
	s.samples = append(s.samples, Sample{name, description, filename})
	return nil
}

func (s *SampleService) GetSamples() []Sample {
	return s.samples
}
