module.exports = {
    siteMetadata: {
        title: `Sistema institucional para IEST públicas y privadas`,
        author: `PaulAntezana`,
        description: `Documentación de la propuesta de desarrollo de un sistema institucional para el instituto Vilcanota Sicuani.`,
        siteUrl: `https://iestpvilcanota.com`,
        // social: {
        //     facebook: `paul antezana`,
        // },
    },
    plugins: [
        `gatsby-plugin-react-helmet`,
        {
            resolve: 'gatsby-plugin-less',
            options: {
                javascriptEnabled: true,
                modifyVars: {
                    'primary-color': '#9A40D6',
                    'layout-header-height': '50px'
                }
            },
        },
        {
            resolve: 'gatsby-plugin-antd',
            options: {
                style: true,
            },
        },
        `gatsby-plugin-sass`,
        {
            resolve: `gatsby-source-filesystem`,
            options: {
                name: `src`,
                path: `${__dirname}/src/`,
            },
        },
        {
            resolve: `gatsby-transformer-remark`,
            options: {
              plugins: [
                `gatsby-remark-autolink-headers`,
                {
                    resolve: `gatsby-remark-prismjs`,
                    options: {
                        classPrefix: "language-",
                        inlineCodeMarker: null,
                        aliases: {},
                        showLineNumbers: true,
                        noInlineHighlight: false,
                    },
                },
                `gatsby-remark-external-links`,
                `gatsby-remark-responsive-iframe`,
              ],
            },
        },
    ]
}