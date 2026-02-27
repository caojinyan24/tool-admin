import { createRouter, createWebHistory } from 'vue-router'

// 基础路由定义
const routes = [
  {
    path: '/',
    redirect: '/scheduler'
  },
  {
    path: '/scheduler',
    name: 'Scheduler',
    component: () => import('@/views/TaskList.vue'),
    meta: {
      title: '调度任务管理'
    }
  },
  {
    path: '/tasks',
    name: 'TaskList',
    component: () => import('@/views/TaskList.vue')
  },
  {
    path: '/tasks/add',
    name: 'AddTask',
    component: () => import('@/views/TaskForm.vue')
  },
  {
    path: '/tasks/edit/:taskName',
    name: 'EditTask',
    component: () => import('@/views/TaskForm.vue'),
    props: true,
    meta: {
      title: '编辑任务'
    }
  },
  {
    path: '/tasks/history',
    name: 'TaskHistory',
    component: () => import('@/views/TaskHistory.vue'),
    meta: {
      title: '任务执行历史'
    }
  },
  {
    path: '/channel-mock-config',
    name: 'channel-mock-config',
    component: () => import('@/views/ChannelMockConfig.vue'),
    meta: {
      title: '渠道模拟响应配置'
    }
  },
  {
    path: '/channel-mock-response-config',
    name: 'channel-mock-response-config',
    component: () => import('@/views/channel-mock/ChannelMockResponseConfig.vue'),
    props: route => ({ filterConfigId: route.query.configId }),
    meta: {
      title: '渠道响应配置'
    }
  },
  {
    path: '/payment-tool',
    name: 'payment-tool',
    component: () => import('@/views/PaymentTool.vue'),
    meta: {
      title: '支付工具'
    }
  }
  ,
  {
    path: '/error-code-manager',
    name: 'ErrorCodeManager',
    component: () => import('@/views/err-code/ErrorCodeManager.vue'),
    meta: {
      title: '错误码管理平台'
    }
  },
  {
    path: '/file-processor',
    name: 'file-processor',
    component: () => import('@/views/file-processor/FileProcessor.vue'),
    meta: {
      title: '文件处理'
    }
  }
]

// 判断是否为测试环境
const isTestEnvironment = import.meta.env.VITE_API_URL != 'http://pay-scheduler.alphicloud.com'

// 添加支付数据路由 - 直接在routes数组中定义，确保直接URL访问支持
routes.unshift(
  {
    path: '/pay-data',
    name: 'pay-data',
    component: () => import('@/views/PayDataList.vue'),
    meta: {
      title: '支付数据列表'
    }
  },
  {
    path: '/pay-data/detail',
    name: 'pay-data-detail',
    component: () => import('@/views/PayDataDetail.vue'),
    meta: {
      title: '支付数据详情'
    }
  },
  {
    path: '/pay-data/channel-detail',
    name: 'pay-data-channel-detail',
    component: () => import('@/views/PayDataChannelDetail.vue'),
    meta: {
      title: '渠道订单详情'
    }
  }
  // ,
  // {
  //   path: '/exchange-rate-config',
  //   name: 'exchange-rate-config',
  //   component: () => import('@/views/ExchangeRateConfig.vue'),
  //   meta: {
  //     title: '汇率配置管理'
  //   }
  // }
)

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes: [...routes],
  // 添加滚动行为配置
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 支付数据路由已在routes数组中直接定义，确保直接URL访问时组件可用

// 仅在测试环境下添加支付工具和渠道模拟响应配置路由
if (isTestEnvironment) {
  import('@/views/PaymentTool.vue').then(module => {
    router.addRoute({
      path: '/payment-tool',
      name: 'payment-tool',
      component: module.default,
      meta: {
        title: '支付工具'
      }
    })
  })
  
  import('@/views/ChannelMockConfig.vue').then(module => {
    router.addRoute({
      path: '/channel-mock-config',
      name: 'channel-mock-config',
      component: module.default,
      meta: {
        title: '渠道模拟响应配置'
      }
    })
  })
  
  import('@/views/channel-mock/ChannelMockResponseConfig.vue').then(module => {
    router.addRoute({
      path: '/channel-mock-response-config',
      name: 'channel-mock-response-config',
      component: module.default,
      props: route => ({ filterConfigId: route.query.configId }),
      meta: {
        title: '渠道响应配置'
      }
    })
  })
}


// 全局前置守卫 - 处理页面标题和权限
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title
  }
  next()
})

export default router