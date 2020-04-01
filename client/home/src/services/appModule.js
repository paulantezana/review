import request from '../utils/request';

const MODULE_API = '/core/app/module';

// Get all module
export async function modulePaginate(body) {
    return request(`${MODULE_API}/paginate`, {
        method: 'POST',
        body,
    });
}

// Get ByID module
export async function moduleById(body) {
    return request(`${MODULE_API}/by/id`, {
        method: 'POST',
        body,
    });
}

// Create module
export async function moduleCreate(body) {
    return request(`${MODULE_API}/create`, {
        method: 'POST',
        body,
    });
}

// Update module
export async function moduleUpdate(body) {
    return request(`${MODULE_API}/update`, {
        method: 'PUT',
        body,
    });
}
