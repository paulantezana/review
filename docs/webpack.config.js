const HtmlWebpackPlugin = require('html-webpack-plugin');
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const path = require('path');
const fs = require('fs');

// Leyendo todos los archivos del directorio indicado
const walkSyncFiles = (dir, filelist = []) => {
    fs.readdirSync(dir).forEach(file => {
        filelist = fs.statSync(path.join(dir, file)).isDirectory()
            ? walkSyncFiles(path.join(dir, file), filelist, file)
            : filelist.concat([`${dir}/${file}`]);
    });
    return filelist;
};
const pageDirFiles = walkSyncFiles(path.resolve(__dirname,'src/pug/pages')); // Obteniendo todo los directorios

const isPugFile = pageDirFiles.filter(item=>item.split('.').pop() === "pug"); // Filtrando solo los archivos pug

// Quitando la extencion del nombre de los archivos
const pageNames = isPugFile.map(item=>{
    let currentDir = item.replace(path.resolve(__dirname,'src/pug/pages'),"");
    currentDir = currentDir.replace(/\\/g, "/");
    currentDir = currentDir.substring(1);
    return currentDir.split('.').slice(0,-1).join('.');
});

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

            {
                test: /\.(jpe?g|png|gif|svg|webp)$/i,
                use: [
                    {
                        loader: 'file-loader',
                        options: {
                            name: 'assets/[name].[ext]',
                        },
                    },
                    {
                        loader: 'image-webpack-loader',
                        options: {
                            bypassOnDebug: true, // webpack@1.x
                            disable: true, // webpack@2.x and newer
                        },
                    },
                ],
            },

            {
                test: /\.(ttf|eot|woff2?|mp4|mp3|txt|xml|pdf)$/i,
                use: [{
                    loader: 'file-loader',
                    options: {
                        name: 'assets/[name].[ext]',
                    },
                }],
            },
            
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