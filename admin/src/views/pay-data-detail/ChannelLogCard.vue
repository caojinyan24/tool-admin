<template>
  <el-card class="inner-card">
    <template #header>
      <div class="card-header">
        <div class="title-with-toggle" :class="{ collapsed: isCollapsed }" @click="toggleCollapse">
          <el-icon class="toggle-icon"><ChevronDown v-if="!isCollapsed" /><ChevronRight v-else /></el-icon>
          <span>Channel Log #{{ logIndex + 1 }}</span>
        </div>
        <el-button type="primary" size="small" @click="viewChannelLogJson">查看详情</el-button>
      </div>
    </template>
    
    <div v-if="!isCollapsed">
      <div class="info-grid">
        <div class="info-item">
          <label>订单ID：</label>
          <span>{{ channelLog.order_id || 'N/A' }}</span>
        </div>
        <div class="info-item">
          <label>创建时间：</label>
          <span>{{ channelLog.created_at }}</span>
        </div>
        <div class="info-item event-type-highlight">
          <label>事件类型：</label>
          <span>{{ channelLog.event_type }}</span>
        </div>
      </div>
        <div class="info-grid">
        <div class="info-item">
          <label>响应内容：</label>
          <div class="expandable-text" @click="toggleLogExpand">
            <div v-if="isLogExpanded">
              <pre class="formatted-json">{{ formattedResponse }}</pre>
            </div>
            <div v-else>
              {{ truncateText(channelLog.response, 50) }}
            </div>
            <span class="expand-hint">{{ channelLog.response && channelLog.response.length > 50 ? (isLogExpanded ? '点击收起' : '点击展开') : '' }}</span>
          </div>
        </div>
        </div>
      </div>
    
  </el-card>
</template>

<script>
import * as Icons from '@element-plus/icons-vue';

export default {
  name: 'ChannelLogCard',
  components: {
    ChevronDown: Icons.ArrowDown,
    ChevronRight: Icons.ArrowRight
  },
  props: {
    channelLog: {
      type: Object,
      required: true
    },
    paymentIndex: {
      type: Number,
      required: true
    },
    coreIndex: {
      type: Number,
      required: true
    },
    logIndex: {
      type: Number,
      required: true
    },
    isCollapsed: {
      type: Boolean,
      default: false
    }
  },
  data() {
    return {
      isLogExpanded: false,
      formattedResponse: ''
    };
  },
  methods: {
    toggleCollapse() {
      this.$emit('toggle-collapse', this.paymentIndex, this.coreIndex, this.logIndex);
    },
    viewChannelLogJson() {
      this.$emit('view-json', this.channelLog);
    },
    truncateText(text, maxLength) {
      if (!text) return '';
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + '...';
    },
    toggleLogExpand() {
      this.isLogExpanded = !this.isLogExpanded;
      
      // 当展开时格式化响应内容
      if (this.isLogExpanded && this.channelLog.response) {
        try {
          // 尝试解析为JSON并格式化
          this.formattedResponse = JSON.stringify(JSON.parse(this.channelLog.response), null, 2);
        } catch (e) {
          // 如果不是有效的JSON，则使用原始内容
          this.formattedResponse = this.channelLog.response;
        }
      }
      
      // 仍然发射事件，以便父组件知道状态变化
      const index = `channel_log_${this.paymentIndex}_${this.coreIndex}_${this.logIndex}`;
      this.$emit('toggle-log-expand', index, this.channelLog.response);
    }
  }
};
</script>

<style scoped>
.inner-card {
  margin-left: 20px;
  margin-bottom: 10px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title-with-toggle {
  display: flex;
  align-items: center;
  cursor: pointer;
  user-select: none;
  padding: 4px 8px;
  border-radius: 4px;
  transition: all 0.2s;
}

.title-with-toggle:hover {
  color: #409EFF;
  background-color: #ecf5ff;
}

.toggle-icon {
  margin-right: 5px;
  font-size: 12px;
  transition: transform 0.2s;
}

.title-with-toggle.collapsed {
  font-weight: normal;
  color: #606266;
  background-color: #f4f4f5;
  border-left: 3px solid #dcdfe6;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 15px;
  margin-top: 10px;
}

.info-item {
  display: flex;
  align-items: flex-start;
}

.info-item label {
  font-weight: normal;
  margin-right: 10px;
  min-width: 100px;
  color: #909399;
  font-size: 14px;
}

.info-item span {
  color: #303133;
  word-break: break-all;
  font-size: 14px;
}

.expandable-text {
  cursor: pointer;
  color: #409EFF;
  display: inline-block;
  word-break: break-all;
  max-width: calc(100% - 110px);
}

.expandable-text:hover {
  text-decoration: underline;
}

.expand-hint {
  color: #909399;
  font-size: 12px;
  margin-left: 5px;
}

.formatted-json {
  white-space: pre-wrap;
  word-break: break-all;
  background: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
  margin: 5px 0 0 0;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
  line-height: 1.5;
  max-height: 400px;
  overflow: auto;
}
</style>