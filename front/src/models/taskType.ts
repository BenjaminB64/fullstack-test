import {DtosTaskType} from "../apiClient/generated";

export enum TaskType {
    GetWeather = "get_weather",
    GetPontChabanSchedule = "get_pont_chaban_schedule",
}

export function JobTypeToDTOType(type: TaskType): DtosTaskType {
    switch (type) {
        case TaskType.GetWeather:
            return DtosTaskType.TaskType_Weather;
        case TaskType.GetPontChabanSchedule:
            return DtosTaskType.TaskType_GetChabanDelmasBridgeSchedule;
        default:
            return DtosTaskType.TaskType_Unknown;
    }
}