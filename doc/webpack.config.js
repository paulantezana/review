const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const path = require('path');
const fs = require('fs');

const pageDirFiles = fs.readdirSync(path.resolve(__dirname, 'src/pug/pages')); // Leyendo todos los archivos del directorio indicado
const pageNames = pageDirFiles.map(item=>item.split('.').slice(0,-1).join('.')) ; // quitando la extencion del nombre de los archivos

// Creando un plugin por cada archivo
const pagePlugin = pageNames.map(page=>{
    return new HtmlWebpackPlugin({
        template    : `./src/pug/pages/${page}.pug`,
        filename    : `${page}.html`,
        minify: {
            html5: true,
            collapseWhitespace: true,
            caseSensitive: true,
            removeComments: true
        },
        hash        : true,
    })
});

module.exports = {
    entry: './src/index.js',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'main.js'
    },
    devServer: {
        contentBase : path.join(__dirname, "dist"),
        compress    : true,
        port        : 3000,
        open        : true,
        stats       : 'errors-only',
    },
    module: {
        rules: [
            {
                test: /\.(css|scss)$/,
                use: [
                    'style-loader',
                    MiniCssExtractPlugin.loader,
                    {
                        loader: 'css-loader',
                        options: {
                            importLoaders: 1,
                            // url: false,
                            // minimize: true,
                            // sourceMap: true,
                            // modules: true,
                            // localIdentName: '[local]__[hash:base64:5]',
                        }   
                    },
                    {
                        loader: 'postcss-loader'
                    },
                    {
                        loader: 'sass-loader',
                        options: {
                            sourceMap: true
                        }
                    }
                ]
            },

            // {
            //     test: /\.(jpe?g|png|gif|svg|webp)$/i,
            //     use: [
            //         'file-loader',
            //     ]
            // },

            // {
            //     test: /\.(ttf|eot|woff2?|mp4|mp3|txt|xml|pdf)$/i,
            //     use: 'file-loader?name=assets/[name].[ext]'
            // },
            
            {
                test    : /\.pug$/,
                use     : ['html-loader','pug-html-loader']
            },

            { 
                test: /\.js$/, 
                exclude: /node_modules/, 
                loader: "babel-loader" 
            }
        ]
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "[name].css",
            chunkFilename: "[id].css"
        }),
    ].concat(pagePlugin)
};