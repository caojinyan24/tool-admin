import request from '@/utils/request'

/**
 * 查询短信验证码
 * @param {Object} params - 查询参数
 * @param {string} params.scene - 查询场景
 * @param {string} params.userName - 用户名
 * @returns {Promise}
 */
export function getPhoneCode(params) {
  return request({
    url: '/v1/payment-tool/phone-code',
    method: 'get',
    params
  })
}