# SIEM 日志审计系统

基于 OpenSearch 3.4.0 的企业级日志审计系统，提供高并发日志采集、存储、检索和可视化功能。

## 系统架构

- **后端**: Go 语言高性能服务，处理日志接收和查询
- **前端**: 轻量级 Web UI，基于 Nginx + 原生 JavaScript
- **存储**: OpenSearch 3.4.0 分布式搜索引擎
- **部署**: Docker 容器化 + GitHub Actions CI/CD

## 快速开始

### 前提条件

- Docker 20.10+
- Docker Compose 2.0+

### 本地部署

1. 克隆项目并配置环境变量：

```bash
git clone <repository-url>
cd SIEM
cp .env.example .env
```

2. 修改 `.env` 文件中的镜像配置：

```bash
IMAGE_TAG=latest
DOCKER_REGISTRY=ghcr.io
GITHUB_REPOSITORY=your-org/siem
```

3. 执行部署脚本：

```bash
chmod +x scripts/deploy.sh
./scripts/deploy.sh pull
```

4. 访问系统：

- 前端界面: http://localhost
- 后端 API: http://localhost:8080
- OpenSearch: http://localhost:9200

## API 接口

### 日志接收

```bash
POST /api/logs
Content-Type: application/json

{
  "timestamp": "2025-12-23T10:00:00Z",
  "level": "info",
  "source": "app-server-01",
  "message": "User login successful",
  "metadata": {
    "user_id": "12345",
    "ip": "192.168.1.100"
  }
}
```

### 日志搜索

```bash
POST /api/logs/search
Content-Type: application/json

{
  "query": "level:error",
  "start_time": "2025-12-23T00:00:00Z",
  "end_time": "2025-12-23T23:59:59Z",
  "size": 100,
  "from": 0
}
```

### 健康检查

```bash
GET /api/health
```

## CI/CD 流程

### GitHub Actions 自动构建

项目使用 GitHub Actions 自动构建和发布 Docker 镜像：

1. **触发条件**：
   - 推送到 main/develop 分支
   - 创建版本标签 (v*)
   - Pull Request

2. **构建流程**：
   - 编译后端 Go 服务
   - 构建前端静态资源
   - 推送镜像到 GitHub Container Registry
   - 发布版本时自动创建部署包

3. **镜像标签规则**：
   - `main` 分支: `main-backend`, `main-frontend`
   - 标签版本: `v1.0.0-backend`, `v1.0.0-frontend`
   - Git SHA: `sha-abc123-backend`

### 生产部署

1. 下载发布包：

```bash
wget https://github.com/your-org/siem/releases/download/v1.0.0/siem-deployment-v1.0.0.tar.gz
tar -xzf siem-deployment-v1.0.0.tar.gz
cd deployment
```

2. 配置环境变量：

```bash
export IMAGE_TAG=v1.0.0
export DOCKER_REGISTRY=ghcr.io
export GITHUB_REPOSITORY=your-org/siem
```

3. 执行部署：

```bash
./deploy.sh pull
```

## 项目结构

```
SIEM/
├── backend/              # Go 后端服务
│   ├── main.go          # 主程序
│   ├── go.mod           # Go 依赖
│   └── Dockerfile       # 后端镜像
├── frontend/            # Web 前端
│   ├── index.html       # 主页面
│   ├── nginx.conf       # Nginx 配置
│   └── Dockerfile       # 前端镜像
├── scripts/             # 部署脚本
│   ├── deploy.sh        # 部署脚本
│   └── stop.sh          # 停止脚本
├── .github/
│   └── workflows/
│       └── build.yml    # CI/CD 配置
├── docker-compose.yml   # 容器编排
└── .env.example         # 环境变量模板
```

## 功能特性

### 已实现

- ✅ 高并发日志接收 (Go Goroutine)
- ✅ 实时日志搜索和过滤
- ✅ 时间范围查询
- ✅ 日志级别统计
- ✅ 分页展示
- ✅ 响应式 Web UI
- ✅ Docker 容器化部署
- ✅ GitHub Actions CI/CD

### 规划中

- 📋 日志聚合分析
- 📋 告警规则引擎
- 📋 用户权限管理
- 📋 数据可视化图表
- 📋 日志导出功能
- 📋 多租户支持

## 运维管理

### 查看日志

```bash
docker-compose logs -f
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f opensearch
```

### 停止服务

```bash
./scripts/stop.sh
```

### 数据备份

```bash
docker exec siem-opensearch /usr/share/opensearch/bin/opensearch-snapshot-restore
```

## 性能优化

- Go 后端使用 Gorilla Mux 路由，支持高并发
- OpenSearch 单节点部署，内存配置 512MB (可根据需求调整)
- Nginx 静态资源缓存
- Docker 镜像多阶段构建，减小镜像体积

## 安全建议

- 生产环境启用 OpenSearch 安全插件
- 配置 HTTPS/TLS 加密
- 使用强密码和访问控制
- 定期更新依赖和镜像
- 限制网络访问端口

## 许可证

本项目仅供学习和商业内部使用。

## 技术支持

如有问题，请提交 Issue 或联系开发团队。
