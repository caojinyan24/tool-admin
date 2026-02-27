<template>
    <div class="search-form">
      <el-form :inline="true" :model="searchForm" size="small">
        <el-form-item label="任务名称">
          <el-input v-model="searchForm.taskName" placeholder="请输入任务名称" style="width: 200px;"></el-input>
        </el-form-item>
        <el-form-item label="任务状态">
          <el-select v-model="searchForm.status" placeholder="请选择任务状态" style="width: 120px;">
            <el-option label="全部" value=""></el-option>
            <el-option label="启用" value="on"></el-option>
            <el-option label="停用" value="off"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleAddTask">
            <el-icon>
              <Plus />
            </el-icon> 添加任务
          </el-button>
        </el-form-item>
      </el-form>

    <el-table v-loading="loading" :data="taskList" style="width: 100%">
      <el-table-column prop="id" label="ID" width="80"></el-table-column>
      <el-table-column prop="taskName" label="任务名称" width="180">
      </el-table-column>
      <el-table-column prop="cronExpr" label="调度时间" width="150"></el-table-column>
      <el-table-column prop="status" label="任务状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.status === 'on' ? 'success' : 'danger'">
            {{ scope.row.status === 'on' ? '启用' : '停用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="apiType" label="接口类型" width="100"></el-table-column>

      <el-table-column prop="apiAddr" label="接口地址" min-width="250"></el-table-column>
      <el-table-column prop="apiName" label="接口名称" width="120"></el-table-column>

      <el-table-column prop="taskParam" label="接口参数" min-width="200">
        <template #default="scope">
          <el-popover trigger="hover" placement="top">
            <div style="max-width: 400px; word-break: break-all;">{{ scope.row.taskParam }}</div>
            <template #reference>
              <span>{{ scope.row.taskParam.length > 20 ? scope.row.taskParam.substring(0, 20) + '...' :
                scope.row.taskParam }}</span>
            </template>
          </el-popover>
        </template>
      </el-table-column>

      <el-table-column prop="retryInterval" label="重试间隔(s)" width="120" />
      <el-table-column prop="retryCount" label="重试次数" width="100" />
      
      <el-table-column prop="notifier" label="失败通知列表" min-width="200">
        <template #default="scope">
          <el-popover trigger="hover" placement="top">
            <div style="max-width: 300px; word-break: break-all;">{{ scope.row.notifier || '无' }}</div>
            <template #reference>
              <span>{{ (scope.row.notifier || '').length > 20 ? (scope.row.notifier).substring(0, 20) + '...' :
                (scope.row.notifier || '无') }}</span>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      
      <el-table-column prop="description" label="任务描述" min-width="180">
        <template #default="scope">
          <el-popover trigger="hover" placement="top">
            <div style="max-width: 300px; word-break: break-all;">{{ scope.row.description || '无' }}</div>
            <template #reference>
              <span>{{ (scope.row.description || '').length > 20 ? (scope.row.description).substring(0, 20) + '...' :
                (scope.row.description || '无') }}</span>
            </template>
          </el-popover>
        </template>
      </el-table-column>
      <!-- <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column prop="updatedAt" label="更新时间" width="180">
        <template #default="scope">
          {{ formatDate(scope.row.updatedAt) }}
        </template>
      </el-table-column> -->
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="scope">
          <el-button type="primary" size="small" @click="handleEditTask(scope.row.taskName)">编辑</el-button>
          <el-button type="danger" size="small" @click="handleDeleteTask(scope.row.taskName)">删除</el-button>
          <el-button size="small" @click="handleTriggerTask(scope.row.taskName)">触发</el-button>
          <el-button size="small" @click="handleViewHistory(scope.row.taskName)">执行历史</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination v-model:current-page="pagination.currentPage" v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]" layout="total, sizes, prev, pager, next, jumper" :total="pagination.total"
        @size-change="handleSizeChange" @current-change="handleCurrentChange"></el-pagination>
    </div>

    <!-- 任务触发弹窗 -->
    <el-dialog v-model="triggerDialogVisible" title="触发任务" width="500px" :before-close="handleTriggerDialogClose">
      <el-form :model="triggerForm" ref="triggerFormRef" label-width="80px">
        <el-form-item label="任务名称" prop="taskName" :disabled="true">
          <el-input v-model="triggerForm.taskName" disabled></el-input>
        </el-form-item>
        <el-form-item label="任务参数" prop="taskParam">
          <el-input v-model="triggerForm.taskParam" type="textarea" :rows="4" placeholder="请输入任务参数（JSON格式）"></el-input>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleTriggerDialogClose">取消</el-button>
          <el-button type="primary" @click="handleConfirmTrigger">提交</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import taskApi from '@/api/taskApi'

export default {
  name: 'TaskList',
  components: {
    Plus
  },
  setup() {
    const router = useRouter()
    const taskList = ref([])
    const loading = ref(false)
    const searchForm = reactive({
      taskName: '',
      status: ''
    })
    const pagination = ref({
      currentPage: 1,
      pageSize: 10,
      total: 0
    })

    // 触发任务弹窗相关
    const triggerDialogVisible = ref(false)
    const triggerFormRef = ref(null)
    const triggerForm = ref({
      taskId: '',
      taskName: '',
      taskParam: ''
    })

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

    // 处理notifier字段展示格式
    const formatNotifier = (notifierData) => {
      if (!notifierData) return ''
      
      try {
        // 尝试将字符串解析为JSON对象
        let notifierObj
        if (typeof notifierData === 'string') {
          notifierObj = JSON.parse(notifierData)
        } else {
          notifierObj = notifierData
        }
        
        // 提取所有key并去掉邮箱后缀
        const formattedKeys = Object.keys(notifierObj).map(key => {
          const atIndex = key.indexOf('@')
          return atIndex > 0 ? key.substring(0, atIndex) : key
        })
        
        return formattedKeys.join(', ')
      } catch (error) {
        // 如果解析失败，可能已经是字符串格式，直接返回
        return notifierData
      }
    }

    // 获取任务列表
    const getTaskList = async () => {
      loading.value = true
      try {
        const response = await taskApi.getAllTasks()

        // 确保获取到正确的数据数组 - 数据在response.data.data中
        const tasksData = Array.isArray(response.data) ? response.data : (response.data?.data || [])

        // 将后端返回的蛇形命名数据转换为前端使用的驼峰命名
        const camelCaseTasks = tasksData.map(task => ({
          id: task.id,
          taskName: task.task_name || task.taskName || '',
          apiAddr: task.api_addr || task.apiAddr || '',
          taskParam: task.task_param || task.taskParam || '',
          apiType: task.api_type || task.apiType || '',
          apiName: task.api_name || task.apiName || '',
          cronExpr: task.cron_expr || task.cronExpr || '',
          status: task.status,
          retryInterval: task.retry_interval || task.retryInterval || 0,
          retryCount: task.retry_count || task.retryCount || 0,
          description: task.description || '',
          notifier: formatNotifier(task.notifier || ''),
          createdAt: task.created_at || task.createdAt,
          updatedAt: task.updated_at || task.updatedAt
        }))

        // 过滤搜索条件
        let filteredTasks = camelCaseTasks
        if (searchForm.taskName) {
          filteredTasks = filteredTasks.filter(task =>
            task.taskName.includes(searchForm.taskName)
          )
        }
        if (searchForm.status !== '') {
          filteredTasks = filteredTasks.filter(task =>
            task.status === searchForm.status
          )
        }

        // 模拟分页
        pagination.value.total = filteredTasks.length
        const start = (pagination.value.currentPage - 1) * pagination.value.pageSize
        const end = start + pagination.value.pageSize
        if (start < filteredTasks.length) {
          if (end > filteredTasks.length) {
            taskList.value = filteredTasks.slice(start)
          } else {
            taskList.value = filteredTasks.slice(start, end)
          }
        } else {
          taskList.value = []
        }
      } catch (error) {
        // 显示包含具体错误信息的提示
        const errorMsg = error.msg || '获取任务列表失败'
        const errorDetail = error.error ? ` (${error.error})` : ''
        ElMessage.error(`${errorMsg}${errorDetail}`)
        console.error('获取任务列表失败:', error)
      } finally {
        loading.value = false
      }
    }

    // 添加任务
    const handleAddTask = () => {
      router.push('/tasks/add')
    }

    // 编辑任务
    const handleEditTask = (taskName) => {
      router.push(`/tasks/edit/${taskName}`)
    }

    // 删除任务
    const handleDeleteTask = async (taskName) => {
      try {
        await ElMessageBox.confirm(
          `确定要删除任务「${taskName}」吗？`,
          '确认删除',
          {
            confirmButtonText: '确定',
            cancelButtonText: '取消',
            type: 'warning'
          }
        )

        await taskApi.deleteTask(taskName)
        ElMessage.success('删除成功')
        getTaskList() // 重新加载任务列表
      } catch (error) {
        if (error !== 'cancel') {
          // 显示包含具体错误信息的提示
          const errorMsg = error.msg || '删除失败'
          const errorDetail = error.error ? ` (${error.error})` : ''
          ElMessage.error(`${errorMsg}${errorDetail}`)
          console.error('删除任务失败:', error)
        }
      }
    }

    // 打开触发任务弹窗
    const handleTriggerTask = (taskName) => {
      const task = taskList.value.find(t => t.taskName === taskName)
      if (task) {
        triggerForm.value = {
          taskId: task.id,
          taskName: taskName,
          taskParam: ''
        }
        triggerDialogVisible.value = true
      }
    }

    // 关闭触发任务弹窗
    const handleTriggerDialogClose = () => {
      triggerDialogVisible.value = false
      if (triggerFormRef.value) {
        triggerFormRef.value.resetFields()
      }
    }

    // 确认触发任务
    const handleConfirmTrigger = async () => {
      try {
        // 验证JSON格式（如果输入了内容）
        if (triggerForm.value.taskParam) {
          try {
            JSON.parse(triggerForm.value.taskParam)
          } catch (e) {
            ElMessage.error('任务参数格式不正确，请输入有效的JSON格式')
            return
          }
        }

        await taskApi.triggerTask(triggerForm.value.taskName, triggerForm.value.taskParam)
        ElMessage.success('触发成功')
        handleTriggerDialogClose()
      } catch (error) {
        // 显示包含具体错误信息的提示
        const errorMsg = error.msg || '触发失败'
        const errorDetail = error.error ? `: ${error.error}` : ''
        ElMessage.error(`${errorMsg}, 失败原因 ${errorDetail}`)
        console.error('触发任务失败:', error)
      }
    }

    // 搜索
    const handleSearch = () => {
      pagination.value.currentPage = 1
      getTaskList()
    }

    // 重置搜索条件
    const handleReset = () => {
      searchForm.taskName = ''
      searchForm.status = ''
      pagination.value.currentPage = 1
      getTaskList()
    }

    // 分页大小变化
    const handleSizeChange = (size) => {
      pagination.value.pageSize = size
      getTaskList()
    }

    // 当前页变化
    const handleCurrentChange = (current) => {
      pagination.value.currentPage = current
      getTaskList()
    }

    // 查看任务执行历史
    const handleViewHistory = (taskName) => {
      router.push({ path: '/tasks/history', query: { taskName } })
    }

    // 初始化加载任务列表
    onMounted(() => {
      getTaskList()
    })

    return {
      taskList,
      loading,
      searchForm,
      pagination,
      triggerDialogVisible,
      triggerFormRef,
      triggerForm,
      handleAddTask,
      handleEditTask,
      handleDeleteTask,
      handleTriggerTask,
      handleConfirmTrigger,
      handleTriggerDialogClose,
      handleSearch,
      handleReset,
      handleSizeChange,
      handleCurrentChange,
      handleViewHistory,
      formatDate
    }
  }
}
</script>

<style scoped>
/* 可以添加特定的样式 */
</style>