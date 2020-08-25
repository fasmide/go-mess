package database

type Part struct {
	PNo         int
	Description string
	Type        int
	WPNo        int
	Picture     *string
	BasePallet  int
	MrpType     int
	SafetyStock int
	LotSize     int
}
