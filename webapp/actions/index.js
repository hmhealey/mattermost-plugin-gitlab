import Client from '../client';
import ActionTypes from '../action_types';

import {bindClientFunc} from './helpers';

export function getConnected(reminder = false) {
    return async (dispatch, getState) => {
        let data;
        try {
            data = await Client.getConnected(reminder);
        } catch (error) {
            return {error};
        }

        dispatch({
            type: ActionTypes.RECEIVED_CONNECTED,
            data: data,
        });

        return {data};
    };
}

export function getCreatedMergeRequests() {
    return bindClientFunc(
        Client.getCreatedMergeRequests,
        ActionTypes.RECEIVED_CREATED_MERGE_REQUESTS,
    );
}

export function getAssignedMergeRequests() {
    return bindClientFunc(
        Client.getAssignedMergeRequests,
        ActionTypes.RECEIVED_ASSIGNED_MERGE_REQUESTS,
    );
}

export function getAssignedIssues() {
    return bindClientFunc(
        Client.getAssignedIssues,
        ActionTypes.RECEIVED_ASSIGNED_ISSUES,
    );
}

export function getMentions() {
    return bindClientFunc(
        Client.getMentions,
        ActionTypes.RECEIVED_MENTIONS,
    );
}

export function getTodos() {
    return bindClientFunc(
        Client.getTodos,
        ActionTypes.RECEIVED_TODOS,
    );
}

export function getUnreads() {
    return bindClientFunc(
        Client.getUnreads,
        ActionTypes.RECEIVED_UNREADS,
    );
}
