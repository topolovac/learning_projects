package components

import (
	"fmt"

	"github.com/google/uuid"
)

func userAlreadyVoted(ratings []uuid.UUID, userId uuid.UUID) bool {
	for _, r := range ratings {
		if r == userId {
			return true
		}
	}
	return false
}
