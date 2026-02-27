<template>
  <div class="scene-management">
    <!-- 场景列表区域 -->
    <div class="scene-list-section">
      <div class="section-header">
        <h3>场景列表</h3>
        <el-button type="success" @click="handleAddScene" size="small">新增场景</el-button>
      </div>
      
      <el-table :data="sceneList" style="width: 100%" v-loading="sceneLoading" size="small">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="scene_name" label="场景名称" min-width="200">
          <template v-slot:default="scope">
            <el-button type="text" @click="handleSceneClick(scope.row.scene_name)">{{ scope.row.scene_name }}</el-button>
          </template>
        </el-table-column>
      
        <el-table-column label="操作" width="180" fixed="right">
          <template v-slot:default="scope">
            <el-button type="success" size="small" @click="handleAddRuleToScene(scope.row.scene_name)">添加规则</el-button>
            <el-button type="danger" size="small" @click="handleDeleteScene(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 场景规则列表区域 -->
    <div class="scene-rule-section">
      <div class="section-header">
        <h3>场景规则列表</h3>
        <div class="header-actions">
          <el-select v-model="selectedSceneName" placeholder="请选择场景" style="width: 180px; margin-right: 10px;" size="small" @change="handleSceneChange">
            <el-option v-for="scene in sceneList" :key="scene.scene_name" :label="scene.scene_name" :value="scene.scene_name"></el-option>
          </el-select>
          <el-button type="success" @click="handleAddSceneRule" size="small" :disabled="!selectedSceneName">新增响应</el-button>
        </div>
      </div>
      
      <el-table :data="sceneRules" style="width: 100%" v-loading="ruleLoading" size="small">
        <el-table-column prop="id" label="ID" width="60"></el-table-column>
        <el-table-column prop="scene_name" label="场景名称" width="200"></el-table-column>
        <el-table-column prop="response_id" label="响应ID" min-width="180"></el-table-column>
        <!-- <el-table-column prop="created_at" label="创建时间" width="160"></el-table-column> -->
        <el-table-column label="操作" width="180" fixed="right">
          <template v-slot:default="scope">
            <el-button type="primary" size="small" @click="handleEditSceneRule(scope.row)">编辑</el-button>
            <el-button type="danger" size="small" @click="handleDeleteSceneRule(scope.row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 新增场景对话框 -->
    <el-dialog
      title="新增场景"
      v-model="sceneDialogVisible"
      width="500px"
    >
      <el-form ref="sceneFormRef" :model="sceneFormData" :rules="sceneFormRules" label-width="100px">
        <el-form-item label="场景说明" prop="scene_name">
          <el-input v-model="sceneFormData.scene_name" placeholder="请输入场景说明"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="sceneDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitScene">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新增/编辑场景规则对话框 -->
    <el-dialog
      :title="isEditRule ? '编辑场景规则' : '新增场景规则'"
      v-model="sceneRuleDialogVisible"
      width="500px"
    >
      <el-form ref="sceneRuleFormRef" :model="sceneRuleFormData" :rules="sceneRuleRules" label-width="100px">
        <el-form-item label="场景名称" prop="scene_name">
          <el-input v-model="sceneRuleFormData.scene_name" placeholder="场景名称" :disabled="isEditRule"></el-input>
        </el-form-item>
        <el-form-item label="渠道ID" prop="channel_id">
          <el-select v-model="sceneRuleFormData.channel_id" placeholder="请选择渠道ID" @change="handleChannelIdChange">
            <el-option v-for="option in channelIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型" prop="request_type">
          <el-select v-model="sceneRuleFormData.request_type" placeholder="请选择请求类型" @change="handleRequestTypeChange">
            <el-option label="create" value="create"></el-option>
            <el-option label="pay" value="pay"></el-option>
            <el-option label="query" value="query"></el-option>
            <el-option label="callback" value="callback"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="响应ID" prop="response_id">
          <el-select v-model="sceneRuleFormData.response_id" placeholder="请选择响应ID">
            <el-option v-for="option in responseIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="sceneRuleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitSceneRule">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 添加规则到场景对话框 -->
    <el-dialog
      title="添加规则到场景"
      v-model="addRuleDialogVisible"
      width="500px"
    >
      <el-form ref="addRuleFormRef" :model="addRuleFormData" :rules="addRuleFormRules" label-width="100px">
        <el-form-item label="用户标识" prop="user_identity">
          <el-input v-model="addRuleFormData.user_identity" placeholder="请输入用户标识"></el-input>
        </el-form-item>
        <el-form-item label="过期时间" prop="expire_time">
          <el-date-picker
            v-model="addRuleFormData.expire_time"
            type="datetime"
            placeholder="选择过期时间"
            style="width: 100%"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="规则状态" prop="status">
          <el-select v-model="addRuleFormData.status" placeholder="请选择规则状态">
            <el-option label="开启" value="on"></el-option>
            <el-option label="关闭" value="off"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="addRuleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitAddRule">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { channelMockSceneApi, channelMockResponseApi, channelMockRuleApi } from '@/api/channelMockApi'

