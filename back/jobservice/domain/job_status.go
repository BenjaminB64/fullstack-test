package domain

type JobStatus string

const (
	JobStatus_Unknown    JobStatus = ""
	JobStatus_Pending    JobStatus = "pending"
	JobStatus_InProgress JobStatus = "in_progress"
	JobStatus_Completed  JobStatus = "completed"
	JobStatus_Failed     JobStatus = "failed"
)

func (j JobStatus) IsValid() bool {
	switch j {
	case JobStatus_Pending, JobStatus_InProgress, JobStatus_Completed, JobStatus_Failed:
		return true
	}
	return false
}
