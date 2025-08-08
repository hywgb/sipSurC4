<<<<<<< Current (Your changes)
# sipSurC4
sip sur nankang cla4
=======
# CloudCallCenter - 云联络中心系统

基于云原生架构的现代化联络中心平台，融合了CATI（计算机辅助电话访谈）、智能客服、全渠道接入等功能。

## 项目特性

- 🚀 **高性能**：支持10000+座席并发，100000+并发呼叫
- 🔧 **微服务架构**：基于Go、Node.js、Python的微服务体系
- 🤖 **AI赋能**：集成语音识别、情感分析、智能质检
- 📊 **实时分析**：实时数据大屏，多维度统计分析
- 🔌 **插件化**：灵活的插件系统，支持功能扩展
- 🔐 **企业级安全**：OAuth2.0认证，数据加密，审计日志

## 系统架构

```
┌─────────────────────────────────────────────────────┐
│                   前端应用层                         │
│  Web管理端 | 座席工作台 | 移动端APP | 数据大屏      │
├─────────────────────────────────────────────────────┤
│                   API网关层                          │
│  Kong Gateway | 认证授权 | 限流熔断 | 负载均衡      │
├─────────────────────────────────────────────────────┤
│                  微服务层                            │
│  呼叫服务 | 问卷服务 | 配额服务 | 报表服务 | AI服务  │
├─────────────────────────────────────────────────────┤
│                  基础服务层                          │
│  消息队列 | 缓存服务 | 存储服务 | 日志服务 | 监控服务│
├─────────────────────────────────────────────────────┤
│                  基础设施层                          │
│  Kubernetes | Docker | 云服务(AWS/阿里云) | CI/CD   │
└─────────────────────────────────────────────────────┘
```

## 核心模块

### 后端服务
- **call-center-service**: 呼叫中心核心服务 (Go)
- **survey-service**: 问卷管理服务 (Node.js)
- **quota-service**: 配额管理服务 (Java)
- **analytics-service**: 数据分析服务 (Python)
- **ai-service**: AI能力服务 (Python)
- **agent-service**: 座席管理服务 (Go)

### 前端应用
- **agent-workbench**: 座席工作台 (React)
- **admin-portal**: 管理后台 (React + Ant Design Pro)
- **mobile-app**: 移动端应用 (React Native)
- **data-dashboard**: 数据大屏 (React + ECharts)

## 快速开始

### 环境要求

- Docker 20.10+
- Kubernetes 1.21+
- Node.js 16+
- Go 1.19+
- Python 3.9+
- Java 11+

### 本地开发

1. 克隆项目
```bash
git clone https://github.com/cloudcallcenter/cloudcallcenter.git
cd cloudcallcenter
```

2. 安装依赖
```bash
# 安装前端依赖
cd frontend/agent-workbench && npm install
cd ../admin-portal && npm install

# 安装后端依赖
cd ../../backend/call-center-service && go mod download
cd ../survey-service && npm install
# ... 其他服务
```

3. 启动基础服务
```bash
# 使用 docker-compose 启动依赖服务
docker-compose -f docker/docker-compose.yml up -d
```

4. 启动微服务
```bash
# 启动各个微服务
./scripts/start-services.sh
```

5. 访问应用
- 座席工作台: http://localhost:3000
- 管理后台: http://localhost:3001
- API文档: http://localhost:8080/swagger

### Docker部署

```bash
# 构建镜像
docker build -t cloudcallcenter/call-center-service:latest ./backend/call-center-service

# 运行容器
docker run -p 8080:8080 cloudcallcenter/call-center-service:latest
```

### Kubernetes部署

```bash
# 部署到Kubernetes
kubectl apply -f k8s/

# 查看部署状态
kubectl get pods -n cloudcallcenter
```

## 开发指南

### 项目结构

```
cloudcallcenter/
├── backend/                 # 后端服务
│   ├── call-center-service/ # 呼叫中心服务
│   ├── survey-service/      # 问卷服务
│   ├── quota-service/       # 配额服务
│   └── ...
├── frontend/               # 前端应用
│   ├── agent-workbench/    # 座席工作台
│   ├── admin-portal/       # 管理后台
│   └── ...
├── k8s/                    # Kubernetes配置
├── docker/                 # Docker配置
├── scripts/                # 构建脚本
└── docs/                   # 文档
```

### 编码规范

- 后端遵循各语言的官方编码规范
- 前端遵循 Airbnb JavaScript Style Guide
- 提交信息遵循 Conventional Commits

### API规范

所有API遵循RESTful设计原则，详见 [API文档](./docs/api.md)

## 性能指标

- **并发座席**: 10,000+
- **并发呼叫**: 100,000+
- **API响应**: < 100ms
- **系统可用性**: 99.99%
- **数据处理**: 1000万+条/天

## 技术栈

### 后端
- Go + Gin/gRPC
- Node.js + Express + TypeScript
- Python + FastAPI
- Java + Spring Boot

### 前端
- React 18 + TypeScript
- Redux Toolkit
- Ant Design Pro 5
- ECharts + D3.js

### 基础设施
- Kubernetes + Docker
- PostgreSQL + Redis
- Kafka + Elasticsearch
- Prometheus + Grafana

## 贡献指南

欢迎贡献代码！请查看 [CONTRIBUTING.md](./CONTRIBUTING.md) 了解详情。

## 许可证

本项目采用 [MIT License](./LICENSE) 开源协议。

## 联系我们

- 官网: https://cloudcallcenter.io
- 邮箱: support@cloudcallcenter.io
- 论坛: https://forum.cloudcallcenter.io
>>>>>>> Incoming (Background Agent changes)
