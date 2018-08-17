import ActionTypes from '../action_types';
import Constants from '../constants';
import {getConnected, getReviews, getUnreads, getYourPrs} from '../actions';

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
    }
}

export function handleDisconnect(store) {
    return () => {
        console.log('disconnected from gitlab');
        store.dispatch({
            type: ActionTypes.RECEIVED_CONNECTED,
            data: {
                connected: false,
                gitlab_url: '',
                gitlab_username: '',
                gitlab_client_id: '',
                settings: {},
            }
        });
    }
}

export function handleReconnect(store, reminder = false) {
    return async () => {
        console.log('reconnected to gitlab');
        const {data} = await store.dispatch(getConnected(reminder));

        if (data && data.connected) {
            // store.dispatch(getReviews());
            // store.dispatch(getUnreads());
            // store.dispatch(getYourPrs());
            // store.dispatch(getYourAssignments());
        }
    }
}

export function handleRefresh(store) {
    return () => {
        console.log('refreshing from gitlab');
        if (store.getState()['plugins-gitlab'].connected) {
            // store.dispatch(getReviews());
            // store.dispatch(getUnreads());
            // store.dispatch(getYourPrs());
            // store.dispatch(getYourAssignments());
        }
    }
}
