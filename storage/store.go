package storage

// TaskAttributes is a struct which is used to transfer data from/to stores.
// All task data are converted from/to string to prevent the store from
// worrying about details of converting data to the proper formats.
type TaskAttributes struct {
	Hash        string  `json:"hash"`
	Name        string	`json:"name"`
	LastRun     string	`json:"lastRun"`
	NextRun     string	`json:"nextRun"`
	Duration    string	`json:"duration"`
	IsRecurring string	`json:"isRecurring"`
	Params      string	`json:"params"`
}

// TaskStore is the interface to implement when adding custom task storage.
type TaskStore interface {
	Add(TaskAttributes) error
	Fetch() ([]TaskAttributes, error)
	Remove(TaskAttributes) error
	Close() error
}