const channelIdOptions = [
  { value: 'payermax_sab', label: 'payermax_sab' },
  { value: 'checkout-hk', label: 'checkout-hk' },
  { value: 'checkout-sab', label: 'checkout-sab' },
  { value: 'payermax-test-hk', label: 'payermax-test-hk' },
  { value: 'stcpay-test', label: 'stcpay-test' }
]

const sceneList = ref([])
const sceneRules = ref([])
const sceneLoading = ref(false)
const ruleLoading = ref(false)
const selectedSceneName = ref('')

const sceneDialogVisible = ref(false)
const sceneRuleDialogVisible = ref(false)
const addRuleDialogVisible = ref(false)
const isEditRule = ref(false)

const sceneFormRef = ref(null)
const sceneRuleFormRef = ref(null)
const addRuleFormRef = ref(null)

const sceneFormData = reactive({
  scene_name: ''
})

const sceneRuleFormData = reactive({
  id: null,
  scene_name: '',
  channel_id: '',
  request_type: '',
  response_id: ''
})

const addRuleFormData = reactive({
  user_identity: '',
  expire_time: '',
  status: 'on'
})

const sceneFormRules = {
  scene_name: [{ required: true, message: '请输入场景说明', trigger: 'blur' }]
}

const sceneRuleRules = {
  scene_name: [{ required: true, message: '场景名称不能为空', trigger: 'blur' }],
  channel_id: [{ required: true, message: '请选择渠道ID', trigger: 'change' }],
  request_type: [{ required: true, message: '请选择请求类型', trigger: 'change' }],
  response_id: [{ required: true, message: '请选择响应ID', trigger: 'change' }]
}

const addRuleFormRules = {
  user_identity: [{ required: true, message: '请输入用户标识', trigger: 'blur' }],
  expire_time: [{ required: true, message: '请选择过期时间', trigger: 'change' }],
  status: [{ required: true, message: '请选择规则状态', trigger: 'change' }]
}

const responseIdOptions = ref([])
const currentSceneName = ref('')

const fetchSceneList = async () => {
  sceneLoading.value = true
  try {
    const response = await channelMockSceneApi.getCustomSceneList()
    if (response.code === 200) {
      sceneList.value = response.data || []
    } else {
      ElMessage.error(response.msg || '获取场景列表失败')
    }
  } catch (error) {
    // 处理错误，显示后端返回的错误信息
    let errorMsg = '获取场景列表失败'
    if (error.msg) {
      errorMsg = error.msg
    } else if (error.response && error.response.data && error.response.data.msg) {
      errorMsg = error.response.data.msg
    }
    ElMessage.error(errorMsg)
    console.error('获取场景列表失败:', error)
  } finally {
    sceneLoading.value = false
  }
}

