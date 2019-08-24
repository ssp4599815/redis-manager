module.exports = {
    // 设置跨域
    devServer: {
        proxy: {
            '/api/v1': {
                target: 'https://www.easy-mock.com/mock/xxx',
                changeOrigin: true
            }
        }
    }
};
