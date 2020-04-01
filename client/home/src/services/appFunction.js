import request from '../utils/request';

const MODULEFUNCTION_API = '/core/module/function';

// Get all moduleFunction
export async function moduleFunctionPaginate(body) {
    return request(`${MODULEFUNCTION_API}/paginate`, {
        method: 'POST',
        body,
    });
}

// Get ByID moduleFunction
export async function moduleFunctionById(body) {
    return request(`${MODULEFUNCTION_API}/by/id`, {
        method: 'POST',
        body,
    });
}

// Create moduleFunction
export async function moduleFunctionCreate(body) {
    return request(`${MODULEFUNCTION_API}/create`, {
        method: 'POST',
        body,
    });
}

// Update moduleFunction
export async function moduleFunctionUpdate(body) {
    return request(`${MODULEFUNCTION_API}/update`, {
        method: 'PUT',
        body,
    });
}
