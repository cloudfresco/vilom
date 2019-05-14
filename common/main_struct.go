package common

import (
	"time"
)

// StatusDates - Used for all structs
type StatusDates struct {
	Statusc      uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedDay   uint
	CreatedWeek  uint
	CreatedMonth uint
	CreatedYear  uint
	UpdatedDay   uint
	UpdatedWeek  uint
	UpdatedMonth uint
	UpdatedYear  uint
}
