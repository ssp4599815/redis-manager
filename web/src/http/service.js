import axios from 'axios'
import config from './config'

const service = axios.create(config);

// 添加请求拦截器
service.interceptors.request.use(
    req => {
        console.log(req);
        return req
    },
    error => {
        return Promise.reject(error)
    }
);

// 添加响应拦截器（返回状态判断）
service.interceptors.response.use(
    res => {
        return res
    },
    error => {
        return Promise.reject(error.response.data) || {
            error: error.message
        }
    }
);

export default service