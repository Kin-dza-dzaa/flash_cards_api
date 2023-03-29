package entity

import (
	"time"
)

type Collection struct {
	UserID     string    `json:"-" validate:"required"`
	Word       string    `json:"word" validate:"required"`
	Name       string    `json:"collection_name" validate:"required"`
	LastRepeat time.Time `json:"last_repeat,omitempty"`
	// Duration which should be added to LastRepeat.
	// Each time TimeDiff should be incrementet: 2*TimeDiff + 1
	TimeDiff time.Duration `json:"time_diff,omitempty"`
}
