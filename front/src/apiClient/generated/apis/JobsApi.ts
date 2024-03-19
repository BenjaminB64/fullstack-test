/* tslint:disable */
/* eslint-disable */
/**
 * Jobs API Service
 * This service provides a RESTful API for managing jobs that can be executed asynchronously
 *
 * The version of the OpenAPI document: 1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import * as runtime from '../runtime';
import type {
  DtosApiError,
  DtosCreateJobRequest,
  DtosJobResponse,
  DtosUpdateJobRequest,
} from '../models';
import {
    DtosApiErrorFromJSON,
    DtosApiErrorToJSON,
    DtosCreateJobRequestFromJSON,
    DtosCreateJobRequestToJSON,
    DtosJobResponseFromJSON,
    DtosJobResponseToJSON,
    DtosUpdateJobRequestFromJSON,
    DtosUpdateJobRequestToJSON,
} from '../models';

export interface JobsIdDeleteRequest {
    id: number;
}

export interface JobsIdGetRequest {
    id: number;
}

export interface JobsIdPutRequest {
    id: number;
    updateJobRequest: DtosUpdateJobRequest;
}

export interface JobsPostRequest {
    createJobRequest: DtosCreateJobRequest;
}

/**
 * 
 */
export class JobsApi extends runtime.BaseAPI {

    /**
     * Get jobs
     * Get jobs
     */
    async jobsGetRaw(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<Array<DtosJobResponse>>> {
        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/jobs`,
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => jsonValue.map(DtosJobResponseFromJSON));
    }

    /**
     * Get jobs
     * Get jobs
     */
    async jobsGet(initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<Array<DtosJobResponse>> {
        const response = await this.jobsGetRaw(initOverrides);
        return await response.value();
    }

    /**
     * Delete a job by ID
     * Delete a job by ID
     */
    async jobsIdDeleteRaw(requestParameters: JobsIdDeleteRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling jobsIdDelete.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/jobs/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'DELETE',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * Delete a job by ID
     * Delete a job by ID
     */
    async jobsIdDelete(requestParameters: JobsIdDeleteRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.jobsIdDeleteRaw(requestParameters, initOverrides);
    }

    /**
     * Get a job by ID
     * Get a job by ID
     */
    async jobsIdGetRaw(requestParameters: JobsIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DtosJobResponse>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling jobsIdGet.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        const response = await this.request({
            path: `/jobs/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'GET',
            headers: headerParameters,
            query: queryParameters,
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DtosJobResponseFromJSON(jsonValue));
    }

    /**
     * Get a job by ID
     * Get a job by ID
     */
    async jobsIdGet(requestParameters: JobsIdGetRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DtosJobResponse> {
        const response = await this.jobsIdGetRaw(requestParameters, initOverrides);
        return await response.value();
    }

    /**
     * Update a job by ID
     * Update a job by ID
     */
    async jobsIdPutRaw(requestParameters: JobsIdPutRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<void>> {
        if (requestParameters.id === null || requestParameters.id === undefined) {
            throw new runtime.RequiredError('id','Required parameter requestParameters.id was null or undefined when calling jobsIdPut.');
        }

        if (requestParameters.updateJobRequest === null || requestParameters.updateJobRequest === undefined) {
            throw new runtime.RequiredError('updateJobRequest','Required parameter requestParameters.updateJobRequest was null or undefined when calling jobsIdPut.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/jobs/{id}`.replace(`{${"id"}}`, encodeURIComponent(String(requestParameters.id))),
            method: 'PUT',
            headers: headerParameters,
            query: queryParameters,
            body: DtosUpdateJobRequestToJSON(requestParameters.updateJobRequest),
        }, initOverrides);

        return new runtime.VoidApiResponse(response);
    }

    /**
     * Update a job by ID
     * Update a job by ID
     */
    async jobsIdPut(requestParameters: JobsIdPutRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<void> {
        await this.jobsIdPutRaw(requestParameters, initOverrides);
    }

    /**
     * Create a job
     * Create a job
     */
    async jobsPostRaw(requestParameters: JobsPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<runtime.ApiResponse<DtosJobResponse>> {
        if (requestParameters.createJobRequest === null || requestParameters.createJobRequest === undefined) {
            throw new runtime.RequiredError('createJobRequest','Required parameter requestParameters.createJobRequest was null or undefined when calling jobsPost.');
        }

        const queryParameters: any = {};

        const headerParameters: runtime.HTTPHeaders = {};

        headerParameters['Content-Type'] = 'application/json';

        const response = await this.request({
            path: `/jobs`,
            method: 'POST',
            headers: headerParameters,
            query: queryParameters,
            body: DtosCreateJobRequestToJSON(requestParameters.createJobRequest),
        }, initOverrides);

        return new runtime.JSONApiResponse(response, (jsonValue) => DtosJobResponseFromJSON(jsonValue));
    }

    /**
     * Create a job
     * Create a job
     */
    async jobsPost(requestParameters: JobsPostRequest, initOverrides?: RequestInit | runtime.InitOverrideFunction): Promise<DtosJobResponse> {
        const response = await this.jobsPostRaw(requestParameters, initOverrides);
        return await response.value();
    }

}
