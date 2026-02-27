<template>
  <div class="channel-mock-response-config">
    <div class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item label="渠道ID">
          <el-select v-model="searchForm.channel_id" placeholder="请选择渠道ID" style="width: 180px;">
            <el-option v-for="option in channelIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型">
          <el-select v-model="searchForm.request_type" placeholder="请选择请求类型" style="width: 120px;">
            <el-option label="create" value="create"></el-option>
            <el-option label="pay" value="pay"></el-option>
            <el-option label="query" value="query"></el-option>
            <el-option label="callback" value="callback"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchResponseList">查询</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button type="success" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="responseList" style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
            <el-table-column prop="response_id" label="响应ID" width="180"></el-table-column>
      <el-table-column prop="channel_id" label="渠道ID" width="120"></el-table-column>
      <el-table-column prop="request_type" label="请求类型" width="100"></el-table-column>
      <el-table-column prop="mocked_body" label="响应内容" >
        <template v-slot:default="scope">
          <div class="response-content">
            {{ scope.row.mocked_body }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="response_param" label="响应参数" width="180"></el-table-column>
      <el-table-column label="响应状态" width="150">
        <template v-slot:default="scope">
          <span :class="getStatusClass(scope.row.status)">{{ scope.row.status }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template v-slot:default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row.id)">删除</el-button>
          <el-button type="success" size="small" @click="handleAddRule(scope.row)">添加规则</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination-container">
      <el-pagination
        background
        :current-page="pagination.currentPage"
        :page-sizes="[10, 20, 50, 100]"
        :page-size="pagination.pageSize"
        layout="total, sizes, prev, pager, next, jumper"
        :total="pagination.total"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      ></el-pagination>
    </div>

    <!-- 添加规则配置对话框 -->
    <el-dialog
      title="添加规则配置"
      v-model="ruleDialogVisible"
      width="500px"
    >
      <el-form ref="ruleFormRef" :model="ruleFormData" :rules="ruleRules" label-width="100px">
        <el-form-item label="用户标识" prop="user_identity">
          <el-input v-model="ruleFormData.user_identity" placeholder="请输入用户标识"></el-input>
        </el-form-item>
        <el-form-item label="过期时间" prop="expire_time">
          <el-date-picker
            v-model="ruleFormData.expire_time"
            type="datetime"
            placeholder="请选择过期时间"
            style="width: 100%;"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="规则状态" prop="status">
          <el-select v-model="ruleFormData.status" placeholder="请选择规则状态">
            <el-option label="启用" value="on"></el-option>
            <el-option label="禁用" value="off"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="ruleDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmitRule">确定</el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      title="新增/编辑响应配置"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="渠道ID" prop="channel_id">
          <el-select v-model="formData.channel_id" placeholder="请选择渠道ID" :disabled="isEdit">
            <el-option v-for="option in channelIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型" prop="request_type">
          <el-select v-model="formData.request_type" placeholder="请选择请求类型" :disabled="isEdit">
            <el-option label="create" value="create"></el-option>
            <el-option label="pay" value="pay"></el-option>
            <el-option label="query" value="query"></el-option>
            <el-option label="callback" value="callback"></el-option>
          </el-select>
        </el-form-item>
         <el-form-item label="响应内容" prop="mocked_body">
          <el-input
            type="textarea"
            v-model="formData.mocked_body"
            placeholder="请输入响应内容"
            rows="8"
            style="font-family: monospace;"
          ></el-input>
        </el-form-item>
            <el-form-item label="响应状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择响应状态" :disabled="isEdit">
            <el-option label="PENDING" value="PENDING"></el-option>
            <el-option label="SUCCESS" value="SUCCESS"></el-option>
            <el-option label="FAILED" value="FAILED"></el-option>
            <el-option label="UNKNOWN" value="UNKNOWN"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="响应ID">
          <el-row :gutter="10">
            <el-col :span="12">
              <el-input v-model="formData.response_id_auto" placeholder="自动生成" disabled></el-input>
            </el-col>
            <el-col :span="12">
              <el-form-item prop="response_id_user" style="margin-bottom: 0;">
                <el-input v-model="formData.response_id_user" placeholder="请输入自定义部分" :disabled="isEdit"></el-input>
              </el-form-item>
            </el-col>
          </el-row>
        </el-form-item>
        
        <el-form-item label="响应参数" prop="response_param">
          <el-select v-model="formData.response_param" placeholder="请选择响应参数" multiple>
            <el-option v-for="option in responseParamOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
    
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { channelMockResponseApi, channelMockRuleApi } from '@/api/channelMockApi'

// 渠道ID选项列表
const channelIdOptions = [
  { value: 'payermax_sab', label: 'payermax_sab' },
  { value: 'checkout-hk', label: 'checkout-hk' },
  { value: 'checkout-sab', label: 'checkout-sab' },
  { value: 'payermax-test-hk', label: 'payermax-test-hk' },
  { value: 'stcpay-test', label: 'stcpay-test' }
]

// 响应参数选项列表
const responseParamOptions = [
  { value: 'order_id', label: 'order_id' },
  { value: 'channel_trade_id', label: 'channel_trade_id' }
]

// 搜索表单
const searchForm = reactive({
  channel_id: '',
  request_type: '',
  status: ''
})

// 响应配置列表
const responseList = ref([])

// 加载状态
const loading = ref(false)

// 分页信息
const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

// 对话框状态
const dialogVisible = ref(false)
const ruleDialogVisible = ref(false)

// 编辑模式标志
const isEdit = ref(false)

// 表单引用
const formRef = ref(null)
const ruleFormRef = ref(null)

// 表单数据
const formData = reactive({
  id: null,
  channel_id: '',
  request_type: '',
  response_id: '',
  // 响应ID的自动拼接部分（渠道ID_请求类型_响应状态）
  response_id_auto: '',
  // 响应ID的用户输入部分
  response_id_user: '',
  mocked_body: '',
  // 响应参数（数组形式，用于多选框）
  response_param: [],
  status: ''
})

// 添加规则配置的弹窗数据
const ruleFormData = reactive({
  user_identity: '',
  response_id: '',
  expire_time: '',
  status: 'on'
})

// 规则表单验证规则
const ruleRules = {
  user_identity: [{ required: true, message: '请输入用户标识', trigger: 'blur' }],
  expire_time: [{ required: true, message: '请选择过期时间', trigger: 'change' }],
  status: [{ required: true, message: '请选择规则状态', trigger: 'change' }]
}

// 表单验证规则
const rules = {
  channel_id: [{ required: true, message: '请输入渠道ID', trigger: 'blur' }],
  request_type: [{ required: true, message: '请选择请求类型', trigger: 'change' }],
  response_id_user: [{ required: true, message: '请输入响应ID的自定义部分', trigger: 'blur' }],
  mocked_body: [{ required: true, message: '请输入响应内容', trigger: 'blur' }],
  response_param: [
    {
      validator: (rule, value, callback) => {
        if (formData.request_type === 'callback' && !value) {
          callback(new Error('请求类型为callback时, 响应参数必填'));
        } else {
          callback();
        }
      },
      trigger: ['blur', 'change']
    }
  ],
  status: [{ required: true, message: '请选择响应状态', trigger: 'change' }]
}

// 根据状态获取样式类名
const getStatusClass = (status) => {
  if (!status) return '';
  
  const statusMap = {
    'SUCCESS': 'status-success',
    'FAILED': 'status-failed',
    'PENDING': 'status-pending'
  };
  
  return statusMap[status] || '';
};

// 防抖函数
const debounce = (fn, delay) => {
  let timer = null;
  return function() {
    const context = this;
    const args = arguments;
    clearTimeout(timer);
    timer = setTimeout(() => {
      fn.apply(context, args);
    }, delay);
  };
};

// 更新自动拼接的响应ID部分（带防抖）
const updateResponseIdAuto = debounce(() => {
  const { channel_id, request_type, status } = formData;
  if (channel_id && request_type && status) {
    formData.response_id_auto = `${channel_id}_${request_type}_${status}`;
  } else {
    formData.response_id_auto = '';
  }
}, 300);

// 监听渠道ID变化，更新自动拼接的响应ID部分
watch(() => formData.channel_id, () => {
  updateResponseIdAuto();
});

// 监听请求类型变化，更新自动拼接的响应ID部分
watch(() => formData.request_type, (newVal) => {
  updateResponseIdAuto();
  // 当变为callback时触发响应参数验证
  if (newVal === 'callback' && formRef.value) {
    formRef.value.validateField('response_param');
  }
});

// 监听响应状态变化，更新自动拼接的响应ID部分
watch(() => formData.status, () => {
  updateResponseIdAuto();
});

// 获取响应配置列表
const fetchResponseList = async () => {
  loading.value = true;
  try {
    const response = await channelMockResponseApi.getResponseConfigList(searchForm);
    if (response.code === 200) {
      responseList.value = response.data || [];
      pagination.total = response.total || 0;
    } else {
      ElMessage.error(response.msg || '获取响应配置列表失败');
    }
  } catch (error) {
    ElMessage.error('获取响应配置列表失败');
    console.error('获取响应配置列表失败:', error);
  } finally {
    loading.value = false;
  }
};

// 重置表单
const resetForm = () => {
  searchForm.channel_id = '';
  searchForm.request_type = '';
  fetchResponseList();
};

// 处理分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  pagination.currentPage = 1;
  fetchResponseList();
};