const fetchSceneRules = async (sceneName) => {
  ruleLoading.value = true
  try {
    const response = await channelMockSceneApi.getCustomSceneRulesBySceneName({
      scene_name: sceneName
    })
    if (response.code === 200) {
      sceneRules.value = response.data || []
    } else {
      ElMessage.error(response.msg || '获取场景规则列表失败')
    }
  } catch (error) {
    // 处理错误，显示后端返回的错误信息
    let errorMsg = '获取场景规则列表失败'
    if (error.msg) {
      errorMsg = error.msg
    } else if (error.response && error.response.data && error.response.data.msg) {
      errorMsg = error.response.data.msg
    }
    ElMessage.error(errorMsg)
    console.error('获取场景规则列表失败:', error)
  } finally {
    ruleLoading.value = false
  }
}

const handleSceneChange = (sceneName) => {
  if (sceneName) {
    fetchSceneRules(sceneName)
  } else {
    sceneRules.value = []
  }
}

const handleSceneClick = (sceneName) => {
  selectedSceneName.value = sceneName
  fetchSceneRules(sceneName)
}

const handleAddScene = () => {
  Object.assign(sceneFormData, {
    scene_name: ''
  })
  sceneDialogVisible.value = true
}

const handleSubmitScene = async () => {
  if (!sceneFormRef.value) return
  
  try {
    await sceneFormRef.value.validate()
    
    const sceneName = `${sceneFormData.scene_name}`
    
    const response = await channelMockSceneApi.addCustomScene({
      scene_name: sceneName
    })
    
    if (response.code === 200) {
      ElMessage.success('新增场景成功')
      sceneDialogVisible.value = false
      fetchSceneList()
    } else {
      ElMessage.error(response.msg || '新增场景失败')
    }
  } catch (error) {
    if (error !== 'cancel') {
      // 处理错误，显示后端返回的错误信息
      let errorMsg = '新增场景失败'
      if (error.msg) {
        errorMsg = error.msg
      } else if (error.response && error.response.data && error.response.data.msg) {
        errorMsg = error.response.data.msg
      }
      ElMessage.error(errorMsg)
      console.error('新增场景失败:', error)
    }
  }
}

const handleDeleteScene = (id) => {
  ElMessageBox.confirm('确定要删除这个场景吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await channelMockSceneApi.deleteCustomScene({ id })
      if (response.code === 200) {
        ElMessage.success('删除场景成功')
        fetchSceneList()
        if (selectedSceneName.value) {
          fetchSceneRules(selectedSceneName.value)
        }
      } else {
        ElMessage.error(response.msg || '删除场景失败')
      }
    } catch (error) {
      // 处理错误，显示后端返回的错误信息
      let errorMsg = '删除场景失败'
      if (error.msg) {
        errorMsg = error.msg
      } else if (error.response && error.response.data && error.response.data.msg) {
        errorMsg = error.response.data.msg
      }
      ElMessage.error(errorMsg)
      console.error('删除场景失败:', error)
    }
  }).catch(() => {})
}

const handleAddSceneRule = () => {
  Object.assign(sceneRuleFormData, {
    id: null,
    scene_name: selectedSceneName.value,
    channel_id: '',
    request_type: '',
    response_id: ''
  })
  isEditRule.value = false
  
  responseIdOptions.value = []
  
  sceneRuleDialogVisible.value = true
}

const handleEditSceneRule = (row) => {
  Object.assign(sceneRuleFormData, row)
  isEditRule.value = true
  
  // 加载响应ID选项
  if (row.channel_id && row.request_type) {
    fetchResponseIdOptions(row.channel_id, row.request_type)
  }
  
  sceneRuleDialogVisible.value = true
}

