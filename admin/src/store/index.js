import { createPinia } from 'pinia'

// 创建Pinia实例
const pinia = createPinia()

// 可以在这里添加全局插件或其他配置

// 导出Pinia实例，供main.js使用
export default pinia

// 定义全局状态Store示例
import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    // 应用配置状态
    appConfig: {
      sidebarCollapse: false,
      currentTheme: 'light'
    },
    // 用户信息状态
    // userInfo: {
    //   username: '',
    //   role: '',
    //   permissions: []
    // },
    // // 加载状态
    // loading: {
    //   global: false,
    //   routes: [],
    //   components: {}
    // }
  }),
  
  getters: {
    // 获取侧边栏状态
    isSidebarCollapse: (state) => state.appConfig.sidebarCollapse,
    
    // 获取当前主题
    currentTheme: (state) => state.appConfig.currentTheme,
    
    // 检查用户是否有指定权限
    hasPermission: (state) => (permission) => {
      return state.userInfo.permissions.includes(permission)
    },
    
    // 检查全局加载状态
    isGlobalLoading: (state) => state.loading.global
  },
  
  actions: {
    // 切换侧边栏状态
    toggleSidebar() {
      this.appConfig.sidebarCollapse = !this.appConfig.sidebarCollapse
    },
    
    // 设置侧边栏状态
    setSidebarCollapse(collapse) {
      this.appConfig.sidebarCollapse = collapse
    },
    
    // 设置主题
    setTheme(theme) {
      this.appConfig.currentTheme = theme
      // 可以在这里添加主题切换的具体实现
    },
    
    // 设置用户信息
    setUserInfo(userInfo) {
      this.userInfo = { ...this.userInfo, ...userInfo }
    },
    
    // 清除用户信息
    clearUserInfo() {
      this.userInfo = {
        username: '',
        role: '',
        permissions: []
      }
    },
    
    // 设置全局加载状态
    setGlobalLoading(loading) {
      this.loading.global = loading
    },
    
    // 设置指定路由的加载状态
    setRouteLoading(routePath, loading) {
      if (loading) {
        if (!this.loading.routes.includes(routePath)) {
          this.loading.routes.push(routePath)
        }
      } else {
        this.loading.routes = this.loading.routes.filter(path => path !== routePath)
      }
    },
    
    // 设置组件加载状态
    setComponentLoading(componentName, loading) {
      this.loading.components[componentName] = loading
    },
    
    // 检查路由是否在加载中
    isRouteLoading(routePath) {
      return this.loading.routes.includes(routePath)
    },
    
    // 检查组件是否在加载中
    isComponentLoading(componentName) {
      return this.loading.components[componentName] || false
    }
  }
})

export const useTaskStore = defineStore('task', {
  state: () => ({
    // 任务相关状态
    tasks: [],
    currentTask: null,
    taskLoading: false,
    taskStats: {
      total: 0,
      active: 0,
      inactive: 0,
      executing: 0
    },
    // 任务日志相关状态
    logs: [],
    logsLoading: false,
    logStats: {
      todayTotal: 0,
      todaySuccess: 0,
      todayFailed: 0
    }
  }),
  
  getters: {
    // 获取所有任务
    allTasks: (state) => state.tasks,
    
    // 获取激活的任务
    activeTasks: (state) => state.tasks.filter(task => task.status === 1),
    
    // 获取停用的任务
    inactiveTasks: (state) => state.tasks.filter(task => task.status === 0),
    
    // 获取当前任务
    getCurrentTask: (state) => state.currentTask,
    
    // 获取任务数量统计
    taskStatistics: (state) => state.taskStats,
    
    // 获取日志列表
    allLogs: (state) => state.logs,
    
    // 获取日志统计
    logStatistics: (state) => state.logStats
  },
  
  actions: {
    // 设置任务列表
    setTasks(tasks) {
      this.tasks = tasks
      this.updateTaskStats()
    },
    
    // 添加新任务
    addTask(task) {
      this.tasks.push(task)
      this.updateTaskStats()
    },
    
    // 更新任务
    updateTask(updatedTask) {
      const index = this.tasks.findIndex(task => task.id === updatedTask.id)
      if (index !== -1) {
        this.tasks[index] = { ...this.tasks[index], ...updatedTask }
        this.updateTaskStats()
      }
    },
    
    // 删除任务
    deleteTask(taskId) {
      this.tasks = this.tasks.filter(task => task.id !== taskId)
      this.updateTaskStats()
    },
    
    // 设置当前任务
    setCurrentTask(task) {
      this.currentTask = task
    },
    
    // 设置任务加载状态
    setTaskLoading(loading) {
      this.taskLoading = loading
    },
    
    // 更新任务统计
    updateTaskStats() {
      this.taskStats = {
        total: this.tasks.length,
        active: this.tasks.filter(task => task.status === 1).length,
        inactive: this.tasks.filter(task => task.status === 0).length,
        executing: 0 // 实际项目中可能需要实时更新执行中的任务数量
      }
    },
    
    // 设置日志列表
    setLogs(logs) {
      this.logs = logs
    },
    
    // 添加新日志
    addLog(log) {
      this.logs.unshift(log) // 新日志添加到列表开头
      // 如果日志数量太多，可以考虑限制
      if (this.logs.length > 1000) {
        this.logs = this.logs.slice(0, 1000)
      }
    },
    
    // 设置日志加载状态
    setLogsLoading(loading) {
      this.logsLoading = loading
    },
    
    // 更新日志统计
    updateLogStats(stats) {
      this.logStats = { ...this.logStats, ...stats }
    },
    
    // 清空任务数据
    clearTasks() {
      this.tasks = []
      this.currentTask = null
      this.taskStats = {
        total: 0,
        active: 0,
        inactive: 0,
        executing: 0
      }
    },
    
    // 清空日志数据
    clearLogs() {
      this.logs = []
      this.logStats = {
        todayTotal: 0,
        todaySuccess: 0,
        todayFailed: 0
      }
    }
  }
})