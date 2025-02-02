package models

type Configurations struct {
	interval int
	sequential bool
	coordination bool
	ActionAfterProcess string
	MoveAfterProcess string
	FileURI string
	MoveAfterFailure string
	FileNamePattern string
	ContentType string
	ActionAfterFailure string
}