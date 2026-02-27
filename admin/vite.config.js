import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // 加载环境变量
  const env = loadEnv(mode, process.cwd())
  return {
    plugins: [vue()],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src')
      }
    },
    server: {
      port: 3000,
      proxy: {
        '/api': {
          // 只从环境变量中加载API URL
          target: env.VITE_API_URL ,
          changeOrigin: true,
          // 不需要重写路径，因为后端API已经有v1前缀
        },
        '/logs': {
          // 只从环境变量中加载API URL
          target: env.VITE_API_URL ,
          changeOrigin: true,
          // 不需要重写路径，直接转发到后端的/logs端点
        }
      },
      historyApiFallback: true
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            'vue-vendor': ['vue', 'vue-router', 'pinia', 'axios'],
            'element-plus': ['element-plus'],
            'icons': ['@element-plus/icons-vue']
          }
        }
      },
      base: './',
      outDir: 'dist',
      assetsDir: 'assets'
    },
    define: {
      'process.env': env
    }
  }
})