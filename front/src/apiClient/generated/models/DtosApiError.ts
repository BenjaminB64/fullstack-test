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

import { exists, mapValues } from '../runtime';
import type { ValidatorFieldError } from './ValidatorFieldError';
import {
    ValidatorFieldErrorFromJSON,
    ValidatorFieldErrorFromJSONTyped,
    ValidatorFieldErrorToJSON,
} from './ValidatorFieldError';

/**
 * 
 * @export
 * @interface DtosApiError
 */
export interface DtosApiError {
    /**
     * 
     * @type {string}
     * @memberof DtosApiError
     */
    error?: string;
    /**
     * 
     * @type {{ [key: string]: ValidatorFieldError; }}
     * @memberof DtosApiError
     */
    fields?: { [key: string]: ValidatorFieldError; };
}

/**
 * Check if a given object implements the DtosApiError interface.
 */
export function instanceOfDtosApiError(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function DtosApiErrorFromJSON(json: any): DtosApiError {
    return DtosApiErrorFromJSONTyped(json, false);
}

export function DtosApiErrorFromJSONTyped(json: any, ignoreDiscriminator: boolean): DtosApiError {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'error': !exists(json, 'error') ? undefined : json['error'],
        'fields': !exists(json, 'fields') ? undefined : (mapValues(json['fields'], ValidatorFieldErrorFromJSON)),
    };
}

export function DtosApiErrorToJSON(value?: DtosApiError | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'error': value.error,
        'fields': value.fields === undefined ? undefined : (mapValues(value.fields, ValidatorFieldErrorToJSON)),
    };
}
