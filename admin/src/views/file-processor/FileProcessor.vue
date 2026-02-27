<template>
  <div class="file-processor-container">
    <h2>文件处理</h2>
    <el-tabs v-model="activeTab">
      <el-tab-pane label="PDF合并" name="pdf-merge">
        <div class="pdf-merge-content">
          <el-upload
            class="upload-demo"
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :file-list="fileList"
            accept=".pdf"
            multiple
          >
            <el-button type="primary">选择PDF文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                请选择多个PDF文件进行合并
              </div>
            </template>
          </el-upload>
          
          <div 
            class="file-list-container" 
            style="margin-top: 20px;"
            ref="fileListContainer"
          >
            <div 
              v-for="(item, index) in fileList" 
              :key="item.uid"
              class="file-item"
              @click="previewPdf(item)"
              draggable="true"
              @dragstart="onDragStart($event, index)"
              @dragover.prevent
              @drop="onDrop($event, index)"
            >
              <el-icon class="file-icon"><Document /></el-icon>
              <span class="file-name">{{ item.name }}</span>
              <el-icon class="drag-icon"><Rank /></el-icon>
            </div>
          </div>
          
          <el-button 
            type="success" 
            :disabled="fileList.length < 2"
            @click="mergePdf"
            style="margin-top: 20px"
          >
            合并PDF
          </el-button>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import axios from 'axios';
import { Document, Rank } from '@element-plus/icons-vue';

export default {
  name: 'FileProcessor',
  components: {
    Document,
    Rank
  },
  data() {
    return {
      activeTab: 'pdf-merge',
      fileList: [],
      draggedIndex: -1
    }
  },
  methods: {
    handleFileChange(file, fileList) {
      // 只保留PDF文件
      this.fileList = fileList.filter(item => item.raw.type === 'application/pdf');
    },
    mergePdf() {
      // 创建FormData对象
      const formData = new FormData();
      
      // 添加文件到FormData
      this.fileList.forEach(file => {
        formData.append('files', file.raw);
      });
      
      // 调用后端API进行PDF合并
      axios.post('/api/v1/file-processor/merge-pdf', formData, {
        responseType: 'blob',
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      .then(response => {
        // 创建下载链接
        const url = window.URL.createObjectURL(new Blob([response.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'merged.pdf');
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
      })
      .catch(error => {
        console.error('合并PDF失败:', error);
        this.$message.error('合并PDF失败，请重试');
      });
    },
    previewPdf(file) {
      // 创建预览链接
      const url = window.URL.createObjectURL(file.raw);
      // 在新窗口中打开预览
      window.open(url, '_blank');
    },
    onDragStart(event, index) {
      this.draggedIndex = index;
      event.target.style.opacity = '0.5';
    },
    onDrop(event, index) {
      event.preventDefault();
      event.target.style.opacity = '';
      
      if (this.draggedIndex !== -1 && this.draggedIndex !== index) {
        // 移除拖动的元素
        const draggedItem = this.fileList.splice(this.draggedIndex, 1)[0];
        // 插入到新位置
        this.fileList.splice(index, 0, draggedItem);
        // 重置拖动索引
        this.draggedIndex = -1;
      }
    }
  }
}
</script>

<style scoped>
.file-processor-container {
  padding: 20px;
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.pdf-merge-content {
  padding: 20px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

h2 {
  margin-bottom: 20px;
  color: #303133;
}

.file-list-container {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 10px;
  min-height: 100px;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  margin-bottom: 10px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.file-item:hover {
  background-color: #f5f7fa;
  border-color: #c0c4cc;
}

.file-icon {
  margin-right: 10px;
  color: #409eff;
}

.file-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.drag-icon {
  margin-left: 10px;
  color: #c0c4cc;
  cursor: move;
}
</style>