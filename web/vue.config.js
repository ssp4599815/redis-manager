module.exports = {
    // 设置跨域
    devServer: {
        proxy: {
            '/api/v1': {
                target: 'http://127.0.0.1:8089/',
                changeOrigin: true
            }
        }
    }
};
