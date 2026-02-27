import axios from 'axios'

// 创建专门用于错误码管理的axios实例
const errCodeRequest = axios.create({
  baseURL: import.meta.env.VITE_ADMIN_API_URL  + '/payment-admin', // 使用ADMIN_API_URL，默认http://localhost:6081
  timeout: 15000 // 请求超时时间
})

// 响应拦截器配置（与默认request实例保持一致）
const errorHandler = (res) => {
  const { error_code, message } = res
  
  // 统一错误处理
  switch (error_code) {
    case 401:
      // 未登录或登录已过期
      console.error('未登录或登录已过期，请重新登录')
      break
    case 403:
      // 没有权限
      console.error('没有权限访问该资源')
      break
    case 500:
      // 服务器内部错误
      console.error('服务器内部错误')
      break
    default:
      // 其他错误
      console.error('请求错误', message || '未知错误')
  }
  
  return Promise.reject(res)
}

errCodeRequest.interceptors.response.use(
  response => {
    const res = response.data
    console.log('响应数据', res,res.error_code)
    // 检查响应状态码
    if (res.error_code === 10000) {
      console.log('成功响应', res)
      // 统一数据格式处理：将数据包装成标准格式
      return {
        code: res.error_code,
        data: res.data || [],
        msg: res.message || 'success',
        total: res.data ? res.data.length : 0 // 为分页组件提供总条数
      }
    } else {
      return errorHandler(res)
    }
  },
  error => {
    console.error('请求错误', error)
    // 处理网络错误等情况
    return Promise.reject(error)
  }
)

// 渠道错误码相关API
export const channelErrorApi = {
  // 获取渠道错误码列表
  getChannelErrorList: (params) => {
    return errCodeRequest({
      url: '/err-code/channel-error-list',
      method: 'get',
      params
    })
  },
  // 新增渠道错误码
  addChannelError: (data) => {
    return errCodeRequest({
      url: '/err-code/channel-error',
      method: 'post',
      data
    })
  },
  // 更新渠道错误码
  updateChannelError: (data) => {
    return errCodeRequest({
      url: '/err-code/channel-error',
      method: 'put',
      data
    })
  },
  // 删除渠道错误码
  deleteChannelError: (params) => {
    return errCodeRequest({
      url: '/err-code/channel-error',
      method: 'delete',
      params
    })
  }
}

// Jaco错误码相关API
export const jacoErrorApi = {
  // 获取Jaco错误码列表
  getJacoErrorList: () => {
    return errCodeRequest({
      url: '/err-code/jaco-error-list',
      method: 'get'
    })
  },
  // 新增Jaco错误码
  addJacoError: (data) => {
    return errCodeRequest({
      url: '/err-code/jaco-error',
      method: 'post',
      data
    })
  },
  // 更新Jaco错误码
  updateJacoError: (data) => {
    return errCodeRequest({
      url: '/err-code/jaco-error',
      method: 'put',
      data
    })
  },
  // 删除Jaco错误码
  deleteJacoError: (params) => {
    return errCodeRequest({
      url: '/err-code/jaco-error',
      method: 'delete',
      params
    })
  }
}

// 错误码映射相关API
export const errorMapApi = {
  // 获取错误码映射列表
  getErrorMapList: () => {
    return errCodeRequest({
      url: '/err-code/error-map-list',
      method: 'get'
    })
  },
  // 获取错误码映射详情
  getErrorMapDetails: (params) => {
    return errCodeRequest({
      url: '/err-code/error-map-details',
      method: 'get',
      params
    })
  },
  // 新增错误码映射
  addErrorMap: (data) => {
    return errCodeRequest({
      url: '/err-code/error-map',
      method: 'post',
      data
    })
  },
  // 更新错误码映射
  updateErrorMap: (data) => {
    return errCodeRequest({
      url: '/err-code/error-map',
      method: 'put',
      data
    })
  },
  // 删除错误码映射
  deleteErrorMap: (params) => {
    return errCodeRequest({
      url: '/err-code/error-map',
      method: 'delete',
      params
    })
  }
}

// 用户提示相关API
export const userPromptApi = {
  // 获取用户提示列表
  getUserPromptList: () => {
    return errCodeRequest({
      url: '/err-code/user-prompt-list',
      method: 'get'
    })
  },
  // 新增用户提示
  addUserPrompt: (data) => {
    return errCodeRequest({
      url: '/err-code/user-prompt',
      method: 'post',
      data
    })
  },
  // 更新用户提示
  updateUserPrompt: (data) => {
    return errCodeRequest({
      url: '/err-code/user-prompt',
      method: 'put',
      data
    })
  },
  // 删除用户提示
  deleteUserPrompt: (params) => {
    return errCodeRequest({
      url: '/err-code/user-prompt',
      method: 'delete',
      params
    })
  }
}