// 处理当前页码变化
const handleCurrentChange = (currentPage) => {
  pagination.currentPage = currentPage;
  fetchResponseList();
};

// 处理编辑
const handleEdit = (row) => {
  // 复制行数据
  Object.assign(formData, row);
  
  // 在编辑模式下，拆分响应ID为自动部分和用户部分
  if (row.response_id) {
    // 尝试拆分响应ID
    const parts = row.response_id.split('_');
    if (parts.length >= 4) {
      // 自动部分：渠道ID_请求类型_响应状态
      formData.response_id_auto = `${parts[0]}_${parts[1]}_${parts[2]}`;
      // 用户部分：剩余的部分
      formData.response_id_user = parts.slice(3).join('_');
    } else {
      // 如果无法拆分，保持原样
      formData.response_id_auto = '';
      formData.response_id_user = row.response_id;
    }
  } else {
    formData.response_id_auto = '';
    formData.response_id_user = '';
  }
  
  // 在编辑模式下，拆分响应参数为数组
  if (row.response_param) {
    formData.response_param = row.response_param.split(',');
  } else {
    formData.response_param = [];
  }
  
  // 设置为编辑模式
  isEdit.value = true;
  
  // 打开对话框
  dialogVisible.value = true;
};

// 处理删除
const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这条响应配置吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await channelMockResponseApi.deleteResponseConfig({ id });
      if (response.code === 200) {
        ElMessage.success('删除成功');
        fetchResponseList();
      } else {
        ElMessage.error(response.msg || '删除失败');
      }
    } catch (error) {
      ElMessage.error('删除失败');
      console.error('删除失败:', error);
    }
  }).catch(() => {
    // 用户取消删除
  });
};

