package database

import "time"

type Order struct {
	ONo         int
	PlanedStart *time.Time
	PlanedEnd   *time.Time
	Start       *time.Time
	End         *time.Time
	CNo         int
	StateID     int
	State       State
	Enabled     bool
	Release     *time.Time

	Positions []OrderPos
}
