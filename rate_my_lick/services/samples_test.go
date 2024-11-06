package services

import (
	"testing"

	"github.com/google/uuid"
)

func TestRemoveExistingSample(t *testing.T) {
	test_uuid := uuid.New()

	ratings := Ratings{
		1: {uuid.New(), test_uuid, uuid.New()},
		2: {uuid.New()},
		3: {uuid.New()},
		4: {},
		5: {},
	}

	ratings = removeExistingRatings(ratings, test_uuid)
	if ratings[1][1] == test_uuid {
		t.Errorf("uuid at index 1 (%s) should not be %s", ratings[1][1], test_uuid)
	}

	ratings = Ratings{
		1: {},
		2: {},
		3: {},
		4: {},
		5: {test_uuid},
	}

	ratings = removeExistingRatings(ratings, test_uuid)
	if len(ratings[5]) > 0 {
		t.Errorf("array length should be 0")
	}
}
