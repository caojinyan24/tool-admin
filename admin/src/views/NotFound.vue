<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <el-card class="box-card" style="width: 100%; max-width: 500px;">
        <template #header>
          <div class="card-header">
            <span>404 - 页面不存在</span>
          </div>
        </template>
        <div class="not-found-body">
          <div class="error-code">404</div>
          <div class="error-message">很抱歉，您访问的页面不存在或已被删除</div>
          <div class="error-actions">
            <el-button type="primary" @click="goToHome">
              <el-icon><House /></el-icon>
              <span>返回首页</span>
            </el-button>
            <el-button @click="goBack">
              <el-icon><Back /></el-icon>
              <span>返回上一页</span>
            </el-button>
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { House, Back } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'

export default {
  name: 'NotFound',
  components: {
    House,
    Back
  },
  setup() {
    const router = useRouter()
    
    const goToHome = () => {
      router.push('/')
    }
    
    const goBack = () => {
      router.back()
    }
    
    return {
      goToHome,
      goBack
    }
  }
}
</script>

<style scoped>
.not-found-container {
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  padding: 20px;
}

.not-found-content {
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 18px;
  font-weight: 500;
}

.not-found-body {
  text-align: center;
  padding: 40px 0;
}

.error-code {
  font-size: 72px;
  font-weight: bold;
  color: #e74c3c;
  margin-bottom: 20px;
}

.error-message {
  font-size: 16px;
  color: #666;
  margin-bottom: 30px;
}

.error-actions {
  display: flex;
  justify-content: center;
  gap: 20px;
}
</style>