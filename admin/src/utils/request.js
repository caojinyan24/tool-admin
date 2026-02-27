import axios from 'axios'

// 创建axios实例
const service = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/payment-admin', // api 的 base_url
  timeout: 15000 // 请求超时时间
})



// response 拦截器
const errorHandler = (response) => {
  const { code, msg } = response.data
  
  // 统一错误处理
  switch (code) {
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
      console.error('请求错误', msg)
  }
  
  return Promise.reject(response.data)
}

service.interceptors.response.use(
  response => {
    const res = response.data
    
    // 检查响应状态码
    if (res.code === 200) {
      // 统一数据格式处理：将数据包装成标准格式
      return {
        code: res.code,
        data: res.data || [],
        msg: res.msg || 'success',
        total: res.data ? res.data.length : 0 // 为分页组件提供总条数
      }
    } else {
      return errorHandler(response)
    }
  },
  error => {
    // 处理HTTP状态码错误
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 400:
          console.error('400 Bad Request:', data)
          // 尝试获取服务器返回的具体错误信息
          if (data && data.msg) {
            console.error('错误信息:', data.msg)
          }
          break
        case 401:
          console.error('401 Unauthorized')
          break
        case 403:
          console.error('403 Forbidden')
          break
        case 404:
          console.error('404 Not Found')
          break
        case 500:
          console.error('500 Internal Server Error')
          break
        default:
          console.error(`HTTP Error ${status}`)
      }
    } else if (error.request) {
      console.error('请求发送失败，未收到响应')
    } else {
      console.error('请求配置错误:', error.message)
    }
    
    console.error('完整错误信息:', error)
    return Promise.reject(error)
  }
)

export default service