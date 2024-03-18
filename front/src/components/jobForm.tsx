import {Input, Select, SelectItem} from "@nextui-org/react";
import {useEffect, useMemo, useState} from "react";
import slackWebhookValidation from "../utils/validateSlackWebhook.tsx";

interface CreateJobFormProps {
    onAddJob: () => void;

    apiErrors?: {
        error?: string;
        fields?: {
            name?: string;
            type?: string;
            slackWebhook?: string;
        }
    }
}


const JobForm = ({onAddJob, apiErrors}: CreateJobFormProps) => {

    const [ slackWebhookValue, setSlackWebhookValue ] = useState("");
    const [ slackWebhookError, setSlackWebhookError ] = useState("");

    const slackWebhookIsInvalid = useMemo(() => {
        if (slackWebhookValue === "") {
            return false;
        }
        const isValid = slackWebhookValidation(slackWebhookValue);
        if (!isValid) {
            setSlackWebhookError("invalid Slack webhook");
        } else {
            setSlackWebhookError("");
        }
        return !isValid;
    }, [slackWebhookValue]);
    useEffect(() => {
        if (apiErrors?.fields?.slackWebhook) {
            setSlackWebhookError(apiErrors.fields.slackWebhook);
        }
    }, [apiErrors?.fields?.slackWebhook]);

    return (
        <>
            <form onSubmit={onAddJob}>
                <Input label={"Job Name"} errorMessage={apiErrors?.fields?.name} className={"mb-2"}/>

                <Select label={"Type"} errorMessage={apiErrors?.fields?.type} className={"mb-2"}>
                    <SelectItem key={"get_weather"}>Get weather</SelectItem>
                    <SelectItem key={"pont_chaban"}>Get Pont Chaman Delmas closing times</SelectItem>
                </Select>

                <Input
                    label={"Slack channel webhook"}
                    errorMessage={slackWebhookError}
                    isInvalid={slackWebhookIsInvalid}
                    onChange={(e) => setSlackWebhookValue(e.target.value)}
                />
            </form>
        </>
    )
};

export default JobForm;