// 处理表单提交
const handleSubmit = async () => {
  if (!formRef.value) return;
  
  try {
    await formRef.value.validate();
    
    const isEdit = !!formData.id;
    const apiMethod = isEdit ? channelMockResponseApi.updateResponseConfig : channelMockResponseApi.addResponseConfig;
    
    let submitData = { ...formData };
    
    // 拼接完整的响应ID
    if (submitData.response_id_auto && submitData.response_id_user) {
      submitData.response_id = `${submitData.response_id_auto}_${submitData.response_id_user}`;
    }
    
    // 将响应参数数组通过英文逗号拼接成字符串
    if (Array.isArray(submitData.response_param)) {
      submitData.response_param = submitData.response_param.join(',');
    }
    
    if (!isEdit) {
      // 新增时，从提交数据中移除id字段
      delete submitData.id;
    } else {
      // 编辑时，确保id是数字类型
      submitData.id = Number(formData.id);
    }
    
    const response = await apiMethod(submitData);
    if (response.code === 200) {
      ElMessage.success(isEdit ? '更新成功' : '新增成功');
      dialogVisible.value = false;
      fetchResponseList();
    } else {
      ElMessage.error(response.msg || (isEdit ? '更新失败' : '新增失败'));
    }
  } catch (error) {
    // 表单验证失败
    return false;
  }
};

