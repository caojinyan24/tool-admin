import axios from 'axios'

// 创建axios实例，只从环境变量中加载baseURL
const request = axios.create({
  baseURL: import.meta.env.VITE_API_URL ,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
request.interceptors.request.use(
  config => {
    // 这里可以添加token等认证信息
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  response => {
    const res = response.data
    if (res.code !== 200) {
      // 处理错误，返回包含msg和error的对象
      const errorMsg = res.msg || '请求失败'
      const errorDetail = res.error || ''
      console.error('请求失败:', errorMsg, errorDetail)
      return Promise.reject({
        msg: errorMsg,
        error: errorDetail
      })
    }
    return res
  },
  error => {
    console.error('网络错误:', error)
    return Promise.reject({
      msg: '网络请求失败',
      error: error.message || ''
    })
  }
)

// 任务相关API
const taskApi = {
  // 获取所有任务
  getAllTasks() {
    return request.get('/api/v1/tasks')
  },
  
  // 获取单个任务
  getTaskByTaskName(taskName) {
    return request.get(`/api/v1/tasks/${taskName}`)
  },
  
  // 创建任务
  createTask(task) {
    return request.post('/api/v1/tasks/add', task)
  },
  
  // 更新任务
  updateTask(task) {
    return request.put(`/api/v1/tasks/edit`, task)
  },
  
  // 删除任务
  deleteTask(taskName) {
    return request.delete(`/api/v1/tasks/${taskName}`)
  },
  
  // 手动触发任务
  triggerTask(taskName, taskParam) {
    return request.post(`/api/v1/tasks/trigger`, { task_name: taskName, task_param: taskParam })
  }
}

export default taskApi

// 导出配置好的axios实例，供其他API模块使用
export const requestInstance = request