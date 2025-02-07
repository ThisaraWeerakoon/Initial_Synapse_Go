package models

type Configurations struct {
	Interval int //need to clarify the difference between interval and polling interval.For now consider Interval as polling interval and measured in seconds
	Sequential bool
	Coordination bool
	ActionAfterProcess string //MOVE, DELETE, NONE need to set default to NONE
	MoveAfterProcess string
	FileURI string
	MoveAfterFailure string
	FileNamePattern string
	ContentType string
	ActionAfterFailure string
}

//modify by adding a constructor