package domain

type JobTaskType string

const (
	JobTaskType_Unknown                     JobTaskType = ""
	JobTaskType_Weather                     JobTaskType = "get_weather"
	JobTaskType_GetChabanDelmasBridgeStatus JobTaskType = "get_chaban_delmas_bridge_status"
)

func (j JobTaskType) IsValid() bool {
	switch j {
	case JobTaskType_Weather, JobTaskType_GetChabanDelmasBridgeStatus:
		return true
	}
	return false
}
