import request from '../utils/request';

const USER_API = '/core/user';
// User Login
export async function login(body) {
    return request(`/public${USER_API}/login`, {
        method: 'POST',
        body,
    });
}