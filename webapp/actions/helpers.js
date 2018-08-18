export function bindClientFunc(clientFunc, success) {
    return async (dispatch) => {
        let data;
        try {
            data = await clientFunc();
        } catch (error) {
            return {error};
        }

        let connected = dispatch(checkAndHandleNotConnected(data));
        if (!connected) {
            return {error: data};
        }

        dispatch({
            type: success,
            data,
        });

        return {data};
    };
}

function checkAndHandleNotConnected(data) {
    return (dispatch) => {
        if (!data || data.id !== 'not_connected') {
            return true;
        }

        dispatch({type: ActionTypes.RECEIVED_DISCONNECTED});

        return false;
    };
}