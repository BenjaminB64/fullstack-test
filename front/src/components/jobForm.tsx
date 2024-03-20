import {Input, Select, SelectItem} from "@nextui-org/react";
import React, {forwardRef, useState} from "react";
import {TaskType} from "../models/taskType.ts";

interface CreateJobFormProps {
    handleSubmit: (
            name?: string,
            type?: TaskType,
            slackWebhook?: string,
    ) => void;
    errors?: {
        name?: string;
        type?: string;
        slackWebhook?: string;
    };
}

const JobForm = forwardRef<HTMLFormElement, CreateJobFormProps>(({handleSubmit, errors}: CreateJobFormProps, ref ) => {
    const [ jobName , setJobName ] = useState("");
    const [ jobType, setJobType ] = useState<TaskType|undefined>();
    const [ slackWebhookValue, setSlackWebhookValue ] = useState("");

    const onSubmit = (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        handleSubmit(
            jobName,
            jobType,
            slackWebhookValue,
        );
    };

    return (
        <form onSubmit={onSubmit} ref={ref}>
            <Input
                label={"Job Name"}
                errorMessage={errors?.name}
                className={"mb-2"}
                onChange={(e) => setJobName(e.target.value)}
                isInvalid={!!errors?.name}
            />

            <Select
                label={"Type"}
                errorMessage={errors?.type}
                className={"mb-2"}
                onChange={(e) => setJobType(e.target.value as TaskType)}
                isInvalid={!!errors?.type}
            >
                <SelectItem key={TaskType.GetWeather}>Get weather</SelectItem>
                <SelectItem key={TaskType.GetPontChabanSchedule}>Get Pont Chaman Delmas schedule</SelectItem>
            </Select>

            <Input
                label={"Slack channel webhook"}
                errorMessage={errors?.slackWebhook}
                isInvalid={!!errors?.slackWebhook}
                onChange={(e) => setSlackWebhookValue(e.target.value)}
            />
        </form>
    )
});

export default JobForm;