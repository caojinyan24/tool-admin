<template>
  <div class="error-code-manager">
    <h2>错误码管理平台</h2>
    <el-tabs v-model="activeTab" class="error-code-tabs">
            <!-- 错误码映射tab -->
      <el-tab-pane label="渠道错误码映射" name="map">
        <div class="tab-content">
          <div style="margin-bottom: 16px;">
            <el-button type="primary" @click="exportErrorMaps">导出渠道错误码映射</el-button>
          </div>

          <el-table :data="errorMapDetails" stripe style="width: 100%">
                                  <el-table-column prop="id" label="ID" width="80" sortable />

            <el-table-column 
              prop="channel" 
              label="渠道" 
              width="120" 
              sortable 
              filterable
              :filters="mapChannelFilters"
              :filter-method="filterMapChannel"
            />
            <el-table-column 
              prop="channel_code" 
              label="渠道错误码" 
              width="120" 
              sortable 
              filterable
              :filter-method="filterMapChannelCode"
            />
            <el-table-column 
              prop="channel_msg" 
              label="渠道错误信息" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterMapChannelMsg"
            />
            <el-table-column 
              prop="jaco_code" 
              label="Jaco错误码" 
              width="120" 
              sortable 
              filterable
              :filter-method="filterMapJacoCode"
            />
            <el-table-column 
              prop="jaco_msg" 
              label="Jaco错误信息" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterMapJacoMsg"
            />
            <el-table-column 
              prop="user_prompt_msg" 
              label="英文提示" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterMapDisplayMsg"
            />
            <el-table-column 
              prop="user_prompt_msg_ar" 
              label="阿语提示" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterMapDisplayMsgAr"
            />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showEditMapDialog(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteErrorMap(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      <!-- 渠道错误码tab -->
      <el-tab-pane label="渠道错误码" name="channel">
        <div class="tab-content">
          <el-row :gutter="20" class="mb-4">
            <!-- <el-col :span="6">
              <el-input placeholder="渠道" v-model="channelSearch" clearable />
            </el-col> -->
            <!-- <el-col :span="4">
              <el-button type="primary" @click="searchChannelErrors">查询</el-button>
            </el-col> -->
            <el-col :span="4">
              <el-button type="success" @click="showAddChannelDialog">新增</el-button>
            </el-col>
          </el-row>
          
          <el-table :data="channelErrors" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" sortable />
            <!-- <el-table-column prop="uniq_id" label="UniqID" width="180" sortable /> -->
            <el-table-column 
              prop="channel" 
              label="渠道" 
              width="120" 
              sortable 
              filterable
              :filters="channelFilters"
              :filter-method="filterChannel"
            />
            <el-table-column 
              prop="code" 
              label="错误码" 
              width="120" 
              sortable 
              filterable
              :filter-method="filterCode"
            />
            <el-table-column 
              prop="msg" 
              label="错误信息" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterMsg"
            />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showEditChannelDialog(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteChannelError(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      
      <!-- Jaco错误码tab -->
      <el-tab-pane label="Jaco错误码" name="jaco">
        <div class="tab-content">
          <el-row :gutter="20" class="mb-4">
            <el-col :span="4">
              <el-button type="success" @click="showAddJacoDialog">新增</el-button>
            </el-col>
          </el-row>
          
          <el-table :data="jacoErrors" stripe style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" sortable />
            <el-table-column prop="type" label="类别" width="180" sortable />
            <el-table-column 
              prop="code" 
              label="错误码" 
              width="120" 
              sortable 
              filterable
              :filter-method="filterJacoCode"
            />
            <el-table-column 
              prop="msg" 
              label="错误信息" 
              width="300" 
              sortable 
              filterable
              :filter-method="filterJacoMsg"
            />
            <el-table-column 
              prop="desc" 
              label="描述" 
              width="200" 
              sortable 
              filterable
              :filter-method="filterJacoDesc"
            />
            <el-table-column 
              label="英文提示" 
              width="200" 
              sortable 
            >
              <template #default="scope">
                {{ getUserPromptMsg(scope.row.user_id) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showEditJacoDialog(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteJacoError(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      

      
      <!-- 用户提示tab -->
      <el-tab-pane label="用户提示" name="user">
        <div class="tab-content">
             <el-col :span="4">
              <el-button type="success" @click="showAddChannelDialog">新增</el-button>
            </el-col>
          
          <el-table :data="userPrompts" stripe style="width: 100%">
                        <el-table-column prop="id" label="ID" width="80" sortable />

            <el-table-column 
              prop="msg" 
              label="英文提示" 
              width="300" 
              sortable 
              filterable
            />
            <el-table-column 
              prop="msg_ar" 
              label="阿语提示" 
              width="300" 
              sortable 
              filterable
            />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button size="small" type="primary" @click="showEditUserPromptDialog(scope.row)">编辑</el-button>
                <el-button size="small" type="danger" @click="deleteUserPrompt(scope.row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>
    
    <!-- 新增/编辑渠道错误码对话框 -->
    <el-dialog
      v-model="channelDialogVisible"
      :title="channelDialogTitle"
      width="500px"
    >
      <el-form :model="channelForm" label-width="100px">
        <el-form-item label="渠道">
          <el-input v-model="channelForm.channel" placeholder="请输入渠道" />
        </el-form-item>
        <el-form-item label="错误码">
          <el-input v-model="channelForm.code" placeholder="请输入错误码" />
        </el-form-item>
        <el-form-item label="错误信息">
          <el-input v-model="channelForm.msg" placeholder="请输入错误信息" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="channelDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveChannelError">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 新增/编辑Jaco错误码对话框 -->
    <el-dialog
      v-model="jacoDialogVisible"
      :title="jacoDialogTitle"
      width="500px"
    >
      <el-form :model="jacoForm" label-width="120px">
        <el-form-item label="错误码">
          <el-input v-model="jacoForm.code" placeholder="请输入错误码" />
        </el-form-item>
        <el-form-item label="错误信息">
          <el-input v-model="jacoForm.msg" placeholder="请输入错误信息" type="textarea" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="jacoForm.desc" placeholder="请输入描述" type="textarea" />
        </el-form-item>
        <el-form-item label="英文提示">
          <el-select v-model="jacoForm.user_id" placeholder="请选择用户提示">
            <el-option
              v-for="prompt in userPrompts"
              :key="prompt.id"
              :label="prompt.msg"
              :value="prompt.msg"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="jacoDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveJacoError">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 编辑错误码映射对话框 -->
    <el-dialog
      v-model="mapDialogVisible"
      title="编辑错误码映射"
      width="500px"
    >
      <el-form :model="mapForm" label-width="120px">
        <el-form-item label="渠道">
          <el-input v-model="mapForm.channel" disabled />
        </el-form-item>
        <el-form-item label="渠道错误码">
          <el-input v-model="mapForm.channel_code" disabled />
        </el-form-item>
      
        <el-form-item label="渠道错误信息">
          <el-input v-model="mapForm.channel_msg" disabled />
        </el-form-item>
          <el-form-item label="渠道错误码唯一标识">
          <el-input v-model="mapForm.channel_uniq_id" disabled />
        </el-form-item>
        <el-form-item label="Jaco错误码">
          <el-select v-model="mapForm.jaco_uniq_id" placeholder="请选择Jaco错误码">
            <el-option
              v-for="item in jacoErrors"
              :key="item.uniq_id"
              :label="`${item.code}-${item.msg}`"
              :value="item.uniq_id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="用户提示">
          <el-select v-model="mapForm.user_id" placeholder="请选择用户提示">
            <el-option
              v-for="prompt in userPrompts"
              :key="prompt.id"
              :label="prompt.msg"
              :value="prompt.msg"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="mapDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveErrorMap">确定</el-button>
        </span>
      </template>
    </el-dialog>
    
    <!-- 新增/编辑用户提示对话框 -->
    <el-dialog
      v-model="userPromptDialogVisible"
      :title="userPromptDialogTitle"
      width="500px"
    >
      <el-form :model="userPromptForm" label-width="120px">
        <el-form-item label="英文提示">
          <el-input v-model="userPromptForm.msg" placeholder="请输入英文提示" type="textarea" />
        </el-form-item>
        <el-form-item label="阿语提示">
          <el-input v-model="userPromptForm.msg_ar" placeholder="请输入阿语提示" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="userPromptDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveUserPrompt">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { channelErrorApi, jacoErrorApi, errorMapApi, userPromptApi } from '@/api/errCodeApi'

export default {
  name: 'ErrorCodeManager',
  setup() {
    const router = useRouter()
    
    // 状态管理
    const activeTab = ref('map')
    
    // 渠道错误码
    const channelErrors = ref([])
    const channelSearch = ref('')
    const channelDialogVisible = ref(false)
    const channelDialogTitle = ref('新增渠道错误码')
    const channelForm = reactive({
      id: 0,
      channel: '',
      code: '',
      msg: ''
    })
    
    // 渠道筛选选项
    const channelFilters = ref([])
    
    // 筛选渠道
    const filterChannel = (value, row) => {
      return row.channel === value
    }
    
    // 筛选错误码
    const filterCode = (value, row) => {
      return row.code.includes(value)
    }
    
    // 筛选错误信息
    const filterMsg = (value, row) => {
      return row.msg.includes(value)
    }
    
    // Jaco错误码
    const jacoErrors = ref([])
    const jacoDialogVisible = ref(false)
    const jacoDialogTitle = ref('新增Jaco错误码')
    const jacoForm = reactive({
      id: 0,
      code: '',
      msg: '',
      desc: '',
      user_id: ''
    })
    
    // 筛选Jaco错误码
    const filterJacoCode = (value, row) => {
      return row.code.includes(value)
    }
    
    // 筛选Jaco错误信息
    const filterJacoMsg = (value, row) => {
      return row.msg.includes(value)
    }
    
    // 筛选Jaco描述
    const filterJacoDesc = (value, row) => {
      return row.desc.includes(value)
    }
    
    // 筛选Jaco英文提示
    const filterJacoDisplayMsg = (value, row) => {
      return row.display_msg.includes(value)
    }
    
    // 筛选Jaco阿语提示
    const filterJacoDisplayMsgAr = (value, row) => {
      return row.display_msg_ar.includes(value)
    }
    
    // 错误码映射
    const errorMapDetails = ref([])
    const mapSearchChannel = ref('')
    const mapDialogVisible = ref(false)
    const mapForm = reactive({
      id: 0,
      channel: '',
      channel_code: '',
      channel_msg: '',
      jaco_id: '',
      jaco_uniq_id: '',
      channel_uniq_id: '',
      user_id: ''
    })
    
    // 错误码映射渠道筛选选项
    const mapChannelFilters = ref([])
    
    // 用户提示
    const userPrompts = ref([])
    const userPromptDialogVisible = ref(false)
    const userPromptDialogTitle = ref('新增用户提示')
    const userPromptForm = reactive({
      id: 0,
      msg: '',
      msg_ar: ''
    })
    
    // 筛选映射渠道
    const filterMapChannel = (value, row) => {
      return row.channel === value
    }
    
    // 筛选映射渠道错误码
    const filterMapChannelCode = (value, row) => {
      return row.channel_code.includes(value)
    }
    
    // 筛选映射渠道错误信息
    const filterMapChannelMsg = (value, row) => {
      return row.channel_msg.includes(value)
    }
    
    // 筛选映射Jaco错误码
    const filterMapJacoCode = (value, row) => {
      return row.jaco_code.includes(value)
    }
    
    // 筛选映射Jaco错误信息
    const filterMapJacoMsg = (value, row) => {
      return row.jaco_msg.includes(value)
    }
    
    // 筛选映射英文提示
    const filterMapDisplayMsg = (value, row) => {
      return row.display_msg.includes(value)
    }
    
    // 筛选映射阿语提示
    const filterMapDisplayMsgAr = (value, row) => {
      return row.display_msg_ar.includes(value)
    }
    
    // 过滤Jaco错误码
    const filterJacoError = (query, item) => {
      return item.label.toLowerCase().includes(query.toLowerCase()) || 
             item.value.toLowerCase().includes(query.toLowerCase())
    }
    
    // 查询渠道错误码
    const searchChannelErrors = async () => {
      try {
        const data = await channelErrorApi.getChannelErrorList({ channel: channelSearch.value })
        if (data.code === 10000) {
          channelErrors.value = data.data
          
          // 提取唯一渠道值并更新筛选选项
          const uniqueChannels = [...new Set(data.data.map(item => item.channel))]
          channelFilters.value = uniqueChannels.map(channel => ({
            text: channel,
            value: channel
          }))
        } else {
          ElMessage.error(data.msg || '查询失败')
        }
      } catch (error) {
        ElMessage.error('查询失败：' + (error.msg || error.message))
      }
    }
    
    // 查询Jaco错误码
    const searchJacoErrors = async () => {
      try {
        const data = await jacoErrorApi.getJacoErrorList()
        if (data.code === 10000) {
          jacoErrors.value = data.data
        } else {
          ElMessage.error(data.msg || '查询失败')
        }
      } catch (error) {
        ElMessage.error('查询失败：' + (error.msg || error.message))
      }
    }
    
    // 查询错误码映射
    const searchErrorMaps = async () => {
      try {
        const data = await errorMapApi.getErrorMapDetails({ channel: mapSearchChannel.value })
        if (data.code === 10000) {
          errorMapDetails.value = data.data
          
          // 提取唯一渠道值并更新筛选选项
          const uniqueChannels = [...new Set(data.data.map(item => item.channel))]
          mapChannelFilters.value = uniqueChannels.map(channel => ({
            text: channel,
            value: channel
          }))
        } else {
          ElMessage.error(data.msg || '查询失败')
        }
      } catch (error) {
        ElMessage.error('查询失败：' + (error.msg || error.message))
      }
    }
    
    // 新增渠道错误码
    const showAddChannelDialog = () => {
      channelDialogTitle.value = '新增渠道错误码'
      Object.assign(channelForm, {
        id: 0,
        channel: '',
        code: '',
        msg: ''
      })
      channelDialogVisible.value = true
    }
    
    // 编辑渠道错误码
    const showEditChannelDialog = (row) => {
      channelDialogTitle.value = '编辑渠道错误码'
      Object.assign(channelForm, row)
      channelDialogVisible.value = true
    }
    
    // 保存渠道错误码
    const saveChannelError = async () => {
      try {
        let data
        if (channelForm.id) {
          data = await channelErrorApi.updateChannelError(channelForm)
        } else {
          data = await channelErrorApi.addChannelError(channelForm)
        }
        if (data.code === 10000) {
          ElMessage.success('保存成功')
          channelDialogVisible.value = false
          searchChannelErrors()
        } else {
          ElMessage.error(data.msg || '保存失败')
        }
      } catch (error) {
        ElMessage.error('保存失败：' + (error.msg || error.message))
      }
    }
    
    // 删除渠道错误码
    const deleteChannelError = async (id) => {
      try {
        const data = await channelErrorApi.deleteChannelError({ id })
        if (data.code === 10000) {
          ElMessage.success('删除成功')
          searchChannelErrors()
        } else {
          ElMessage.error(data.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败：' + (error.msg || error.message))
      }
    }
    
    // 新增Jaco错误码
    const showAddJacoDialog = () => {
      jacoDialogTitle.value = '新增Jaco错误码'
      Object.assign(jacoForm, {
        id: 0,
        code: '',
        msg: '',
        desc: '',
        user_id: ''
      })
      jacoDialogVisible.value = true
    }
    
    // 编辑Jaco错误码
    const showEditJacoDialog = (row) => {
      jacoDialogTitle.value = '编辑Jaco错误码'
      Object.assign(jacoForm, {
        id: row.id,
        code: row.code,
        msg: row.msg,
        desc: row.desc,
        user_id: row.user_id || ''
      })
      jacoDialogVisible.value = true
    }
    
    // 保存Jaco错误码
    const saveJacoError = async () => {
      try {
        let data
        if (jacoForm.id) {
          data = await jacoErrorApi.updateJacoError(jacoForm)
        } else {
          data = await jacoErrorApi.addJacoError(jacoForm)
        }
        if (data.code === 10000) {
          ElMessage.success('保存成功')
          jacoDialogVisible.value = false
          searchJacoErrors()
        } else {
          ElMessage.error(data.msg || '保存失败')
        }
      } catch (error) {
        ElMessage.error('保存失败：' + (error.msg || error.message))
      }
    }
    
    // 删除Jaco错误码
    const deleteJacoError = async (id) => {
      try {
        const data = await jacoErrorApi.deleteJacoError({ id })
        if (data.code === 10000) {
          ElMessage.success('删除成功')
          searchJacoErrors()
        } else {
          ElMessage.error(data.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败：' + (error.msg || error.message))
      }
    }
    
    // 编辑错误码映射
    const showEditMapDialog = (row) => {
      Object.assign(mapForm, {
        id: row.id,
        channel: row.channel,
        channel_code: row.channel_code,
        channel_msg: row.channel_msg,
        channel_uniq_id: row.channel_uniq_id,
        jaco_uniq_id: row.jaco_uniq_id || '',
        user_id: row.user_id || ''
      })
      mapDialogVisible.value = true
    }
    
    // 保存错误码映射
    const saveErrorMap = async () => {
      try {
        const data = await errorMapApi.updateErrorMap(mapForm)
        if (data.code === 10000) {
          ElMessage.success('保存成功')
          mapDialogVisible.value = false
          searchErrorMaps()
        } else {
          ElMessage.error(data.msg || '保存失败')
        }
      } catch (error) {
        ElMessage.error('保存失败：' + (error.msg || error.message))
      }
    }
    
    // 删除错误码映射
    const deleteErrorMap = async (id) => {
      try {
        const data = await errorMapApi.deleteErrorMap({ id })
        if (data.code === 10000) {
          ElMessage.success('删除成功')
          searchErrorMaps()
        } else {
          ElMessage.error(data.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败：' + (error.msg || error.message))
      }
    }
    
    // 导出错误码映射
    const exportErrorMaps = () => {
      if (errorMapDetails.value.length === 0) {
        ElMessage.warning('没有数据可以导出')
        return
      }
      
      // CSV 头部
      const headers = ['渠道', '渠道错误码', '渠道错误信息', 'Jaco错误码', 'Jaco错误信息', '英文提示', '阿语提示']
      
      // CSV 内容
      const rows = errorMapDetails.value.map(row => [
        row.channel || '',
        row.channel_code || '',
        row.channel_msg || '',
        row.jaco_code || '',
        row.jaco_msg || '',
        row.user_prompt_msg || '',
        row.user_prompt_msg_ar || ''
      ])
      
      // 转换为 CSV 字符串
      const csvContent = [
        headers.join(','),
        ...rows.map(row => row.map(cell => `"${cell}"`).join(','))
      ].join('\n')
      
      // 创建 Blob
      const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
      
      // 下载文件
      const link = document.createElement('a')
      if (link.download !== undefined) {
        const url = URL.createObjectURL(blob)
        link.setAttribute('href', url)
        link.setAttribute('download', `错误码映射_${new Date().toISOString().slice(0, 10)}.csv`)
        link.style.visibility = 'hidden'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      }
    }
    
    // 查询用户提示
    const searchUserPrompts = async () => {
      try {
        const data = await userPromptApi.getUserPromptList()
        if (data.code === 10000) {
          userPrompts.value = data.data
        } else {
          ElMessage.error(data.msg || '查询失败')
        }
      } catch (error) {
        ElMessage.error('查询失败：' + (error.msg || error.message))
      }
    }
    
    // 根据用户提示ID获取英文提示
    const getUserPromptMsg = (user_id) => {
      if (!user_id) return ''
      const prompt = userPrompts.value.find(item => item.msg == user_id)
      return prompt ? prompt.msg : ''
    }
    
    // 根据用户提示ID获取阿语提示
    const getUserPromptMsgAr = (user_id) => {
      if (!user_id) return ''
      const prompt = userPrompts.value.find(item => item.msg == user_id)
      return prompt ? prompt.msg_ar : ''
    }
    
    // 新增用户提示
    const showAddUserPromptDialog = () => {
      userPromptDialogTitle.value = '新增用户提示'
      Object.assign(userPromptForm, {
        id: 0,
        msg: '',
        msg_ar: ''
      })
      userPromptDialogVisible.value = true
    }
    
    // 编辑用户提示
    const showEditUserPromptDialog = (row) => {
      userPromptDialogTitle.value = '编辑用户提示'
      Object.assign(userPromptForm, row)
      userPromptDialogVisible.value = true
    }
    
    // 保存用户提示
    const saveUserPrompt = async () => {
      try {
        let data
        if (userPromptForm.id) {
          data = await userPromptApi.updateUserPrompt(userPromptForm)
        } else {
          data = await userPromptApi.addUserPrompt(userPromptForm)
        }
        if (data.code === 10000) {
          ElMessage.success('保存成功')
          userPromptDialogVisible.value = false
          searchUserPrompts()
        } else {
          ElMessage.error(data.msg || '保存失败')
        }
      } catch (error) {
        ElMessage.error('保存失败：' + (error.msg || error.message))
      }
    }
    
    // 删除用户提示
    const deleteUserPrompt = async (id) => {
      try {
        const data = await userPromptApi.deleteUserPrompt({ id })
        if (data.code === 10000) {
          ElMessage.success('删除成功')
          searchUserPrompts()
        } else {
          ElMessage.error(data.msg || '删除失败')
        }
      } catch (error) {
        ElMessage.error('删除失败：' + (error.msg || error.message))
      }
    }
    
    // 初始化数据
    onMounted(() => {
      searchChannelErrors()
      searchJacoErrors()
      searchErrorMaps()
      searchUserPrompts()
    })
    
    return {
      activeTab,
      
      // 渠道错误码
      channelErrors,
      channelSearch,
      channelDialogVisible,
      channelDialogTitle,
      channelForm,
      channelFilters,
      searchChannelErrors,
      showAddChannelDialog,
      showEditChannelDialog,
      saveChannelError,
      deleteChannelError,
      filterChannel,
      filterCode,
      filterMsg,
      
      // Jaco错误码
      jacoErrors,
      jacoDialogVisible,
      jacoDialogTitle,
      jacoForm,
      showAddJacoDialog,
      showEditJacoDialog,
      saveJacoError,
      deleteJacoError,
      filterJacoCode,
      filterJacoMsg,
      filterJacoDesc,
      filterJacoDisplayMsg,
      filterJacoDisplayMsgAr,
      
      // 错误码映射
      errorMapDetails,
      mapSearchChannel,
      mapDialogVisible,
      mapForm,
      mapChannelFilters,
      searchErrorMaps,
      showEditMapDialog,
      saveErrorMap,
      deleteErrorMap,
      exportErrorMaps,
      filterMapChannel,
      filterMapChannelCode,
      filterMapChannelMsg,
      filterMapJacoCode,
      filterMapJacoMsg,
      filterMapDisplayMsg,
      filterMapDisplayMsgAr,
      
      // 用户提示
      userPrompts,
      userPromptDialogVisible,
      userPromptDialogTitle,
      userPromptForm,
      searchUserPrompts,
      showAddUserPromptDialog,
      showEditUserPromptDialog,
      saveUserPrompt,
      deleteUserPrompt,
      getUserPromptMsg,
      getUserPromptMsgAr,
      
      // 过滤函数
      filterJacoError
    }
  }
}
</script>

<style scoped>
.error-code-manager {
  padding: 20px;
}

.tab-content {
  margin-top: 20px;
}

.mb-4 {
  margin-bottom: 16px;
}
</style>