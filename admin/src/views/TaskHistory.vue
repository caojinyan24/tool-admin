<template>
  <div class="task-history-container">
    <div class="header">
      <div class="header-actions">
        <el-button @click="handleBack" type="primary" style="margin-right: 20px;">
          返回任务列表
        </el-button>
        <h1>{{ taskName }} - 执行历史</h1>
      </div>
    </div>
    
    <!-- 搜索条件区域 -->
    <div class="search-container" style="margin-bottom: 20px;">
      <el-form :inline="true" :model="searchForm" class="demo-form-inline">
        <el-form-item label="执行状态">
            <el-select v-model="searchForm.status" placeholder="请选择执行状态" clearable style="width: 140px;">
              <el-option
                v-for="option in statusOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              ></el-option>
            </el-select>
          </el-form-item>
          
          <el-form-item label="执行类型">
            <el-select v-model="searchForm.executeType" placeholder="请选择执行类型" clearable style="width: 140px;">
              <el-option
                v-for="option in executeTypeOptions"
                :key="option.value"
                :label="option.label"
                :value="option.value"
              ></el-option>
            </el-select>
          </el-form-item>
        
        <el-form-item label="执行时间">
          <el-date-picker
            v-model="searchForm.startTime"
            type="datetime"
            placeholder="开始时间"
            style="width: 160px; margin-right: 10px;"
            :disabled-date="disabledStartDate"
          />
          <span>-</span>
          <el-date-picker
            v-model="searchForm.endTime"
            type="datetime"
            placeholder="结束时间"
            style="width: 160px; margin-left: 10px;"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    
    <div class="content">
      <el-table
        v-loading="loading"
        :data="logsData"
        style="width: 100%"
        border
      >
        <el-table-column prop="id" label="日志ID" width="80" />
        <el-table-column prop="taskName" label="任务名称" />
        <el-table-column prop="executeTime" label="执行时间" width="200">
          <template #default="{ row }">
            {{ formatDate(row.executeTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="executeType" label="执行类型" width="120">
          <template #default="{ row }">
            {{ row.executeType || '未知' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="执行状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="taskParam" label="任务参数">
          <template #default="{ row }">
            <el-popover trigger="hover" placement="top">
              <div style="max-width: 600px; word-break: break-all;">{{ row.taskParam || '无' }}</div>
              <template #reference>
                <span>{{ (row.taskParam || '').length > 20 ? (row.taskParam).substring(0, 20) + '...' : (row.taskParam || '无') }}</span>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="result" label="执行结果">
          <template #default="{ row }">
            <el-popover trigger="hover" placement="top">
              <div style="max-width: 600px; word-break: break-all;">{{ row.result || '无' }}</div>
              <template #reference>
                <span>{{ (row.result || '').length > 20 ? (row.result).substring(0, 20) + '...' : (row.result || '无') }}</span>
              </template>
            </el-popover>
          </template>
        </el-table-column>
        <el-table-column prop="traceId" label="追踪ID" width="180">
          <template #default="{ row }">
            <el-popover trigger="hover" placement="top">
              <div style="max-width: 400px; word-break: break-all;">{{ row.traceId || '无' }}</div>
              <template #reference>
                <span>{{ (row.traceId || '').length > 15 ? (row.traceId).substring(0, 15) + '...' : (row.traceId || '无') }}</span>
              </template>
            </el-popover>
          </template>
        </el-table-column>
      </el-table>
      
     <div class="pagination-container">
      <el-pagination v-model:current-page="pagination.currentPage" v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange"></el-pagination>
    </div>
    </div>
    
    <!-- 日志详情弹窗 -->
    <el-dialog
      v-model="dialogVisible"
      title="执行详情"
      width="70%"
      :before-close="handleClose"
    >
      <div v-if="currentLogDetail" class="log-detail">
        <el-descriptions column="1" border>
          <el-descriptions-item label="任务名称">{{ currentLogDetail.taskName }}</el-descriptions-item>
          <el-descriptions-item label="执行时间">{{ formatDate(currentLogDetail.executeTime) }}</el-descriptions-item>
          <el-descriptions-item label="执行状态">{{ currentLogDetail.status === 1 ? '成功' : '失败' }}</el-descriptions-item>
          <el-descriptions-item label="执行时长">{{ currentLogDetail.duration }}ms</el-descriptions-item>
          <el-descriptions-item label="服务器IP">{{ currentLogDetail.serverIp }}</el-descriptions-item>
          <el-descriptions-item label="响应时间">{{ currentLogDetail.responseTime }}ms</el-descriptions-item>
          <el-descriptions-item label="执行轨迹">{{ currentLogDetail.executionTrace }}</el-descriptions-item>
          <el-descriptions-item label="请求头">
            <pre>{{ currentLogDetail.requestHeaders }}</pre>
          </el-descriptions-item>
          <el-descriptions-item label="请求体">
            <pre>{{ currentLogDetail.requestBody }}</pre>
          </el-descriptions-item>
          <el-descriptions-item label="执行结果">
            <pre>{{ currentLogDetail.result }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import logApi from '@/api/logApi.js'
import { ElMessage } from 'element-plus'

export default {
  name: 'TaskHistory',
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    // 获取URL参数中的任务名称
    const taskName = ref(route.query.taskName || '')
    const loading = ref(false)
    const logsData = ref([])
    const total = ref(0)
    const dialogVisible = ref(false)
    const currentLogDetail = ref(null)
    
    // 分页配置
    const pagination = ref({
      currentPage: 1,
      pageSize: 10
    })
    
    // 搜索条件
    const searchForm = reactive({
      status: '',
      executeType: '',
      startTime: '',
      endTime: ''
    })
    
    // 执行状态选项 - 简化实现，直接定义为常量数组
    const statusOptions = [
      { label: '全部', value: '' },
      { label: '初始化', value: 'I' },
      { label: '成功', value: 'S' },
      { label: '失败', value: 'F' }
    ]
    
    // 执行类型选项 - 简化实现，直接定义为常量数组
    const executeTypeOptions = [
      { label: '全部', value: '' },
      { label: '手动执行', value: 'manual' },
      { label: '自动调度', value: 'scheduled' }
    ]
    
    // 初始化时设置默认时间范围为最近一个月
    const initDefaultTimeRange = () => {
      const now = new Date()
      const oneMonthAgo = new Date()
      oneMonthAgo.setMonth(now.getMonth() - 1)
      
      // 设置默认开始时间为一个月前，结束时间为当前时间
      searchForm.startTime = oneMonthAgo
      searchForm.endTime = now
    }
    
    // 检查日期是否在一个月内
    const isWithinOneMonth = (date) => {
      const now = new Date()
      const oneMonthAgo = new Date()
      oneMonthAgo.setMonth(now.getMonth() - 1)
      return date >= oneMonthAgo && date <= now
    }
    
    // 禁用早于一个月前的日期
    const disabledStartDate = (date) => {
      const now = new Date()
      const oneMonthAgo = new Date()
      oneMonthAgo.setMonth(now.getMonth() - 1)
      return date < oneMonthAgo
    }
    
    // 格式化日期
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      if (isNaN(date.getTime())) {
        // 如果日期解析失败，尝试处理其他可能的格式
        return dateString
      }
      // 直接使用字符串拼接确保格式一致且只显示到秒
      const year = date.getFullYear()
      const month = String(date.getMonth() + 1).padStart(2, '0')
      const day = String(date.getDate()).padStart(2, '0')
      const hours = String(date.getHours()).padStart(2, '0')
      const minutes = String(date.getMinutes()).padStart(2, '0')
      const seconds = String(date.getSeconds()).padStart(2, '0')
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    }
    
    // 获取状态标签类型
    const getStatusTagType = (status) => {
      switch (status) {
        case 'S':
        case '1':
          return 'success'
        case 'F':
        case '0':
          return 'danger'
        case 'I':
          return 'info'
        default:
          return 'default'
      }
    }
    
    // 获取状态文本
    const getStatusText = (status) => {
      switch (status) {
        case 'S':
        case '1':
          return '成功'
        case 'F':
        case '0':
          return '失败'
        case 'I':
          return '初始化'
        default:
          return '未知'
      }
    }
    
    // 加载任务执行历史
    const loadTaskHistory = async () => {
      if (!taskName.value) {
        ElMessage.warning('任务名称不能为空')
        return
      }
      
      loading.value = true
      try {
        const response = await logApi.getLogs({
          page: pagination.currentPage,
          pageSize: pagination.pageSize,
          taskName: taskName.value,
          status: searchForm.status,
          executeType: searchForm.executeType,
          startTime: searchForm.startTime ? searchForm.startTime.toISOString() : '',
          endTime: searchForm.endTime ? searchForm.endTime.toISOString() : ''
        })
        
        // 确保logsData始终是一个数组
        // 同时处理两种可能的响应格式：后端直接返回的格式和logApi.js中模拟数据的格式
        if (response) {
          // 从不同位置获取状态码
          let code = response.code || (response.data && response.data.code)
          // 获取数据部分
          let data = response.data
          
          if (code === 200 && data) {
            // 处理数据，如果存在嵌套的data部分则使用它
            const innerData = data.data || data
            // 获取日志列表，如果存在list字段则使用list，否则使用data
            const logsList = innerData.list || (Array.isArray(innerData.data) ? innerData.data : null)
            if (logsList && Array.isArray(logsList)) {
              // 转换数据格式以匹配前端期望的字段名和结构
            logsData.value = logsList.map(log => ({
              id: log.id,
              taskName: log.task_name || log.taskName || '',
              executeTime: log.execute_time || log.executeTime || '',
              executeType: log.execute_type || log.executeType || '',
              status: log.status || '',
              taskParam: log.task_param || log.taskParam || '',
              result: log.result || '',
              traceId: log.trace_id || log.traceId || ''
            }))
              
              // 尝试从不同位置获取总数，确保能正确显示
              total.value = 
                (innerData.count || innerData.total) || 
                (data.count || data.total) || 
                (logsList && logsList.length ? logsList.length : 0)
            } else {
              logsData.value = []
              total.value = 0
            }
          } else {
            logsData.value = []
            total.value = 0
            ElMessage.error(response?.message || response?.data?.message || '获取执行历史失败')
          }
        } else {
          // 没有有效数据，设置为空数组
          logsData.value = []
          total.value = 0
          ElMessage.error(response?.message || '获取执行历史失败')
        }
      } catch (error) {
        console.error('获取执行历史失败:', error)
        ElMessage.error('获取执行历史失败: ' + (error.message || '未知错误'))
        // 捕获异常时也要确保logsData是数组
        logsData.value = []
        total.value = 0
      } finally {
        loading.value = false
      }
    }
    
    // 搜索方法
    const handleSearch = () => {
      // 检查并限制开始时间不能早于一个月前
      const now = new Date()
      const oneMonthAgo = new Date()
      oneMonthAgo.setMonth(now.getMonth() - 1)
      
      if (searchForm.startTime && searchForm.startTime < oneMonthAgo) {
        // 强制将开始时间设置为一个月前
        searchForm.startTime = oneMonthAgo
        ElMessage.info('只支持查询最近一个月的日志')
      }
      
      // 重置当前页码为1
      pagination.value.currentPage = 1
      // 执行搜索
      loadTaskHistory()
    }
    
    // 重置搜索条件
    const handleReset = () => {
      searchForm.status = ''
      searchForm.executeType = ''
      searchForm.startTime = ''
      searchForm.endTime = ''
      // 重置当前页码为1
      pagination.value.currentPage = 1
      // 重新加载数据
      loadTaskHistory()
    }
    
    // 查看日志详情
    const handleViewDetail = async (logId) => {
      try {
        const response = await logApi.getLogDetail(logId)
        if (response) {
          // 从不同位置获取状态码和数据
          let code = response.code || (response.data && response.data.code)
          let data = response.data
          
          if (code === 200 && data) {
            // 处理可能的嵌套数据结构
            const logData = data.data || data
            
            // 转换数据格式以匹配前端期望的字段名和结构
            currentLogDetail.value = {
              id: logData.id || '',
              taskName: logData.task_name || logData.taskName || '',
              executeTime: logData.execute_time || logData.executeTime || '',
              status: logData.status === 'S' ? 1 : 0,
              duration: 0, // 后端数据中没有提供duration
              serverIp: '', // 后端数据中没有提供serverIp
              responseTime: 0, // 后端数据中没有提供responseTime
              executionTrace: '', // 后端数据中没有提供executionTrace
              requestHeaders: '', // 后端数据中没有提供requestHeaders
              requestBody: logData.task_param || logData.taskParam || '',
              result: logData.result || ''
            }
            dialogVisible.value = true
          } else {
            ElMessage.error(response?.message || response?.data?.message || '获取日志详情失败')
          }
        } else {
          ElMessage.error('获取日志详情失败')
        }
      } catch (error) {
        console.error('获取日志详情失败:', error)
        ElMessage.error('获取日志详情失败')
      }
    }
    
    // 分页大小变化
    const handleSizeChange = (size) => {
      pagination.value.pageSize = size
      pagination.value.currentPage = 1
      loadTaskHistory()
    }
    
    // 当前页变化
    const handleCurrentChange = (current) => {
      pagination.value.currentPage = current
      loadTaskHistory()
    }
    
    // 关闭弹窗
    const handleClose = () => {
      dialogVisible.value = false
      currentLogDetail.value = null
    }
    
    // 监听任务名称变化
    watch(() => route.query.taskName, (newTaskName) => {
      taskName.value = newTaskName || ''
      if (newTaskName) {
        // 重置搜索条件
        handleReset()
      }
    })
    
    // 初始化加载数据
    onMounted(() => {
      // 设置默认时间范围为最近一个月
      initDefaultTimeRange()
      
      if (taskName.value) {
        loadTaskHistory()
      } else {
        ElMessage.warning('未找到任务名称参数')
      }
    })
    
    // 返回任务列表
    const handleBack = () => {
      router.push('/tasks')
    }

    return {
      taskName,
      loading,
      logsData,
      total,
      dialogVisible,
      currentLogDetail,
      pagination,
      searchForm,
      statusOptions,
      executeTypeOptions,
      formatDate,
      getStatusTagType,
      getStatusText,
      loadTaskHistory,
      handleSearch,
      handleReset,
      handleViewDetail,
      handleSizeChange,
      handleCurrentChange,
      handleClose,
      handleBack,
      initDefaultTimeRange,
      isWithinOneMonth,
      disabledStartDate
    }
  }
}
</script>

<style scoped>
.task-history-container {
  padding: 20px;
}

.header {
  margin-bottom: 20px;
}

.header-actions {
  display: flex;
  align-items: center;
}

.pagination-container {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.log-detail pre {
  white-space: pre-wrap;
  word-break: break-word;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  max-height: 300px;
  overflow-y: auto;
}
</style>