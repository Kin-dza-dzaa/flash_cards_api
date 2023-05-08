package entity

import (
	"time"
)

type Collection struct {
	UserID     string
	Word       string
	Name       string
	LastRepeat time.Time
	// Duration which should be added to LastRepeat.
	// Each time TimeDiff should be incrementet: 2*TimeDiff + 1
	TimeDiff time.Duration
}
