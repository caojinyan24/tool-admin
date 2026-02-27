import request from '@/utils/request'

// 渠道配置相关API
export const channelMockChannelApi = {
  // 获取渠道配置列表
  getChannelConfigList: (params) => {
    return request({
      url: '/api/v1/channel-mock/channel-config',
      method: 'get',
      params
    })
  },
  // 新增渠道配置
  addChannelConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/channel-config',
      method: 'post',
      data
    })
  },
  // 更新渠道配置
  updateChannelConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/channel-config/update',
      method: 'post',
      data
    })
  },
  // 删除渠道配置
  deleteChannelConfig: (params) => {
    return request({
      url: '/api/v1/channel-mock/channel-config/delete',
      method: 'get',
      params
    })
  },
  // 根据ID获取渠道配置
  getChannelConfigByID: (params) => {
    return request({
      url: '/api/v1/channel-mock/channel-config/detail',
      method: 'get',
      params
    })
  }
}

// 响应配置相关API
export const channelMockResponseApi = {
  // 获取响应配置列表
  getResponseConfigList: (params) => {
    return request({
      url: '/api/v1/channel-mock/response-config',
      method: 'get',
      params
    })
  },
  // 新增响应配置
  addResponseConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/response-config',
      method: 'post',
      data
    })
  },
  // 更新响应配置
  updateResponseConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/response-config/update',
      method: 'post',
      data
    })
  },
  // 删除响应配置
  deleteResponseConfig: (params) => {
    return request({
      url: '/api/v1/channel-mock/response-config/delete',
      method: 'get',
      params
    })
  },
  // 根据条件获取响应配置列表（用于规则配置中的response_id选择器）
  getResponseConfigsByCondition: (params) => {
    return request({
      url: '/api/v1/channel-mock/response-config/condition',
      method: 'get',
      params
    })
  }
}

// 规则配置相关API
export const channelMockRuleApi = {
  // 获取规则配置列表
  getRuleConfigList: (params) => {
    return request({
      url: '/api/v1/channel-mock/rule-config',
      method: 'get',
      params
    })
  },
  // 新增规则配置
  addRuleConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/rule-config',
      method: 'post',
      data
    })
  },
  // 更新规则配置
  updateRuleConfig: (data) => {
    return request({
      url: '/api/v1/channel-mock/rule-config/update',
      method: 'post',
      data
    })
  },
  // 删除规则配置
  deleteRuleConfig: (params) => {
    return request({
      url: '/api/v1/channel-mock/rule-config/delete',
      method: 'get',
      params
    })
  },
  // 检查规则是否存在
  checkRuleExists: (params) => {
    return request({
      url: '/api/v1/channel-mock/rule-config/check',
      method: 'get',
      params
    })
  },
  // 新增规则（简化版，用于场景管理）
  addRule: (data) => {
    return request({
      url: '/api/v1/channel-mock/rule-config',
      method: 'post',
      data
    })
  },
  // 获取回调数据
  getCallbackData: (params) => {
    return request({
      url: '/api/v1/channel-mock/get-callback-data',
      method: 'get',
      params
    })
  },
  // 发送回调
  sendCallback: (params) => {
    return request({
      url: '/api/v1/channel-mock/send-callback',
      method: 'get',
      params
    })
  }
}

// 场景管理相关API
export const channelMockSceneApi = {
  // 获取场景列表
  getCustomSceneList: () => {
    return request({
      url: '/api/v1/channel-mock/custom-scene/list',
      method: 'get'
    })
  },
  // 新增场景
  addCustomScene: (data) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene/add',
      method: 'post',
      data
    })
  },
  // 删除场景
  deleteCustomScene: (params) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene/delete',
      method: 'get',
      params
    })
  },
  // 获取场景规则列表
  getCustomSceneRulesBySceneName: (params) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene-rule/list',
      method: 'get',
      params
    })
  },
  // 新增场景规则
  addCustomSceneRule: (data) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene-rule/add',
      method: 'post',
      data
    })
  },
  // 更新场景规则
  updateCustomSceneRule: (data) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene-rule/update',
      method: 'post',
      data
    })
  },
  // 删除场景规则
  deleteCustomSceneRule: (params) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene-rule/delete',
      method: 'get',
      params
    })
  },
  // 根据ID获取场景规则
  getCustomSceneRuleByID: (params) => {
    return request({
      url: '/api/v1/channel-mock/custom-scene-rule/detail',
      method: 'get',
      params
    })
  }
}