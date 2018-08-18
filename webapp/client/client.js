import request from 'superagent';

export default class Client {
    constructor() {
        this.url = '/plugins/gitlab/api/v1';
    }

    getConnected = (reminder = false) => {
        return this.doGet(`${this.url}/connected?reminder=` + reminder);
    }

    // getReviews = () => {
    //     return this.doGet(`${this.url}/reviews`);
    // }

    getCreatedMergeRequests = () => {
        return this.doGet(`${this.url}/merge_requests/created`);
    }

    getAssignedMergeRequests = () => {
        return this.doGet(`${this.url}/merge_requests/assigned`);
    }

    getAssignedIssues = () => {
        return this.doGet(`${this.url}/issues/assigned`);
    }

    getTodos = () => {
        return this.doGet(`${this.url}/todos`);
    }

    // getMentions = () => {
    //     return this.doGet(`${this.url}/mentions`);
    // }

    // getUnreads = () => {
    //     return this.doGet(`${this.url}/unreads`);
    // }

    doGet = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        try {
            const response = await request.
                get(url).
                set(headers).
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }

    doPost = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        try {
            const response = await request.
                post(url).
                send(body).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }

    doDelete = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        try {
            const response = await request.
                delete(url).
                send(body).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }

    doPut = async (url, body, headers = {}) => {
        headers['X-Requested-With'] = 'XMLHttpRequest';
        headers['X-Timezone-Offset'] = new Date().getTimezoneOffset();

        try {
            const response = await request.
                put(url).
                send(body).
                set(headers).
                type('application/json').
                accept('application/json');

            return response.body;
        } catch (err) {
            throw err;
        }
    }
}
