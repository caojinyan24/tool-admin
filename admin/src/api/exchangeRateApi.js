import axios from 'axios';

// 创建axios实例
const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL + '/api/v1',
  timeout: 10000
})

console.log('Exchange Rate API Base URL:', import.meta.env.VITE_API_URL + '/api/v1')

// 请求拦截器
request.interceptors.request.use(
  config => {
    console.log('Exchange Rate API请求发送:', config.method?.toUpperCase(), config.url, '参数:', config.params || config.data);
    return config;
  },
  error => {
    console.error('Exchange Rate API请求配置错误:', error);
    return Promise.reject(error);
  }
);

// 响应拦截器
request.interceptors.response.use(
  response => {
    console.log('Exchange Rate API响应接收:', response.config.url, '状态码:', response.status);
    const res = response.data;
    if (res.code !== 200) {
      console.error('Exchange Rate API返回错误:', res.code, res.msg);
      return Promise.reject(new Error(res.msg || 'Error'));
    }
    console.log('Exchange Rate API返回成功数据:', res.data);
    return res.data;
  },
  error => {
    console.error('Exchange Rate API请求失败:', error.message);
    if (error.response) {
      console.error('Exchange Rate API响应错误详情:', {
        status: error.response.status,
        data: error.response.data,
        headers: error.response.headers
      });
    } else if (error.request) {
      console.error('Exchange Rate API请求已发送但未收到响应:', error.request);
    } else {
      console.error('Exchange Rate API请求配置错误:', error.message);
    }
    return Promise.reject(error);
  }
);

// 汇率配置相关API
export const exchangeRateApi = {
  /**
   * 查询汇率配置列表
   * @param {Object} params 查询参数
   * @param {string} params.source 数据源
   * @param {string} params.from_currency 源货币
   * @param {string} params.to_currency 目标货币
   * @param {number} params.page 页码
   * @param {number} params.page_size 每页大小
   */
  queryExchangeRateList: (params) => {
    return request.get('/exchange/exchange-rate-list', { params });
  },

  /**
   * 新增汇率配置
   * @param {Object} data 配置数据
   * @param {string} data.source 数据源
   * @param {string} data.from_currency 源货币
   * @param {string} data.to_currency 目标货币
   * @param {number} data.exchange_rate 汇率
   */
  addExchangeRate: (data) => {
    return request.post('/exchange/exchange-rate', data);
  },

  /**
   * 修改汇率配置
   * @param {Object} data 修改数据
   * @param {string} data.source 数据源
   * @param {string} data.from_currency 源货币
   * @param {string} data.to_currency 目标货币
   * @param {number} data.exchange_rate 新汇率
   */
  updateExchangeRate: (data) => {
    return request.put('/exchange/exchange-rate', data);
  },

  /**
   * 删除汇率配置（软删除，设置status为inactive）
   * @param {Object} data 删除数据
   * @param {string} data.source 数据源
   * @param {string} data.from_currency 源货币
   * @param {string} data.to_currency 目标货币
   */
  deleteExchangeRate: (data) => {
    return request.delete('/exchange/exchange-rate', { data });
  }
};

export default exchangeRateApi;