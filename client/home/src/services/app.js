import request from '../utils/request';

const APP_API = '/core/app';


// Get ByID app
export async function appById(body) {
    return request(`${APP_API}/by/id`, {
        method: 'POST',
        body,
    });
}


// Update app
export async function appUpdate(body) {
    return request(`${APP_API}/update`, {
        method: 'PUT',
        body,
    });
}
