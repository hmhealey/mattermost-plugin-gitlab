import {combineReducers} from 'redux';

import ActionTypes from '../action_types';
import Constants from '../constants';

function gitlabURL(state = '', action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.gitlab_url;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return '';
    default:
        return state;
    }
}

function connected(state = false, action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.connected;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return false;
    default:
        return state;
    }
}

function username(state = '', action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.gitlab_username;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return '';
    default:
        return state;
    }
}

function userId(state = '', action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.gitlab_user_id;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return '';
    default:
        return state;
    }
}

function settings(state = {sidebar_buttons: Constants.SETTING_BUTTONS_TEAM, daily_reminder: true, notifications: true}, action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.settings;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return {};
    default:
        return state;
    }
}

function clientId(state = '', action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CONNECTED:
        return action.data.gitlab_client_id;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return '';
    default:
        return state;
    }
}

function createdMergeRequests(state = [], action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_CREATED_MERGE_REQUESTS:
        return action.data;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return [];
    default:
        return state;
    }
}

function assignedMergeRequests(state = [], action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_ASSIGNED_MERGE_REQUESTS:
        return action.data;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return [];
    default:
        return state;
    }
}

function assignedIssues(state = [], action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_ASSIGNED_ISSUES:
        return action.data;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return [];
    default:
        return state;
    }
}

function mentions(state = [], action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_MENTIONS:
        return action.data;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return [];
    default:
        return state;
    }
}

function todos(state = [], action) {
    switch(action.type) {
        case ActionTypes.RECEIVED_TODOS:
            return action.data;
        case ActionTypes.RECEIVED_DISCONNECTED:
            return [];
        default:
            return state;
    }
}

function unreads(state = [], action) {
    switch(action.type) {
    case ActionTypes.RECEIVED_UNREADS:
        return action.data;
    case ActionTypes.RECEIVED_DISCONNECTED:
        return [];
    default:
        return state;
    }
}

export default combineReducers({
    gitlabURL,
    connected,
    username,
    userId,
    settings,
    clientId,
    createdMergeRequests,
    assignedMergeRequests,
    assignedIssues,
    mentions,
    unreads,
    todos,
});
