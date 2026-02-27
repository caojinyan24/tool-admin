<template>
  <div>
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑任务' : '添加任务' }}</span>
        </div>
      </template>
      
      <el-form 
        ref="taskFormRef" 
        :model="taskForm" 
        :rules="rules" 
        label-width="100px"
      >
      
        <el-form-item label="任务标识" prop="taskName">
          <el-input v-model="taskForm.taskName" placeholder="请输入任务名称" @input="handleFieldChange"></el-input>
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            任务唯一标识, 创建完成后不可变更
          </div>
        </el-form-item>
         <el-form-item label="任务描述" prop="description">
          <el-input
            v-model="taskForm.description"
            type="textarea"
            placeholder="请输入任务描述"
            :rows="1"
            @input="handleFieldChange"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="接口地址" prop="apiAddr">
          <el-input v-model="taskForm.apiAddr" placeholder="请输入接口地址" @input="handleFieldChange"></el-input>
           <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            接口地址，如：http://localhost:8080, 不带路径
          </div>
        </el-form-item>
                <el-form-item label="接口名称" prop="apiName">
          <el-input v-model="taskForm.apiName" placeholder="请输入接口名称" @input="handleFieldChange"></el-input>
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            接口名称，如：api/v1/hello
          </div>
        </el-form-item>
        
        <el-form-item label="接口参数" prop="taskParam">
          <el-input
            v-model="taskForm.taskParam"
            type="textarea"
            placeholder="请输入接口参数(JSON格式) POST请求"
            :rows="4"
            @input="handleFieldChange"
          ></el-input>
        </el-form-item>
        
        <el-form-item label="调度时间" prop="cronExpr">
          <el-input v-model="taskForm.cronExpr" placeholder="请输入cron表达式，如：0 */5 * * * *" @input="handleFieldChange"></el-input>
          <div style="color: #909399; font-size: 12px; margin-top: 5px;">
            格式说明：秒 分 时 日 月 周 (如：0 */5 * * * * 表示每5分钟执行一次)
          </div>
        </el-form-item>
        
        <el-form-item label="接口类型" prop="apiType">
          <el-select v-model="taskForm.apiType" placeholder="请选择接口类型" @change="handleFieldChange">
            <el-option label="HTTP" value="http"></el-option>
            <el-option label="gRPC" value="grpc"></el-option>
          </el-select>
        </el-form-item>
        

        
        <el-form-item label="任务状态" prop="status">
          <el-radio-group v-model="taskForm.status" @change="handleFieldChange">
            <el-radio :label="'on'">启用</el-radio>
            <el-radio :label="'off'">停用</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="重试间隔(s)" prop="retryInterval">
          <el-input-number v-model="taskForm.retryInterval" :min="0" placeholder="重试间隔" @change="handleFieldChange"></el-input-number>
        </el-form-item>
        
        <el-form-item label="重试次数" prop="retryCount">
          <el-input-number v-model="taskForm.retryCount" :min="0" placeholder="失败重试次数" @change="handleFieldChange"></el-input-number>
        </el-form-item>
        
        <el-form-item label="失败通知列表" prop="notifier">
          <el-input
              v-model="taskForm.notifier"
              placeholder=" 失败通知列表, 使用邮箱前缀的用户名, 多个用户名用逗号分隔, 示例: jinyan.cao, jinyan.cao1"
              @input="handleNotifierChange"
            ></el-input>
            <div style="color: #909399; font-size: 12px; margin-top: 5px;">
              失败通知列表, 使用邮箱前缀的用户名, 多个用户名用逗号分隔, 示例: jinyan.cao, jinyan.cao1
            </div>
        </el-form-item>
        
       
        
        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :disabled="isEdit && !isFormModified">提交</el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import taskApi from '@/api/taskApi'

