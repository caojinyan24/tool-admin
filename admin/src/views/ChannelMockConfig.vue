<template>
  <div class="channel-mock-config">
    <!-- 操作指引按钮 -->
    <!-- <div class="operation-guide-toggle" @click="toggleGuide" :class="{ 'active': showGuide }">
      <el-button type="primary" plain icon="el-icon-question">
        操作指引
      </el-button>
    </div> -->
    
    <!-- 操作指引面板 -->
    <!-- <div class="operation-guide" v-if="showGuide">
      <el-card class="guide-card">
        <div class="guide-content">
          <!-- <el-collapse v-model="activeNames"> -->
            <!-- <el-collapse-item title="查看操作指引" name="1"> -->
              <!-- <div class="guide-step">
                <div class="step-number">1</div>
                <div class="step-content">
                  <h3>添加渠道响应配置</h3>
                  <p>在<span class="tab-link">渠道响应配置</span>标签页中，点击<span class="action-btn">新增</span>按钮，填写渠道ID、请求类型等信息，创建响应配置记录。</p>
                </div>
              </div>
              <div class="guide-step">
                <div class="step-number">2</div>
                <div class="step-content">
                  <h3>配置匹配规则</h3>
                  <p>切换到<span class="tab-link">规则配置</span>标签页，点击<span class="action-btn">新增</span>按钮，填写渠道ID、请求类型、用户标识等信息，配置过期时间，创建匹配规则。(若过期时间过期或规则状态为禁用, 则规则不生效, 即不会使用模拟的响应)</p>
                </div>
              </div> -->
            <!-- </el-collapse-item> -->
          <!-- </el-collapse> -->
        <!-- </div> -->
      <!-- </el-card> -->
    <!-- </div> -->
     
    <div class="main-content">
      <el-tabs v-model="activeTab" type="border-card" class="config-tabs" @tab-click="handleTabClick">
        <el-tab-pane label="规则配置" name="rule-config">
          <ChannelMockRuleConfig />
        </el-tab-pane>
        <el-tab-pane label="渠道响应配置" name="response-config">
          <ChannelMockResponseConfig />
        </el-tab-pane>
        <el-tab-pane label="场景管理" name="scene-management">
          <ChannelMockSceneManagement />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import ChannelMockResponseConfig from './channel-mock/ChannelMockResponseConfig.vue';
import ChannelMockRuleConfig from './channel-mock/ChannelMockRuleConfig.vue';
import ChannelMockSceneManagement from './channel-mock/ChannelMockSceneManagement.vue';

export default {
  name: 'ChannelMockConfig',
  components: {
    ChannelMockRuleConfig,
    ChannelMockResponseConfig,
    ChannelMockSceneManagement
  },
  data() {
    return {
      activeTab: 'rule-config',
      // 控制操作指引面板的显示/隐藏
      showGuide: false
    };
  },
  created() {
    // 检查路由参数，如果有tab参数则设置默认tab页
    if (this.$route.query.tab) {
      this.activeTab = this.$route.query.tab;
    }
  },
  methods: {
    // 切换操作指引面板的显示状态
    toggleGuide() {
      this.showGuide = !this.showGuide;
    },
    // 处理tab点击事件，确保切换tab时重新加载数据
    handleTabClick(tab) {
      // 当切换到不同tab时，子组件会自动触发created钩子，重新加载数据
    }
  }
};
</script>

<style scoped>
.channel-mock-config {
  padding: 20px;
  position: relative;
}

h1 {
  margin-bottom: 20px;
  color: #303133;
}

.main-content {
  width: 100%; /* 主内容区域恢复为100%宽度 */
}

.config-tabs {
  margin-top: 0;
}

/* 操作指引按钮样式 */
.operation-guide-toggle {
  /* position: fixed; */
  left: 20px;
  /* top: 100px; */
  z-index: 100;
}

.operation-guide-toggle .el-button {
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.15);
  white-space: nowrap;
}

/* 操作指引面板样式 */
.operation-guide {
  /* position: fixed; */
  /* right: 20px; */
  top: 150px;
  /* width: 280px; */
  z-index: 99;
  max-height: calc(100vh - 170px);
  overflow-y: auto;
}

.guide-card {
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.15);
}

.card-header {
  display: flex;
  align-items: center;
  font-weight: bold;
  font-size: 16px;
  color: #303133;
}

.guide-content {
  padding: 10px 0;
}

/* 折叠面板样式 */
:deep(.el-collapse) {
  border: none;
}

:deep(.el-collapse-item__header) {
  font-weight: bold;
  color: #303133;
  background-color: #f5f7fa;
  border-radius: 4px;
  padding: 10px 15px;
}

:deep(.el-collapse-item__content) {
  padding: 15px;
}

.guide-step {
  display: flex;
  margin-bottom: 20px;
  align-items: flex-start;
}

.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background-color: #409eff;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  margin-right: 16px;
  flex-shrink: 0;
}

.step-content h3 {
  margin: 0 0 8px 0;
  color: #303133;
  font-size: 14px;
}

.step-content p {
  margin: 0;
  color: #606266;
  font-size: 13px;
  line-height: 1.6;
}

.tab-link {
  color: #409eff;
  font-weight: bold;
}

.action-btn {
  color: #67c23a;
  font-weight: bold;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .operation-guide-toggle {
    position: static;
    margin-bottom: 20px;
  }
  
  .operation-guide {
    position: static;
    width: 100%;
    max-height: none;
    margin-bottom: 20px;
  }
  
  .main-content {
    width: 100%;
  }
}
</style>