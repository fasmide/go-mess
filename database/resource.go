package database

type Resource struct {
	ResourceID         int
	ResourceName       string
	Description        string
	PlcType            int
	IP                 string
	Picture            string
	ParallelProcessing bool
	Automatic          bool
	WebPage            string
	DefaultBrowser     bool
	TopologyType       int
}
