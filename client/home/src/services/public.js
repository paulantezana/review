import request from '../utils/request';

const PUBLIC_API = '/public';

// User app
export async function getApp(body) {
    return request(`${PUBLIC_API}/app`, {
        method: 'POST',
        body,
    });
}

export async function getAppModule(body) {
    return request(`${PUBLIC_API}/app/module`, {
        method: 'POST',
        body,
    });
}