const fetchResponseIdOptions = async (channelId, requestType) => {
  try {
    const response = await channelMockResponseApi.getResponseConfigsByCondition({
      channel_id: channelId,
      request_type: requestType
    })
    if (response.code === 200) {
      responseIdOptions.value = response.data.map(item => ({
        value: item.response_id,
        label: item.response_id
      }))
    } else {
      responseIdOptions.value = []
    }
  } catch (error) {
    responseIdOptions.value = []
    // 处理错误，显示后端返回的错误信息
    let errorMsg = '获取响应ID选项失败'
    if (error.msg) {
      errorMsg = error.msg
    } else if (error.response && error.response.data && error.response.data.msg) {
      errorMsg = error.response.data.msg
    }
    ElMessage.error(errorMsg)
    console.error('获取响应ID选项失败:', error)
  }
}

const handleChannelIdChange = () => {
  // 重置response_id
  sceneRuleFormData.response_id = ''
  // 清空响应ID选项
  responseIdOptions.value = []
  // 当channel_id和request_type都存在时，加载响应ID选项
  if (sceneRuleFormData.channel_id && sceneRuleFormData.request_type) {
    fetchResponseIdOptions(sceneRuleFormData.channel_id, sceneRuleFormData.request_type)
  }
}

const handleRequestTypeChange = () => {
  // 重置response_id
  sceneRuleFormData.response_id = ''
  // 清空响应ID选项
  responseIdOptions.value = []
  // 当channel_id和request_type都存在时，加载响应ID选项
  if (sceneRuleFormData.channel_id && sceneRuleFormData.request_type) {
    fetchResponseIdOptions(sceneRuleFormData.channel_id, sceneRuleFormData.request_type)
  }
}

const handleSubmitSceneRule = async () => {
  if (!sceneRuleFormRef.value) return
  
  try {
    await sceneRuleFormRef.value.validate()
    
    const apiMethod = isEditRule.value ? channelMockSceneApi.updateCustomSceneRule : channelMockSceneApi.addCustomSceneRule
    
    const response = await apiMethod(sceneRuleFormData)
    if (response.code === 200) {
      ElMessage.success(isEditRule.value ? '更新场景规则成功' : '新增场景规则成功')
      sceneRuleDialogVisible.value = false
      if (selectedSceneName.value) {
        fetchSceneRules(selectedSceneName.value)
      }
    } else {
      ElMessage.error(response.msg || (isEditRule.value ? '更新场景规则失败' : '新增场景规则失败'))
    }
  } catch (error) {
    if (error !== 'cancel') {
      // 处理错误，显示后端返回的错误信息
      let errorMsg = '提交场景规则失败'
      if (error.msg) {
        errorMsg = error.msg
      } else if (error.response && error.response.data && error.response.data.msg) {
        errorMsg = error.response.data.msg
      }
      ElMessage.error(errorMsg)
      console.error('提交场景规则失败:', error)
    }
  }
}

const handleDeleteSceneRule = (id) => {
  ElMessageBox.confirm('确定要删除这条场景规则吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await channelMockSceneApi.deleteCustomSceneRule({ id })
      if (response.code === 200) {
        ElMessage.success('删除场景规则成功')
        if (selectedSceneName.value) {
          fetchSceneRules(selectedSceneName.value)
        }
      } else {
        ElMessage.error(response.msg || '删除场景规则失败')
      }
    } catch (error) {
      // 处理错误，显示后端返回的错误信息
      let errorMsg = '删除场景规则失败'
      if (error.msg) {
        errorMsg = error.msg
      } else if (error.response && error.response.data && error.response.data.msg) {
        errorMsg = error.response.data.msg
      }
      ElMessage.error(errorMsg)
      console.error('删除场景规则失败:', error)
    }
  }).catch(() => {})
}

