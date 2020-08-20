package database

import "time"

type FinOrder struct {
	ONo         int
	PlanedStart time.Time
	PlanedEnd   time.Time
	Start       time.Time
	End         time.Time
	CNo         int
	State       int
	Enabled     bool
	Release     time.Time
}