export default {
  name: 'TaskForm',
  props: {
    id: {
      type: String,
      default: ''
    }
  },
  setup(props) {
    const router = useRouter()
    const route = useRoute()
    const taskFormRef = ref(null)
    const isEdit = ref(false)
    
    const taskForm = reactive({
      id: 0,
      taskName: '',
      apiAddr: '',
      taskParam: '',
      apiType: 'http',
      apiName: '',
      cronExpr: '',
      status: 'on',
      retryInterval: 0,
      retryCount: 0,
      description: '',
      notifier: '',
      originalNotifier: ''
    })
    
    // 表单验证规则
    const rules = {
      taskName: [
        { required: true, message: '请输入任务名称', trigger: 'blur' },
        { min: 2, max: 100, message: '长度在 2 到 100 个字符,尽量使用英文字符', trigger: 'blur' }
      ],
      apiAddr: [
        { required: true, message: '请输入接口地址', trigger: 'blur' }
      ],
      apiType: [
        { required: true, message: '请选择接口类型', trigger: 'change' }
      ],
      apiName: [
        { required: true, message: '请输入接口名称', trigger: 'blur' },
        { min: 2, max: 100, message: '长度在 2 到 100 个字符', trigger: 'blur' }
      ],
      cronExpr: [
        { required: false, message: '请输入调度时间表达式', trigger: 'blur' },
        { pattern: /^[\d\s\*\/\-\,\?LWC#A-Z]+$/i, message: '请输入有效的cron表达式', trigger: 'blur' }
      ]
    }
    
    // 跟踪notifier字段是否被修改
    const notifierModified = ref(false)
    
    // 跟踪整个表单是否被修改
    const isFormModified = ref(false)
    
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

    // 加载任务详情
    const loadTaskDetail = async (taskName) => {
      try {
        const response = await taskApi.getTaskByTaskName(taskName)
        const taskData = response.data
        
        // 重置修改标记
        notifierModified.value = false
        isFormModified.value = false
        
        // 将后端返回的蛇形命名数据转换为前端使用的驼峰命名
        taskForm.id = taskData.id
        taskForm.taskName = taskData.task_name || taskData.taskName
        taskForm.apiAddr = taskData.api_addr || taskData.apiAddr
        taskForm.taskParam = taskData.task_param || taskData.taskParam
        taskForm.apiType = taskData.api_type || taskData.apiType
        taskForm.apiName = taskData.api_name || taskData.apiName
        taskForm.cronExpr = taskData.cron_expr || taskData.cronExpr
        taskForm.status = taskData.status
        taskForm.retryInterval = taskData.retry_interval || taskData.retryInterval || 0
        taskForm.retryCount = taskData.retry_count || taskData.retryCount || 0
        taskForm.description = taskData.description || ''
        taskForm.notifier = formatNotifier(taskData.notifier || '')
        taskForm.originalNotifier = taskData.notifier || ''
      } catch (error) {
        ElMessage.error('加载任务详情失败')
        console.error('加载任务详情失败:', error)
        router.push('/tasks')
      }
    }
    
    // 处理表单字段变更
    const handleFieldChange = () => {
      isFormModified.value = true
    }
    
    // 专门处理notifier字段变更
    const handleNotifierChange = () => {
      notifierModified.value = true
      isFormModified.value = true
    }
    
    // 将逗号分隔的用户名前缀转换为数据库所需的JSON格式
    const formatNotifierForSave = (notifierText) => {
      if (!notifierText) return '{}'
      
      // 分割逗号分隔的用户名
      const usernames = notifierText.split(',').map(name => name.trim()).filter(name => name)
      
      // 构建JSON对象
      const notifierObj = {}
      usernames.forEach(username => {
        // 假设域名是jaco.live，实际应该根据需求调整
        const email = `${username}@jaco.live`
        notifierObj[email] = {}
      })
      
      return JSON.stringify(notifierObj)
    }

    // 提交表单
    const handleSubmit = async () => {
      if (!taskFormRef.value) return
      
      try {
        await taskFormRef.value.validate()
        
        // 验证JSON格式的参数
        if (taskForm.taskParam) {
          try {
            JSON.parse(taskForm.taskParam)
          } catch (error) {
            ElMessage.error('接口参数格式不正确，请输入有效的JSON字符串')
            return
          }
        }
        
        // 转换驼峰命名为蛇形命名
        const snakeForm = {
          id: taskForm.id,
          task_name: taskForm.taskName,
          api_addr: taskForm.apiAddr,
          task_param: taskForm.taskParam,
          api_type: taskForm.apiType,
          api_name: taskForm.apiName,
          cron_expr: taskForm.cronExpr,
          status: taskForm.status,
          retry_interval: taskForm.retryInterval,
          retry_count: taskForm.retryCount,
          description: taskForm.description,
          notifier: isEdit.value && !notifierModified.value ? taskForm.originalNotifier : formatNotifierForSave(taskForm.notifier)
        }
        
        if (isEdit.value) {
          // 更新任务
          await taskApi.updateTask(snakeForm)
          ElMessage.success('更新成功')
        } else {
          // 创建任务
          await taskApi.createTask(snakeForm)
          ElMessage.success('创建成功')
        }
        
        // 返回任务列表
        router.push('/tasks')
      } catch (error) {
        if (error !== false) { // false 是validate失败时的返回值
          // 显示包含具体错误信息的提示
          const errorMsg = error.msg || '操作失败'
          const errorDetail = error.error ? ` (${error.error})` : ''
          ElMessage.error(`${errorMsg}${errorDetail}`)
          console.error('提交表单失败:', error)
        }
      }
    }
    
    // 取消操作
    const handleCancel = () => {
      router.push('/tasks')
    }
    
    // 初始化
    onMounted(() => {
      // 从路由参数或props中获取id
      const taskName = props.taskName || route.params.taskName
      if (taskName) {
        isEdit.value = true
        loadTaskDetail(taskName)
      }
    })
    
    return {
      taskFormRef,
      isEdit,
      taskForm,
      rules,
      notifierModified,
      isFormModified,
      handleSubmit,
      handleCancel,
      handleFieldChange,
      handleNotifierChange
    }
  }
}
</script>

<style scoped>
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>