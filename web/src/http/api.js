import http from '@/http/service'

// 获取 集群 信息
const getClusterNodesApi = params => {
    return http.post(`api/v1/nodes`,params)
};


export {
    getClusterNodesApi
}