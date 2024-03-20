import {Chip, Table, TableBody, TableCell, TableColumn, TableHeader, TableRow} from "@nextui-org/react";
import {DtosJobResponse, DtosStatus, DtosTaskType} from "../apiClient/generated";
import {useEffect, useMemo, useState} from "react";
import {getApiClient} from "../apiClient/apiClient.tsx";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCircleNotch} from "@fortawesome/free-solid-svg-icons";

const statusColorMap: Map<string,  "success" | "warning" | "secondary" | "default" | "primary" | "danger" | undefined> = new Map();
statusColorMap.set(DtosStatus.Status_Completed, "success");
statusColorMap.set(DtosStatus.Status_Failed, "danger");
statusColorMap.set(DtosStatus.Status_InProgress, "warning");
statusColorMap.set(DtosStatus.Status_Pending, "secondary");
statusColorMap.set(DtosStatus.Status_Unknown, "secondary");

const taskTypeMap: Map<string, string> = new Map();
taskTypeMap.set(DtosTaskType.TaskType_Weather, "Get Weather");
taskTypeMap.set(DtosTaskType.TaskType_GetChabanDelmasBridgeSchedule, "Get Pont Chaban Delmas schedule");

const dateFormater = new Intl.DateTimeFormat(["en-US", "fr-FR"], {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "numeric",
    minute: "numeric",
    second: "numeric",
});

function JobTable() {
    const [ status, setStatus ] = useState<"idle" | "loading" | "error" | "success">("idle");
    const [ jobs, setJobs ] = useState<DtosJobResponse[]>([]);
    const apiClient = useMemo(() => getApiClient(), []);
    const refreshJobs = async  () => {
        setStatus("loading");
        try {
            const jobs = await apiClient.getJobs();
            setJobs(jobs);
            setStatus("success");
        } catch (e) {
            console.error(e);
            setStatus("error");
        }
    }
    useEffect(() => {
        refreshJobs();
    }, []);

    return (
        <div>
            <div
                className={(status === "loading" ? "text-primary border-primary p-2 border-1 rounded-medium" : "hidden")}>
                <FontAwesomeIcon spin={true} icon={faCircleNotch} className={"mr-2"}/> Loading...
            </div>
            <div
                className={(status === "error" ? "text-danger border-danger p-2 border-1 rounded-medium" : "hidden")}>
                Error occurred while fetching jobs
            </div>
            <Table>
                <TableHeader>
                    <TableColumn>Name</TableColumn>
                    <TableColumn>Type</TableColumn>
                    <TableColumn>Status</TableColumn>
                    <TableColumn>Create Date</TableColumn>
                    <TableColumn>Slack Webhook</TableColumn>
                </TableHeader>
                <TableBody isLoading={status === "loading"} emptyContent={"No job to display"}>
                    {
                        (jobs || []).map((job, index) => {
                            return (
                                <TableRow key={index}>
                                    <TableCell>{job.name}</TableCell>
                                    <TableCell>{taskTypeMap.get(job.taskType!)}</TableCell>
                                    <TableCell>
                                        <Chip className="capitalize" color={statusColorMap.get(job.status!)} size="sm" variant="flat">
                                            {job.status}
                                        </Chip>
                                    </TableCell>
                                    <TableCell>{dateFormater.format(Date.parse(job.createdAt!))}</TableCell>
                                    <TableCell><div className="text-sm max-w-72 overflow-auto">{job.slackWebhook}</div></TableCell>
                                </TableRow>
                            )
                        })
                    }
                </TableBody>
            </Table>
        </div>
    );
}

export default JobTable;