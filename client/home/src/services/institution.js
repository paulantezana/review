import request from '../utils/request';

const INSTITUTION_API = '/core/institution';

// Get all institution
export async function institutionPaginate(body) {
    return request(`${INSTITUTION_API}/paginate`, {
        method: 'POST',
        body,
    });
}

// Get ByID institution
export async function institutionById(body) {
    return request(`${INSTITUTION_API}/by/id`, {
        method: 'POST',
        body,
    });
}

// Create institution
export async function institutionCreate(body) {
    return request(`${INSTITUTION_API}/create`, {
        method: 'POST',
        body,
    });
}

// Update institution
export async function institutionUpdate(body) {
    return request(`${INSTITUTION_API}/update`, {
        method: 'PUT',
        body,
    });
}
