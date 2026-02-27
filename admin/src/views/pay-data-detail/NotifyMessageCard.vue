<template>
  <el-card class="inner-card">

    <template #header>
      <div class="card-header">
        <div class="title-with-toggle" :class="{ collapsed: isCollapsed }" @click="toggleCollapse">
          <el-icon class="toggle-icon"><svg v-if="!isCollapsed" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"
              width="16" height="16">
              <path fill="currentColor"
                d="M8 12a.5.5 0 0 0 .5-.5V5.707l2.146 2.147a.5.5 0 0 0 .708-.708l-3-3a.5.5 0 0 0-.708 0l-3 3a.5.5 0 1 0 .708.708L7.5 5.707V11.5a.5.5 0 0 0 .5.5z" />
            </svg><svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16" width="16" height="16">
              <path fill="currentColor"
                d="M4.5 11.5a.5.5 0 0 0 .5.5h8.793l-2.147 2.146a.5.5 0 0 0 .708.708l3-3a.5.5 0 0 0 0-.708l-3-3a.5.5 0 0 0-.708.708L13.293 11H5a.5.5 0 0 0-.5.5z" />
            </svg></el-icon>
          <span>通知消息 #{{ msgIndex + 1 }}</span>
        </div>
        <el-button type="primary" size="small" @click="viewNotifyMessageJson">查看详情</el-button>
      </div>
    </template>

    <div v-if="!isCollapsed">
      <div class="info-grid">
        <div class="info-item">
          <label>消息ID：</label>
          <span>{{ notifyMessage.id }}</span>
        </div>
        <div class="info-item">
          <label>业务ID：</label>
          <span>{{ notifyMessage.biz_id }}</span>
        </div>
        <div class="info-item">
          <label>发送方：</label>
          <span>{{ notifyMessage.msg_from }}</span>
        </div>
        <div class="info-item">
          <label>接收方：</label>
          <span>{{ notifyMessage.msg_to }}</span>
        </div>
        <div class="info-item">
          <label>发送状态：</label>
          <span :class="getStatusClass(notifyMessage.send_status)">{{ notifyMessage.send_status }}</span>
        </div>
        <div class="info-item">
          <label>创建时间：</label>
          <span>{{ notifyMessage.created_at }}</span>
        </div>
      </div>
      <div class="info-grid">
        <div class="info-item">
          <label>消息内容：</label>
          <div class="expandable-text"
            @click="toggleLogExpand(`pay_notify_${paymentIndex}_${msgIndex}`, notifyMessage.msg_body)">
            <div v-if="expandedLogs.includes(`pay_notify_${paymentIndex}_${msgIndex}`)">
              <pre class="formatted-json">{{ formattedResponses[`pay_notify_${paymentIndex}_${msgIndex}`] }}</pre>
            </div>
            <div v-else>
              {{ truncateText(notifyMessage.msg_body, 50) }}
            </div>
            <span class="expand-hint">{{ notifyMessage.msg_body && notifyMessage.msg_body.length > 50 ?
              (expandedLogs.includes(`pay_notify_${paymentIndex}_${msgIndex}`) ? '点击收起' : '点击展开') : '' }}</span>
          </div>
        </div>
      </div>
      <!-- </div> -->
    </div>

    <!-- 通知历史记录 -->
    <!-- <div v-if="notifyHistories.length > 0">
      <div class="title-with-toggle">
        <h5 class="subsection-title">通知历史记录 ({{ notifyHistories.length }})</h5>
      </div>
      <notify-history-card v-for="(notifyHistory, historyIndex) in notifyHistories"
        :key="notifyHistory.id || historyIndex" :notify-history="notifyHistory" :payment-index="paymentIndex"
        :msg-index="msgIndex" :history-index="historyIndex" :is-collapsed="isNotifyHistoryCollapsed(historyIndex)"
        :expanded-logs="expandedLogs" :formatted-responses="formattedResponses" @toggle-collapse="toggleNotifyHistory"
        @view-json="viewNotifyHistoryJson" @toggle-log-expand="toggleLogExpand" />
    </div> -->
    <!-- </div> -->
  </el-card>
</template>

<script>
import NotifyHistoryCard from './NotifyHistoryCard.vue';
import { ElIcon } from 'element-plus';

export default {
  name: 'NotifyMessageCard',
  components: {
    NotifyHistoryCard,
    ElIcon
  },
  data() {
    return {
      // 移除内部状态管理，改为通过父组件管理折叠状态
    };
  },
  props: {
    notifyMessage: {
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
    notifyHistories: {
      type: Array,
      default: () => []
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
  mounted() {
    // 组件挂载时的逻辑
    console.log('NotifyMessageCard组件挂载 - 消息ID:', this.notifyMessage.id);
    console.log('NotifyMessageCard组件挂载 - 通知历史数量:', this.notifyHistories.length);
    console.log('NotifyMessageCard组件挂载 - isCollapsed状态:', this.isCollapsed);
  },
  methods: {
    getStatusClass(status) {
      if (!status) return '';

      const statusMap = {
        'SUCCESS': 'status-success',
        'FAILED': 'status-fail',
        'PENDING': 'status-pending',
      };

      return statusMap[status.toUpperCase()] || '';
    },
    truncateText(text, maxLength) {
      if (!text) return '';
      if (text.length <= maxLength) return text;
      return text.substring(0, maxLength) + '...';
    },
    toggleCollapse() {
      this.$emit('toggle-collapse', this.paymentIndex, this.msgIndex);
    },
    viewNotifyMessageJson() {
      this.$emit('view-json', this.notifyMessage);
    },
    toggleNotifyHistory(paymentIndex, msgIndex, historyIndex) {
      // 直接通知父组件切换折叠状态，不再在组件内部管理状态
      this.$emit('toggle-history-collapse', paymentIndex, msgIndex, historyIndex);
    },
    isNotifyHistoryCollapsed(historyIndex) {
      // 尝试从父组件获取折叠状态
      let collapsed = false;
      if (this.$parent && typeof this.$parent.isNotifyHistoryCollapsed === 'function') {
        collapsed = this.$parent.isNotifyHistoryCollapsed(this.paymentIndex, this.msgIndex, historyIndex);
      }
      // 如果父组件没有提供，则尝试从父组件的父组件获取
      else if (this.$parent && this.$parent.$parent && typeof this.$parent.$parent.isNotifyHistoryCollapsed === 'function') {
        collapsed = this.$parent.$parent.isNotifyHistoryCollapsed(this.paymentIndex, this.msgIndex, historyIndex);
      }
      // 默认展开
      console.log('NotifyMessageCard - 历史记录折叠状态:', collapsed, '历史索引:', historyIndex);
      return collapsed;
    },
    viewNotifyHistoryJson(notifyHistory) {
      this.$emit('view-json', notifyHistory);
    },
    toggleLogExpand(index, response) {
      // 直接向父组件发射事件
      this.$emit('toggle-log-expand', index, response);
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

.subsection-title {
  margin: 20px 0 10px 0;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  padding: 6px 12px;
  border-radius: 6px;
  transition: all 0.2s;
  background-color: #f0f9ff;
  border-left: 4px solid #409EFF;
}

.subsection-title:hover {
  color: #409EFF;
  background-color: #ecf5ff;
  border-left-color: #66b1ff;
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