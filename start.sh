#!/bin/bash

# 启动脚本：并行启动前端和后端服务，并确保两个服务都成功启动

# 设置颜色变量
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # 无颜色

# 清理旧的日志文件
rm -f frontend.log backend.log
rm -f  server/nohup.out


# 记录开始时间
START_TIME=$(date +%s)

# 定义检查后端服务是否启动的函数
check_backend() {
    local max_attempts=30
    local attempt=0
    local delay=2
    
    echo -e "${YELLOW}等待后端服务启动...${NC}"
    
    while [ $attempt -lt $max_attempts ]; do
        # 检查后端是否有日志输出，匹配Hertz框架的实际日志格式
        if grep -q "HERTZ: Using network library" backend.log || curl -s http://localhost:8080/api/v1/health > /dev/null 2>&1; then
            echo -e "${GREEN}后端服务启动成功！${NC}"
            return 0
        fi
        
        # 检查后端是否有错误输出
        if grep -q "Fatal" backend.log || grep -q "Error" backend.log || grep -q "failed" backend.log || grep -q "panic" backend.log; then
            echo -e "${RED}后端服务启动失败，请查看backend.log了解详情${NC}"
            return 1
        fi
        
        attempt=$((attempt+1))
        echo -e "${YELLOW}尝试 $attempt/$max_attempts: 后端服务尚未启动...${NC}"
        sleep $delay
    done
    
    echo -e "${RED}后端服务启动超时！${NC}"
    return 1
}

# 定义检查前端服务是否启动的函数
check_frontend() {
    local max_attempts=30
    local attempt=0
    local delay=2
    
    echo -e "${YELLOW}等待前端服务启动...${NC}"
    
    while [ $attempt -lt $max_attempts ]; do
        # 检查前端是否有日志输出，这里假设前端服务启动成功会有特定日志
        if grep -q "VITE" frontend.log && (grep -q "ready in" frontend.log || grep -q "running at" frontend.log); then
            # 尝试通过HTTP请求验证前端服务是否真的可用
            if curl -s http://localhost:3000 > /dev/null 2>&1; then
                echo -e "${GREEN}前端服务启动成功！${NC}"
                return 0
            fi
        fi
        
        # 检查前端是否有严重错误输出，但忽略代理错误
        if grep -q "error" frontend.log && ! grep -q "http proxy error" frontend.log; then
            echo -e "${RED}前端服务启动失败，请查看frontend.log了解详情${NC}"
            return 1
        fi
        
        attempt=$((attempt+1))
        echo -e "${YELLOW}尝试 $attempt/$max_attempts: 前端服务尚未启动...${NC}"
        sleep $delay
    done
    
    echo -e "${RED}前端服务启动超时！${NC}"
    return 1
}

# 启动后端服务（Go服务）
cd server || { echo -e "${RED}无法进入server目录${NC}"; exit 1; }
# 注意：这里使用了2>&1来重定向标准错误到标准输出，然后写入日志文件
export StartEnv=test
nohup go build -o server && ./server > ../backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# 先等待后端服务完全启动
check_backend
if [ $? -ne 0 ]; then
    echo -e "${RED}后端服务启动失败，取消启动前端服务${NC}"
    kill $BACKEND_PID 2>/dev/null
    exit 1
fi

# 后端服务启动成功后，再启动前端服务
cd admin || { echo -e "${RED}无法进入admin目录${NC}"; kill $BACKEND_PID; exit 1; }
# 注意：这里使用了2>&1来重定向标准错误到标准输出，然后写入日志文件
npm install
nohup npm run serve > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# 显示服务启动信息
cat << EOF

${YELLOW}服务启动信息${NC}
--------------------------------------------------
后端服务PID: $BACKEND_PID
前端服务PID: $FRONTEND_PID
后端日志: backend.log
前端日志: frontend.log
--------------------------------------------------
EOF

# 检查前端服务是否启动成功
frontend_ready=0
check_frontend && frontend_ready=1

# 计算总启动时间
END_TIME=$(date +%s)
TOTAL_TIME=$((END_TIME-START_TIME))

# 检查两个服务是否都启动成功
if [ $frontend_ready -eq 1 ]; then
    echo -e "\n${GREEN}==================================================${NC}"
    echo -e "${GREEN}🎉 恭喜！所有服务已成功启动！${NC}"
    echo -e "${GREEN}🚀 总启动时间: ${TOTAL_TIME}秒${NC}"
    echo -e "${GREEN}🌐 前端服务: http://localhost:3000/ (或3001，如果3000被占用)${NC}"
    echo -e "${GREEN}⚙️  后端服务: http://localhost:8080/ (或根据配置)${NC}"
    echo -e "${GREEN}==================================================${NC}"
    
    # 创建一个函数来清理后台进程
    cleanup() {
        echo -e "\n${YELLOW}正在停止所有服务...${NC}"
        kill $BACKEND_PID 2>/dev/null
        kill $FRONTEND_PID 2>/dev/null
        echo -e "${GREEN}所有服务已停止${NC}"
    }
    
    # 捕获SIGINT和SIGTERM信号，用于优雅退出
    trap cleanup SIGINT SIGTERM
    
    # 等待用户输入Ctrl+C来停止服务
    echo -e "\n${YELLOW}按 Ctrl+C 停止所有服务${NC}"
    wait  # 等待直到收到终止信号
else
    echo -e "\n${RED}==================================================${NC}"
    echo -e "${RED}❌ 启动失败！前端服务未能成功启动。${NC}"
    echo -e "${RED}请查看日志文件了解详细信息：${NC}"
    echo -e "${RED}  - backend.log: 后端服务日志${NC}"
    echo -e "${RED}  - frontend.log: 前端服务日志${NC}"
    echo -e "${RED}==================================================${NC}"
    
    # 清理失败的进程
    kill $BACKEND_PID 2>/dev/null
    kill $FRONTEND_PID 2>/dev/null
    
    exit 1
fi


