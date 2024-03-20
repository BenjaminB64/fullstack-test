package domainProtoConverters

import (
	"github.com/BenjaminB64/fullstack-test/back/common/domain"
	"github.com/BenjaminB64/fullstack-test/back/common/protobuf/jobs-proto"
)

func DomainJobToGRPCJob(job *domain.Job) *jobs_proto.Job {
	return &jobs_proto.Job{
		Id:     int32(job.ID),
		Type:   DomainTaskTypeToGRPCTaskType(job.TaskType),
		Status: DomainStatusToGRPCStatus(job.Status),
		Result: nil,
	}
}

func DomainTaskTypeToGRPCTaskType(taskType domain.JobTaskType) jobs_proto.JobType {
	switch taskType {
	case domain.JobTaskType_Weather:
		return jobs_proto.JobType_JOB_TYPE_GET_WEATHER
	case domain.JobTaskType_GetChabanDelmasBridgeStatus:
		return jobs_proto.JobType_JOB_TYPE_GET_CHABAN_DELMAS_BRIDGE_STATUS
	default:
		return jobs_proto.JobType_JOB_TYPE_UNKNOWN
	}
}

func DomainStatusToGRPCStatus(status domain.JobStatus) jobs_proto.JobStatus {
	switch status {
	case domain.JobStatus_Pending:
		return jobs_proto.JobStatus_JOB_STATUS_PENDING
	case domain.JobStatus_InProgress:
		return jobs_proto.JobStatus_JOB_STATUS_IN_PROGRESS
	case domain.JobStatus_Completed:
		return jobs_proto.JobStatus_JOB_STATUS_COMPLETED
	case domain.JobStatus_Failed:
		return jobs_proto.JobStatus_JOB_STATUS_FAILED
	default:
		return jobs_proto.JobStatus_JOB_STATUS_UNKNOWN
	}
}

func GrpcJobToDomainJob(job *jobs_proto.Job) *domain.Job {
	return &domain.Job{
		ID:       int(job.Id),
		Status:   GrpcStatusToDomainStatus(job.Status),
		TaskType: GrpcTypeToDomainType(job.Type),
	}
}

func GrpcTypeToDomainType(taskType jobs_proto.JobType) domain.JobTaskType {
	switch taskType {
	case jobs_proto.JobType_JOB_TYPE_GET_WEATHER:
		return domain.JobTaskType_Weather
	case jobs_proto.JobType_JOB_TYPE_GET_CHABAN_DELMAS_BRIDGE_STATUS:
		return domain.JobTaskType_GetChabanDelmasBridgeStatus
	default:
		return domain.JobTaskType_Unknown
	}
}

func GrpcStatusToDomainStatus(status jobs_proto.JobStatus) domain.JobStatus {
	switch status {
	case jobs_proto.JobStatus_JOB_STATUS_PENDING:
		return domain.JobStatus_Pending
	case jobs_proto.JobStatus_JOB_STATUS_IN_PROGRESS:
		return domain.JobStatus_InProgress
	case jobs_proto.JobStatus_JOB_STATUS_COMPLETED:
		return domain.JobStatus_Completed
	case jobs_proto.JobStatus_JOB_STATUS_FAILED:
		return domain.JobStatus_Failed
	default:
		return domain.JobStatus_Unknown
	}
}
