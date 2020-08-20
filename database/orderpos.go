package database

import "time"

type OrderPos struct {
	ONo      int
	OPos     int
	WPNo     int
	StepNo   int
	MainOPos int
	State    int
	OpNo     int
	WONo     int
	PNo      int

	subOrderBlocked bool
	Error           bool

	PlanedStart *time.Time
	PlanedEnd   *time.Time
	Start       *time.Time
	End         *time.Time

	ResourceID int
	Resource   Resource
	Carrier    Carrier
}
