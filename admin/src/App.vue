<template>
  <el-container style="height: 100vh; border: 1px solid #eee;">
    <el-aside :width="isSidebarCollapse ? '60px' : '200px'" style="background-color: #304156; color: white; transition: width 0.3s ease;">
      <div style="padding: 20px; font-size: 18px; font-weight: bold; display: flex; justify-content: space-between; align-items: center;">
        <span v-show="!isSidebarCollapse">管理后台</span>
        <el-button 
          type="text" 
          size="small" 
          @click="toggleSidebar"
          style="color: white; padding: 0; min-width: 0;"
        >
          <el-icon v-if="!isSidebarCollapse"><Fold /></el-icon>
          <el-icon v-else><Expand /></el-icon>
        </el-button>
      </div>
      <el-menu 
        :default-openeds="['1']" 
        background-color="#304156" 
        text-color="#fff" 
        active-text-color="#ffd04b" 
        router
      >
        <!-- <el-menu-item index="/tasks">
          <el-icon><Document /></el-icon>
          <template #title>
            <span v-show="!isSidebarCollapse">调度任务管理</span>
          </template>
        </el-menu-item> -->
        <el-menu-item index="/file-processor">
          <el-icon><Document /></el-icon>
          <template #title>
            <span v-show="!isSidebarCollapse">文件处理</span>
          </template>
        </el-menu-item>
    
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header style="background-color: #fff; padding: 0;">
        <div style="height: 60px; line-height: 60px; padding: 0 20px; border-bottom: 1px solid #eee;">
          <span style="font-size: 16px; font-weight: 500;">{{ currentPath }}</span>
        </div>
      </el-header>
      
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script>
import { ref, onMounted, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Document, Setting, Fold, Expand, Search, Coin, Key } from '@element-plus/icons-vue'
import { useAppStore } from '@/store/index.js'

export default {
  name: 'App',
  components: {
    Document,
    Setting,
    Fold,
    Expand,
    Search,
    Coin,
    Key
  },
  setup() {
    const route = useRoute()
    const currentPath = ref('')
    const appStore = useAppStore()
    const isSidebarCollapse = computed(() => appStore.isSidebarCollapse)
    
    // 判断是否为测试环境
    const isTestEnvironment = import.meta.env.VITE_API_URL != 'http://pay-scheduler.alphicloud.com'
    
    const toggleSidebar = () => {
      appStore.toggleSidebar()
    }
    
    const updateCurrentPath = () => {
      if (route.path === '/tasks' || route.path.startsWith('/tasks/')) {
        currentPath.value = '调度任务管理'
      } else if (route.path === '/file-processor') {
        currentPath.value = '文件处理'
      }
    
    }
    
    onMounted(() => {
      updateCurrentPath()
    })
    
    watch(() => route.path, () => {
      updateCurrentPath()
    })
    
    return {
      currentPath,
      isSidebarCollapse,
      toggleSidebar,
      isTestEnvironment
    }
  }
}
</script>

<style scoped>
.el-header {
  padding: 0;
}

.el-main {
  padding: 20px;
  background-color: #f5f5f5;
}

.el-container {
  height: 100vh;
}
</style>