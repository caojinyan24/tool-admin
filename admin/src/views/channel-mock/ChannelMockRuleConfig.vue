<template>
  <div class="channel-mock-rule-config">
    <div class="search-bar">
      <el-form :model="searchForm" inline>
        <el-form-item label="渠道ID">
          <el-select v-model="searchForm.channel_id" placeholder="请选择渠道ID" style="width: 150px;">
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
        <el-form-item label="用户标识">
          <el-input v-model="searchForm.user_identity" placeholder="请输入用户标识" style="width: 180px;"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchRuleList">查询</el-button>
          <el-button @click="resetForm">重置</el-button>
          <el-button type="success" @click="handleAdd">新增</el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-table :data="ruleList" style="width: 100%" v-loading="loading" :row-class-name="tableRowClassName">
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="channel_id" label="渠道ID" width="120"></el-table-column>
      <el-table-column prop="request_type" label="请求类型" width="100"></el-table-column>
      <el-table-column prop="user_identity" label="用户标识" width="180"></el-table-column>

      <!-- <el-table-column prop="response_id" label="响应ID" width="100"></el-table-column> -->
      <el-table-column prop="mocked_body" label="响应内容" min-width="500" >
        <template v-slot:default="scope">
          <div class="response-content">
            {{ scope.row.mocked_body }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="response_param" label="响应参数" width="180"></el-table-column>
      <el-table-column prop="order_status" label="响应状态" width="100">
<template v-slot:default="scope">
      <span :class="getStatusClass(scope.row.order_status)">{{ scope.row.order_status }}</span>
        </template>

      </el-table-column>
      <el-table-column prop="expire_time" label="过期时间" width="180" show-overflow-tooltip>
        <template v-slot:default="scope">
          {{ formatDate(scope.row.expire_time) }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="80">
        <template v-slot:default="scope">
          <el-tag :type="scope.row.status === 'on' ? 'success' : 'danger'">
            {{ scope.row.status === 'on' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="250" fixed="right">
        <template v-slot:default="scope">
          <el-button type="primary" size="small" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDelete(scope.row.id)">删除</el-button>
          <el-button v-if="scope.row.request_type === 'callback'" type="warning" size="small" @click="handleSendCallback(scope.row)">发送回调</el-button>
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

    <!-- 新增/编辑对话框 -->
    <el-dialog
      title="规则配置"
      v-model="dialogVisible"
      width="600px"
    >
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="渠道ID" prop="channel_id">
          <el-select v-model="formData.channel_id" placeholder="请选择渠道ID">
            <el-option v-for="option in channelIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="请求类型" prop="request_type">
          <el-select v-model="formData.request_type" placeholder="请选择请求类型">
            <el-option label="create" value="create"></el-option>
            <el-option label="pay" value="pay"></el-option>
            <el-option label="query" value="query"></el-option>
            <el-option label="callback" value="callback"></el-option>
          </el-select>
        </el-form-item>
           <el-form-item label="响应ID" prop="response_id">
          <el-select v-model="formData.response_id" placeholder="请选择响应ID" filterable :loading="responseIdLoading">
            <el-option v-for="option in responseIdOptions" :key="option.value" :label="option.label" :value="option.value"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="用户标识" prop="user_identity">
          <el-input v-model="formData.user_identity" placeholder="请输入用户标识"></el-input>
        </el-form-item>
     
        <el-form-item label="过期时间" prop="expire_time">
          <el-date-picker
            v-model="formData.expire_time"
            type="datetime"
            placeholder="请选择过期时间"
            style="width: 100%;"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="规则状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择规则状态">
            <el-option label="启用" value="on"></el-option>
            <el-option label="禁用" value="off"></el-option>
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

    <!-- 发送回调对话框 -->
    <el-dialog
      title="发送回调"
      v-model="callbackDialogVisible"
      width="600px"
    >
      <el-form ref="callbackFormRef" :model="callbackFormData" :rules="callbackRules" label-width="100px">
        <el-form-item v-if="showOrderIdInput" label="订单ID" prop="order_id">
          <el-input v-model="callbackFormData.order_id" placeholder="请输入订单ID"></el-input>
        </el-form-item>
        <el-form-item label="回调参数" prop="callback_params">
          <el-input
            type="textarea"
            v-model="callbackFormData.callback_params"
            placeholder="回调参数"
            rows="8"
            style="font-family: monospace;"
          ></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="callbackDialogVisible = false">取消</el-button>
          <el-button v-if="!callbackFormData.callback_params && showOrderIdInput" type="primary" @click="handleGetCallbackInfo">获取回调信息</el-button>
          <el-button v-else type="success" @click="handleSendCallbackRequest">发送</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, watch, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { channelMockRuleApi, channelMockResponseApi } from '@/api/channelMockApi'

// 渠道ID选项列表
const channelIdOptions = [
  { value: 'payermax_sab', label: 'payermax_sab' },
  { value: 'checkout-hk', label: 'checkout-hk' },
  { value: 'checkout-sab', label: 'checkout-sab' },
  { value: 'payermax-test-hk', label: 'payermax-test-hk' },
  { value: 'stcpay-test', label: 'stcpay-test' }
]

// 响应ID选项列表
const responseIdOptions = ref([])

// 响应ID加载状态
const responseIdLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  channel_id: '',
  request_type: '',
  user_identity: ''
})

// 规则列表
const ruleList = ref([])

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
const callbackDialogVisible = ref(false)

// 表单引用
const formRef = ref(null)
const callbackFormRef = ref(null)

// 表单数据
const formData = reactive({
  id: '',
  channel_id: '',
  request_type: '',
  user_identity: '',
  response_id: '',
  expire_time: '',
  status: ''
})

// 回调表单数据
const callbackFormData = reactive({
  order_id: '',
  response_id: '',
  channel_id: '',
  callback_params: ''
})

// 显示订单ID输入框的标志
const showOrderIdInput = ref(true)

// 表单验证规则
const rules = {
  channel_id: [{ required: true, message: '请输入渠道ID', trigger: 'blur' }],
  request_type: [{ required: true, message: '请选择请求类型', trigger: 'change' }],
  user_identity: [{ required: true, message: '请输入用户标识', trigger: 'blur' }],
  response_id: [{ required: true, message: '请输入响应ID', trigger: 'blur' }],
  status: [{ required: true, message: '请选择规则状态', trigger: 'change' }]
}

// 回调表单验证规则
const callbackRules = {
  order_id: [{ required: true, message: '请输入订单ID', trigger: 'blur' }]
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  if (isNaN(date.getTime())) return '-'
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
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

// 表格行样式
const tableRowClassName = ({ row }) => {
  // 判断是否过期：如果过期时间存在且小于当前时间
  const isExpired = row.expire_time && new Date(row.expire_time) < new Date();
  // 判断状态是否为禁用：如果状态存在且为'off'
  const isDisabled = row.status && row.status.toLowerCase() === 'off';
  
  if (isExpired) {
    return 'row-expired';
  } else if (isDisabled) {
    return 'row-disabled';
  } else {
    return '';
  }
}

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

// 根据渠道ID和请求类型获取响应ID选项（带防抖）
const fetchResponseIds = debounce(async () => {
  const { channel_id, request_type } = formData;
  
  // 只有当渠道ID和请求类型都不为空时，才获取响应ID选项
  if (!channel_id || !request_type) {
    responseIdOptions.value = [];
    formData.response_id = '';
    return;
  }
  
  // 显示加载状态
  responseIdLoading.value = true;
  
  try {
    // 调用API获取响应配置列表
    const response = await channelMockResponseApi.getResponseConfigList({ channel_id, request_type });
    if (response.code === 200 && response.data && response.data.length > 0) {
      // 转换响应数据为选项格式
      responseIdOptions.value = response.data.map(item => ({
        value: item.response_id,
        label: item.response_id
      }));
    } else {
      responseIdOptions.value = [];
      formData.response_id = '';
    }
  } catch (error) {
    console.error('获取响应ID选项失败:', error);
    responseIdOptions.value = [];
    formData.response_id = '';
  } finally {
    // 隐藏加载状态
    responseIdLoading.value = false;
  }
}, 300);

// 监听渠道ID变化，获取响应ID选项
watch(() => formData.channel_id, () => {
  fetchResponseIds();
});

// 监听请求类型变化，获取响应ID选项
watch(() => formData.request_type, () => {
  fetchResponseIds();
});

// 获取规则配置列表
const fetchRuleList = async () => {
  loading.value = true;
  try {
    const response = await channelMockRuleApi.getRuleConfigList(searchForm);
    if (response.code === 200) {
      ruleList.value = response.data || [];
      pagination.total = response.total || 0;
    } else {
      ElMessage.error(response.msg || '获取规则配置列表失败');
    }
  } catch (error) {
    ElMessage.error('获取规则配置列表失败');
    console.error('获取规则配置列表失败:', error);
  } finally {
    loading.value = false;
  }
};

// 重置表单
const resetForm = () => {
  searchForm.channel_id = '';
  searchForm.request_type = '';
  searchForm.user_identity = '';
  fetchRuleList();
};

// 处理分页大小变化
const handleSizeChange = (size) => {
  pagination.pageSize = size;
  pagination.currentPage = 1;
  fetchRuleList();
};

// 处理当前页码变化
const handleCurrentChange = (currentPage) => {
  pagination.currentPage = currentPage;
  fetchRuleList();
};

// 处理编辑
const handleEdit = (row) => {
  Object.assign(formData, row);
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

// 处理新增
const handleAdd = () => {
  // 重置表单数据
  Object.assign(formData, {
    id: null,
    channel_id: '',
    request_type: '',
    user_identity: '',
    response_id: '',
    expire_time: getDefaultExpireTime(),
    status: 'on'
  });
  
  dialogVisible.value = true;
};

// 处理删除
const handleDelete = (id) => {
  ElMessageBox.confirm('确定要删除这条规则配置吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const response = await channelMockRuleApi.deleteRuleConfig({ id });
      if (response.code === 200) {
        ElMessage.success('删除成功');
        fetchRuleList();
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
    const apiMethod = isEdit ? channelMockRuleApi.updateRuleConfig : channelMockRuleApi.addRuleConfig;
    
    // 处理id字段类型问题
    const submitData = { ...formData };
    if (!isEdit) {
      // 新增时删除id字段，避免类型不匹配
      delete submitData.id;
    } else {
      // 编辑时确保id是数字类型
      submitData.id = Number(submitData.id);
    }
    
    // 转换过期时间格式为ISO 8601格式
    if (submitData.expire_time) {
      if (typeof submitData.expire_time === 'string' && !submitData.expire_time.includes('T')) {
        submitData.expire_time = submitData.expire_time.replace(' ', 'T') + 'Z';
      }
    }
    
    const response = await apiMethod(submitData);
    if (response.code === 200) {
      ElMessage.success(isEdit ? '更新成功' : '新增成功');
      dialogVisible.value = false;
      fetchRuleList();
    } else {
      ElMessage.error(response.msg || (isEdit ? '更新失败' : '新增失败'));
    }
  } catch (error) {
    // 表单验证失败
    return false;
  }
};

// 处理发送回调按钮点击事件
const handleSendCallback = (row) => {
  // 重置回调表单数据
  Object.assign(callbackFormData, {
    order_id: '',
    response_id: row.response_id,
    channel_id: row.channel_id,
    callback_params: ''
  });
  
  // 判断是否显示订单ID输入框
  showOrderIdInput.value = !!row.response_param;
  
  // 直接使用当前行的数据来填充回调表单
  // 判断响应参数是否为空
  if (!row.response_param) {
    // 响应参数为空，直接取响应内容的值填入回调参数
    callbackFormData.callback_params = row.mocked_body;
  }
  
  // 打开回调对话框
  callbackDialogVisible.value = true;
};

// 处理获取回调信息
const handleGetCallbackInfo = async () => {
  if (!callbackFormRef.value) return;
  
  try {
    await callbackFormRef.value.validate();
    
    // 调用API获取回调信息
    const response = await channelMockRuleApi.getCallbackData({ 
      order_id: callbackFormData.order_id,
      response_id: callbackFormData.response_id,
      channel_id: callbackFormData.channel_id
    });
    
    if (response.code === 200) {
      // 成功获取回调参数，显示在对话框中
      if (response.data) {
        callbackFormData.callback_params = JSON.stringify(response.data, null, 2);
        ElMessage.success('获取回调信息成功，请确认后发送');
      } else {
        ElMessage.error('获取回调信息失败：响应数据为空');
      }
    } else {
      ElMessage.error(response.msg || '获取回调信息失败');
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('获取回调信息失败');
      console.error('获取回调信息失败:', error);
    }
  }
};

// 处理发送回调请求
const handleSendCallbackRequest = async () => {
  try {
    // 只有当显示订单ID输入框时，才需要验证表单
    if (showOrderIdInput.value && callbackFormRef.value) {
      await callbackFormRef.value.validate();
    }
    
    // 调用API发送回调
    const response = await channelMockRuleApi.sendCallback({ 
      order_id: callbackFormData.order_id,
      response_id: callbackFormData.response_id,
      channel_id: callbackFormData.channel_id,
      callback_params: callbackFormData.callback_params 
    });
    
    if (response.code === 200) {
      ElMessage.success('发送回调成功:' + response.data);
      callbackDialogVisible.value = false;
    } else {
      ElMessage.error(response.msg || '发送回调失败');
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('发送回调失败');
      console.error('发送回调失败:', error);
    }
  }
};

// 组件挂载时获取规则列表
onMounted(() => {
  fetchRuleList();
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

.channel-mock-rule-config {
  padding: 20px;
}

/* 成功状态样式 */
:deep(.status-success) {
  color: #67c23a !important;
  background-color: #f0f9eb !important;
  padding: 2px 8px !important;
  border-radius: 4px !important;
  font-weight: 500 !important;
  display: inline-block !important;
}

/* 失败状态样式 */
:deep(.status-failed) {
  color: #f56c6c !important;
  background-color: #fff1f0 !important;
  padding: 2px 8px !important;
  border-radius: 4px !important;
  font-weight: 500 !important;
  display: inline-block !important;
}

/*  pending状态样式 */
:deep(.status-pending) {
  color: #409eff !important;
  background-color: #ecf5ff !important;
  padding: 2px 8px !important;
  border-radius: 4px !important;
  font-weight: 500 !important;
  display: inline-block !important;
}

/* 行置灰样式实现 */
:deep(.el-table__row.row-disabled) {
  background-color: #f5f5f5 !important;
}

/* 确保表格单元格内容置灰 */
:deep(.el-table__row.row-disabled .el-table__cell) {
  color: #c0c4cc !important;
}

/* 保持按钮原始样式 */
:deep(.el-table__row.row-disabled .el-button--primary) {
  background-color: #409eff !important;
  border-color: #409eff !important;
  color: #fff !important;
  opacity: 1 !important;
}

:deep(.el-table__row.row-disabled .el-button--danger) {
  background-color: #f56c6c !important;
  border-color: #f56c6c !important;
  color: #fff !important;
  opacity: 1 !important;
}

/* 过期行样式 */
:deep(.el-table__row.row-expired) {
  background-color: #fff1f0 !important;
}
/* 
/* 确保表格单元格内容红色 */
/* :deep(.el-table__row.row-expired .el-table__cell) {
  color: #0c0c0c !important;
} */

/* 保持按钮原始样式 */
:deep(.el-table__row.row-expired .el-button--primary) {
  background-color: #409eff !important;
  border-color: #409eff !important;
  color: #fff !important;
  opacity: 1 !important;
} 

:deep(.el-table__row.row-expired .el-button--danger) {
  background-color: #f56c6c !important;
  border-color: #f56c6c !important;
  color: #fff !important;
  opacity: 1 !important;
}

/* 保持标签样式 */
:deep(.el-table__row.row-disabled .el-tag--success) {
  background-color: #f0f9eb !important;
  border-color: #e1f3d8 !important;
  color: #67c23a !important;
  opacity: 0.7 !important;
}

:deep(.el-table__row.row-disabled .el-tag--danger) {
  background-color: #fef0f0 !important;
  border-color: #fbc4c4 !important;
  color: #f56c6c !important;
  opacity: 0.7 !important;
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