// 生成默认过期时间（当前时间后一天）
const getDefaultExpireTime = () => {
  const now = new Date();
  const tomorrow = new Date(now);
  tomorrow.setDate(now.getDate() + 1);
  
  const year = tomorrow.getFullYear();
  const month = String(tomorrow.getMonth() + 1).padStart(2, '0');
  const day = String(tomorrow.getDate()).padStart(2, '0');
  const hours = String(tomorrow.getHours()).padStart(2, '0');
  const minutes = String(tomorrow.getMinutes()).padStart(2, '0');
  const seconds = String(tomorrow.getSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

const handleAddRuleToScene = (sceneName) => {
  currentSceneName.value = sceneName
  Object.assign(addRuleFormData, {
    user_identity: '',
    expire_time: getDefaultExpireTime(),
    status: 'on'
  })
  addRuleDialogVisible.value = true
}

const handleSubmitAddRule = async () => {
  if (!addRuleFormRef.value) return
  
  try {
    await addRuleFormRef.value.validate()
    
    // 获取场景关联的所有响应ID
    const sceneRuleResponse = await channelMockSceneApi.getCustomSceneRulesBySceneName({
      scene_name: currentSceneName.value
    })
    
    if (sceneRuleResponse.code !== 200 || !sceneRuleResponse.data || sceneRuleResponse.data.length === 0) {
      ElMessage.error('该场景下暂无场景规则，请先添加场景规则')
      return
    }
    
    // 遍历场景规则，为每个response_id创建规则
    for (const rule of sceneRuleResponse.data) {
      const responseId = rule.response_id
      
      try {
        // 检查是否已存在相同的user_identity和response_id
        const checkResponse = await channelMockRuleApi.checkRuleExists({
          user_identity: addRuleFormData.user_identity,
          response_id: responseId
        })
        
        if (checkResponse.code === 200 && checkResponse.data.exists) {
          ElMessage.warning(`用户标识 ${addRuleFormData.user_identity} 与响应ID ${responseId} 的规则已存在`)
          continue
        }
        
        // 转换过期时间格式为ISO 8601格式
        let expireTime = addRuleFormData.expire_time;
        if (expireTime && typeof expireTime === 'string' && !expireTime.includes('T')) {
          expireTime = expireTime.replace(' ', 'T') + 'Z';
        }
        
        // 创建规则
        const createResponse = await channelMockRuleApi.addRule({
          user_identity: addRuleFormData.user_identity,
          response_id: responseId,
          status: addRuleFormData.status,
          expire_time: expireTime
        })
        
        if (createResponse.code !== 200) {
          ElMessage.error(`添加规则失败: ${createResponse.msg}`)
        }
      } catch (error) {
        // 处理错误，显示后端返回的错误信息
        let errorMsg = `添加规则失败`
        if (error.msg) {
          errorMsg = `添加规则失败: ${error.msg}`
        } else if (error.response && error.response.data && error.response.data.msg) {
          errorMsg = `添加规则失败: ${error.response.data.msg}`
        }
        ElMessage.error(errorMsg)
        console.error('添加规则失败:', error)
      }
    }
    
    ElMessage.success('规则添加成功')
    addRuleDialogVisible.value = false
  } catch (error) {
    if (error !== 'cancel') {
      // 处理错误，显示后端返回的错误信息
      let errorMsg = '添加规则失败'
      if (error.msg) {
        errorMsg = error.msg
      } else if (error.response && error.response.data && error.response.data.msg) {
        errorMsg = error.response.data.msg
      }
      ElMessage.error(errorMsg)
      console.error('添加规则失败:', error)
    }
  }
}

onMounted(() => {
  fetchSceneList()
})
</script>

<style scoped>
.scene-management {
  display: flex;
  gap: 20px;
  height: calc(100vh - 200px);
}

.scene-list-section {
  flex: 0.4;
  min-width: 300px;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.scene-rule-section {
  flex: 0.6;
  min-width: 0;
  background-color: #ffffff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.section-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
}

.dialog-footer {
  text-align: right;
}

@media screen and (max-width: 1200px) {
  .scene-management {
    flex-direction: column;
    height: auto;
  }
  
  .scene-list-section,
  .scene-rule-section {
    min-height: 500px;
  }
}
</style>