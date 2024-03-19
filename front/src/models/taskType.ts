import {DtosTaskType} from "../apiClient/generated";

export enum TaskType {
    GetWeather = "get_weather",
    GetPontChaban = "get_pont_chaban",
}

export function JobTypeToDTOType(type: TaskType): DtosTaskType {
    switch (type) {
        case TaskType.GetWeather:
            return DtosTaskType.JobTaskType_Weather;
        case TaskType.GetPontChaban:
            return DtosTaskType.JobTaskType_GetChabanDelmasBridgeStatus;
        default:
            return DtosTaskType.JobTaskType_Unknown;
    }
}