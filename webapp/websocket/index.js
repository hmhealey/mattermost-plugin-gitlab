import ActionTypes from '../action_types';
import Constants from '../constants';
import {
    getAssignedIssues,
    getAssignedMergeRequests,
    getCreatedMergeRequests,
    getConnected,
    // getUnreads,
} from '../actions';

export function handleConnect(store) {
    return (msg) => {
        console.log('connected to gitlab');
        if (!msg.data) {
            return;
        }

        store.dispatch({
            type: ActionTypes.RECEIVED_CONNECTED,
            data: {
                ...msg.data,
                settings: {sidebar_buttons: Constants.SETTING_BUTTONS_TEAM, daily_reminder: true}
            },
        });
    };
}

export function handleDisconnect(store) {
    return () => {
        console.log('disconnected from gitlab');
        store.dispatch({type: ActionTypes.RECEIVED_DISCONNECTED});
    };
}

export function handleReconnect(store, reminder = false) {
    return () => {
        console.log('reconnected to gitlab');
        store.dispatch(getConnected(reminder)).then(({data}) => {
            if (data && data.connected) {
                getMergeRequestsAndIssues(store);
            }
        });
    };
}

export function handleRefresh(store) {
    return () => {
        console.log('refreshing from gitlab');
        if (store.getState()['plugins-gitlab'].connected) {
            getMergeRequestsAndIssues(store);
        }
    };
}

function getMergeRequestsAndIssues(store) {
    store.dispatch(getCreatedMergeRequests());
    store.dispatch(getAssignedMergeRequests());
    store.dispatch(getAssignedIssues());
    // store.dispatch(getUnreads());
}