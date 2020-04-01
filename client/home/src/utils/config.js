export const app = {
    name: 'Administrativo',
    key: 'admin',
    description: 'Sistema institucional.',
    version: '0.1.5',
    uri: 'https://iestpvilcanota.com',
    facebook: 'https://www.facebook.com/Paulantezana-764145183607069/inbox',
    youtube: 'https://www.youtube.com/channel/UCwnGqfqlVjDxRZJ-pFjP2oQ?view_as=subscriber',
    twitter: 'https://twitter.com/paulantezana',

    author: 'PaulAntezana',
    authorUri: 'https://paulantezana.com',

    // Certification
    // name: 'Certificación',
    // key: 'certificate',
    // description: 'Sistema de emición de certificaciones modulares',
};

const rootData = {
    domain: 'localhost',
    protocol: 'http',
    socket: 'ws',
    port: '1323',
}

export const service = {
    domain: rootData.domain,
    protocol: rootData.protocol,
    port: rootData.port,
    path: `${rootData.protocol}://${rootData.domain}:${rootData.port}`,
    api_path: `${rootData.protocol}://${rootData.domain}:${rootData.port}/api/v1`,
    socket: `${rootData.socket}://${rootData.domain}:${rootData.port}/ws`,
};

// export const service = {
//     path: 'http://api.localhost:1323',
//     api_path: 'http://api.localhost:1323/api/v1',
//     socket: 'ws://api.localhost:1323/ws',
// };

// export const service = {
//     path: 'https://institutional-server.herokuapp.com',
//     api_path: 'https://institutional-server.herokuapp.com/api/v1',
//     socket: 'wss://institutional-server.herokuapp.com/ws',
// };

// export const service = {
//     path: 'https://api.iestpvilcanota.com',
//     api_path: 'https://api.iestpvilcanota.com/api/v1',
//     socket: 'wss://api.iestpvilcanota.com/ws',
// };
