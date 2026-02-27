const express = require('express');
const path = require('path');
const port = process.env.PORT || 8080;
const app = express();

// 提供静态文件服务
app.use(express.static(path.join(__dirname, 'dist')));

// SPA fallback - 所有非静态资源请求都返回index.html
app.get('*', (req, res) => {
  res.sendFile(path.resolve(__dirname, 'dist', 'index.html'));
});

// 启动服务器
app.listen(port, () => {
  console.log(`Server running on port ${port}`);
  console.log(`To access the application, navigate to http://localhost:${port}`);
});