import {Configuration, DtosCreateJobRequest, DtosUpdateJobRequest, JobsApi} from "./generated";

let apiClient: ApiClient;
export function getApiClient() {
    if (!apiClient) {
        apiClient = new ApiClient();
    }
    return apiClient;
}

export class ApiClient {
    private ApiClient: JobsApi;

    constructor() {
        const baseUrl = import.meta.env.VITE_API_URL;
        const configuration = new Configuration({basePath: baseUrl});
        this.ApiClient = new JobsApi(configuration)
    }

    async getJobs() {
        return this.ApiClient.jobsGet();
    }

    async getJob(id: number) {
        return this.ApiClient.jobsIdGet({id});
    }

    async createJob(createJobRequest: DtosCreateJobRequest) {
        return this.ApiClient.jobsPost({createJobRequest});
    }

    async updateJob(id: number, updateJobRequest: DtosUpdateJobRequest) {
        return this.ApiClient.jobsIdPut({id, updateJobRequest});
    }

    async deleteJob(id: number) {
        return this.ApiClient.jobsIdDelete({id});
    }
}