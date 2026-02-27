import axios from 'axios';

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL + '/api/v1',
  timeout: 10000
})

console.log('API Base URL:', import.meta.env.VITE_API_URL + '/api/v1')

// 请求拦截器
request.interceptors.request.use(
  config => {
    console.log('API请求发送:', config.method?.toUpperCase(), config.url, '参数:', config.params || config.data);
    return config;
  },
  error => {
    console.error('请求配置错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  response => {
    console.log('API响应接收:', response.config.url, '状态码:', response.status);
    const res = response.data;
    if (res.code !== 200) {
      console.error('API返回错误:', res.code, res.msg);
      return Promise.reject(new Error(res.msg || 'Error'));
    }
    console.log('API返回成功数据:', res.data);
    return res.data;
  },
  error => {
    console.error('API请求失败:', error.message);
    if (error.response) {
      console.error('响应错误详情:', {
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers
      });
    } else if (error.request) {
      console.error('请求已发送但未收到响应:', error.request);
    } else {
      console.error('请求配置错误:', error.message);
    }
    return Promise.reject(error);
  }
);

// 支付数据查询相关API
export const payDataApi = {
  /**
   * 查询交易列表
   * @param {Object} params 查询参数
   * @param {string} params.trade_id 交易ID
   * @param {string} params.out_trade_id 外部交易ID
   * @param {string} params.user_identity 用户标识
   * @param {string} params.trade_status 交易状态
   * @param {string} params.month 查询月份（格式：202511）
   * @param {number} params.page 页码
   * @param {number} params.page_size 每页大小
   */
  queryTradeList: (params) => {
    return request.get('/pay-data/trade-list', { params });
  },

  /**
   * 查询详情数据
   * @param {Object} params 查询参数
   * @param {string} params.trade_id 交易ID
   * @param {string} params.month 查询月份
   */
  queryDetail: (params) => {
    console.log('API queryDetail called with params:', params);
    console.log('完整请求URL:', request.defaults.baseURL + '/pay-data/detail');
    return request.get('/pay-data/detail', { params });
  },

  /**
   * 获取单条记录的JSON格式
   * @param {Object} params 查询参数
   * @param {string} params.db_type 数据库类型（trade, core, channel, notify）
   * @param {string} params.table_name 表名
   * @param {string} params.record_id 记录ID
   * @param {string} params.month 查询月份
   */
  getRecordJSON: (params) => {
    return request.get('/pay-data/record-json', { params });
  },

  // 查询渠道订单列表
  queryChannelOrderList: (params) => {
    return request.get('/pay-data/channel-order-list', { params });
  },

  // 查询渠道订单详情
  queryChannelOrderDetail: (params) => {
    return request.get('/pay-data/channel-order-detail', { params });
  },
};


export default payDataApi;