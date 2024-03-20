package domain

type JobTaskType string

const (
	JobTaskType_Unknown                       JobTaskType = ""
	JobTaskType_Weather                       JobTaskType = "get_weather"
	JobTaskType_GetChabanDelmasBridgeSchedule JobTaskType = "get_chaban_delmas_bridge_schedule"
)

func (j JobTaskType) IsValid() bool {
	switch j {
	case JobTaskType_Weather, JobTaskType_GetChabanDelmasBridgeSchedule:
		return true
	}
	return false
}
