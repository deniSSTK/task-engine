package entUtils

import (
	"time"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NewTime() time.Time {
	return time.Now()
}
