// KEYS
const TOKEN_KEY = '/&AV(%LK';
const USER_KEY = 'V?(%A¡T';

class NavigatorStorage {
    static getItem(key, remember = false){
        if(remember){
            if (typeof localStorage != "undefined") {
                return localStorage.getItem(key);
            } else {
                return undefined;
            }
        }else{
            if (typeof sessionStorage != "undefined") {
                return sessionStorage.getItem(key);
            } else {
                return undefined;
            }
        }
    }
    static setItem(key, value, remember = false){
        if(remember){
            if (typeof localStorage != "undefined") {
                return localStorage.setItem(key,value);
            } else {
                return undefined;
            }
        }else{
            if (typeof sessionStorage != "undefined") {
                return sessionStorage.setItem(key,value);
            } else {
                return undefined;
            }
        }
    }
    static clear(){
        if (typeof sessionStorage != "undefined") {
            sessionStorage.clear();
        }
        if (typeof localStorage != "undefined") {
            localStorage.clear();
        }
    }
}


// Recupera el token
export const getToken = () => {
    let token = NavigatorStorage.getItem(TOKEN_KEY,true);
    if (token === null) {
        token = NavigatorStorage.getItem(TOKEN_KEY,false);
    }
    return token;
};

// Recuperal el perfil
const getLicense = ()=> {
    let remember = true;
    let user = NavigatorStorage.getItem(USER_KEY,remember);
    if (user === null) {
        remember = false;
        user = NavigatorStorage.getItem(USER_KEY,remember);
    }
    user = JSON.parse(user);
    return { user, remember };
}

// Set authority
export const setAuthority = ({token, role_id = 0, remember = false, user = {}})=> {
    NavigatorStorage.setItem(TOKEN_KEY, token,remember);
    NavigatorStorage.setItem(USER_KEY, JSON.stringify(user),remember);
}


// Set authority new role
export const getAuthorityLicense = () => {
    const license = getLicense();
    if (license === null) return {};
    return license;
}

// Logout
export const destroy = () => {
    NavigatorStorage.clear();
};
