package domain

type JobProcessor interface {
	ProcessJobs() error
}
