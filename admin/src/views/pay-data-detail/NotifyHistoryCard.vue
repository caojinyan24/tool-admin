<template>
  <el-card class="small-card">
    <template #header>
      <div class="card-header">
        <div class="title-with-toggle" :class="{ collapsed: isCollapsed }" @click="toggleCollapse">
          <el-icon class="toggle-icon"><ChevronDown v-if="!isCollapsed" /><ChevronRight v-else /></el-icon>
          <span>历史记录 #{{ historyIndex + 1 }}</span>
        </div>
        <el-button type="primary" size="small" @click="viewNotifyHistoryJson">查看详情</el-button>
      </div>
    </template>
    
    <div v-if="!isCollapsed">
      <div class="info-grid">
        <div class="info-item">
          <label>返回码：</label>
          <span>{{ notifyHistory.ret_code }}</span>
        </div>
        <div class="info-item">
          <label>创建时间：</label>
          <span>{{ notifyHistory.created_at }}</span>
        </div>
        <div class="info-item">
          <label>返回体：</label>
          <div class="expandable-text" @click="toggleLogExpand(`pay_history_${paymentIndex}_${msgIndex}_${historyIndex}`, notifyHistory.ret_body)">
            <div v-if="expandedLogs.includes(`pay_history_${paymentIndex}_${msgIndex}_${historyIndex}`)">
              <pre class="formatted-json">{{ formattedResponses[`pay_history_${paymentIndex}_${msgIndex}_${historyIndex}`] }}</pre>
            </div>
            <div v-else>
              {{ truncateText(notifyHistory.ret_body, 50) }}
            </div>
            <span class="expand-hint">{{ notifyHistory.ret_body && notifyHistory.ret_body.length > 50 ? (expandedLogs.includes(`pay_history_${paymentIndex}_${msgIndex}_${historyIndex}`) ? '点击收起' : '点击展开') : '' }}</span>
          </div>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script>
import * as Icons from '@element-plus/icons-vue';

export default {
  name: 'NotifyHistoryCard',
  components: {
    ChevronDown: Icons.ArrowDown,
    ChevronRight: Icons.ArrowRight
  },
  props: {
    notifyHistory: {
      type: Object,
      required: true
    },
    paymentIndex: {
      type: Number,
      required: true
    },
    msgIndex: {
      type: Number,
      required: true
    },
    historyIndex: {
      type: Number,
      required: true
    },
    isCollapsed: {
      type: Boolean,
      default: false
    },
    expandedLogs: {
      type: Array,
      default: () => []
    },
    formattedResponses: {
      type: Object,
      default: () => ({})
    }
  },
  methods: {
    truncateText(text, maxLength) {
      if (!text) return '';
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + '...';
    },
    toggleCollapse() {
      this.$emit('toggle-collapse', this.paymentIndex, this.msgIndex, this.historyIndex);
    },
    viewNotifyHistoryJson() {
      this.$emit('view-json', this.notifyHistory);
    },
    toggleLogExpand(index, response) {
      this.$emit('toggle-log-expand', index, response);
    }
  }
};
</script>

<style scoped>
.small-card {
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