import {forwardRef, useCallback, useMemo, useState} from "react";
import JobForm from "./jobForm.tsx";
import {JobTypeToDTOType, TaskType} from "../models/taskType.ts";
import {ApiClient} from "../apiClient/apiClient.tsx";
import {DtosApiErrorFromJSON, ResponseError} from "../apiClient/generated";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCircleNotch} from "@fortawesome/free-solid-svg-icons";

type JobFormContainerProps = unknown;

export const JobFormContainer = forwardRef<HTMLFormElement, JobFormContainerProps>((props, ref) => {
    console.log(props);
    const apiClient = useMemo(() => new ApiClient(), []);

    const [ error, setError ] = useState<string | null>(null);
    const [ status, setStatus ] = useState<"idle" | "loading" | "error" | "success">("idle");

    const [ errors, setErrors ] = useState<{
        name?: string;
        type?: string;
        slackWebhook?: string;
    }>({});

    const handleSubmit = useCallback(async (name?: string, taskType?: TaskType, slackWebhook?: string) => {
        setStatus("loading");
        setError(null);
        setErrors({});

        console.log("Name: ", name);
        console.log("Type: ", taskType);
        console.log("Slack webhook: ", slackWebhook);
        if (!name || !taskType) {
            setStatus("error");
            setError("Name and type are required");
            setErrors({
                name: name ? undefined : "Name is required",
                type: taskType ? undefined : "Type is required",
            });
            return;
        }

        try {
            setStatus("loading");
            await apiClient.createJob({
                name: name,
                taskType: JobTypeToDTOType(taskType),
                slackWebhook: slackWebhook,
            })
            setStatus("success");
        } catch (e) {
            setStatus("error");
            console.error(e);

            let error = "unknown error occurred. Please try again later.";
            if (e instanceof ResponseError) {
                const apiError = DtosApiErrorFromJSON(await e.response.json());
                if (apiError.fields) {
                    setErrors({
                        name: apiError.fields.name?.error,
                        type: apiError.fields.taskType?.error,
                        slackWebhook: apiError.fields.slackWebhook?.error,
                    });
                }
                if (apiError.error) {
                    error = apiError.error;
                }
            }

            setError(error);
        }

    }, [apiClient]);

    return (
        <>
            <div className={(status === "loading" ? "text-primary border-primary p-2 border-1 rounded-medium" : "hidden")}><FontAwesomeIcon spin={true} icon={faCircleNotch} className={"mr-2"} /> Loading...</div>
            <div className={(status === "success" ? "text-success border-success p-2 border-1 rounded-medium" : "hidden")}>Job created successfully</div>
            <div className={(status === "error" ? "text-danger border-danger p-2 border-1 rounded-medium" : "hidden")}>{error}</div>
            <JobForm ref={ref} errors={errors}  handleSubmit={handleSubmit} />
        </>
    );
});