// 处理新增
const handleAdd = () => {
  // 重置表单数据
  Object.assign(formData, {
    id: null,
    channel_id: '',
    request_type: '',
    response_id: '',
    // 响应ID的自动拼接部分（渠道ID_请求类型_响应状态）
    response_id_auto: '',
    // 响应ID的用户输入部分
    response_id_user: '',
    mocked_body: '',
    response_param: [],
    status: ''
  });
  
  // 设置为新增模式
  isEdit.value = false;
  
  // 打开对话框
  dialogVisible.value = true;
};

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

// 处理添加规则按钮点击事件
const handleAddRule = (row) => {
  // 重置规则表单数据
  Object.assign(ruleFormData, {
    user_identity: '',
    response_id: row.response_id,
    expire_time: getDefaultExpireTime(),
    status: 'on'
  });
  
  // 打开规则对话框
  ruleDialogVisible.value = true;
};

// 处理提交规则配置
const handleSubmitRule = async () => {
  if (!ruleFormRef.value) return;
  
  try {
    await ruleFormRef.value.validate();
    
    // 处理过期时间格式
    const submitData = { ...ruleFormData };
    if (submitData.expire_time) {
      if (typeof submitData.expire_time === 'string' && !submitData.expire_time.includes('T')) {
        submitData.expire_time = submitData.expire_time.replace(' ', 'T') + 'Z';
      }
    }
    
    // 调用API添加规则配置
    const response = await channelMockRuleApi.addRuleConfig(submitData);
    if (response.code === 200) {
      ElMessage.success('添加规则配置成功');
      ruleDialogVisible.value = false;
    } else {
      ElMessage.error(response.msg || '添加规则配置失败');
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('添加规则配置失败');
      console.error('添加规则配置失败:', error);
    }
  }
};

// 组件挂载时获取响应配置列表
onMounted(() => {
  fetchResponseList();
});
</script>

<style scoped>
.search-bar {
  margin-bottom: 20px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 4px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.channel-mock-response-config {
  padding: 20px;
}

/* 成功状态样式 */
:deep(.status-success) {
  color: #67c23a;
  background-color: #f0f9eb;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

/* 失败状态样式 */
:deep(.status-failed) {
  color: #f56c6c;
  background-color: #fff1f0;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

/*  pending状态样式 */
:deep(.status-pending) {
  color: #409eff;
  background-color: #ecf5ff;
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 500;
}

/* 响应内容样式 */
.response-content {
  font-family: monospace;
  font-size: 12px;
  line-height: 1.4;
  color: #303133;
  word-break: break-all;
  white-space: pre-wrap;
  max-height: 120px;
  overflow: auto;
  padding: 8px;
  background-color: #f5f7fa;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  box-shadow: inset 0 1px 2px rgba(0, 0, 0, 0.05);
}

/* 表格单元格样式 */
:deep(.el-table__cell) {
  padding: 8px 12px !important;
}

/* 表格行高 */
:deep(.el-table__row) {
  height: auto !important;
  min-height: 48px !important;
}
</style>