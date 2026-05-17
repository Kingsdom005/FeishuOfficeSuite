# 飞书企业办公套件 (Feishu Office Suite)

基于 Go + Kratos 框架打造的企业级飞书办公套件后端服务，支持用户管理、消息推送、日历、审批等核心功能。

## 技术栈

| 组件 | 技术 | 说明 |
|------|------|------|
| 框架 | Kratos | Go 微服务框架 |
| 通信 | gRPC + Protobuf | 高性能 RPC 通信 |
| 数据库 | MySQL + Ent | ORM 数据访问层 |
| 缓存 | Redis | 分布式缓存 |
| 消息队列 | Kafka + Asynq | 异步任务处理 |
| 权限 | Casbin | 基于角色的访问控制 |
| 可观测性 | OpenTelemetry | 链路追踪 |
| 监控 | Prometheus + Grafana | 指标监控 |
| 容器化 | Docker + Kubernetes | 部署环境 |

## 项目结构

```
FeishuOfficeSuite/
├── api/                          # Protobuf 定义
│   └── feishu/
│       └── v1/                   # API v1 版本
│           ├── user.proto        # 用户服务
│           ├── message.proto     # 消息服务
│           └── calendar.proto    # 日历/审批服务
├── cmd/                          # 应用入口
│   ├── server/                   # 主服务
│   │   └── main.go
│   └── worker/                   # 异步任务 worker
├── configs/                      # 配置文件
│   └── config.yaml
├── internal/                     # 内部包
│   ├── data/                     # 数据层
│   ├── domain/                   # 领域模型
│   ├── handler/                  # gRPC Handler
│   ├── middleware/               # 中间件
│   ├── job/                      # 异步任务
│   ├── kafka/                    # Kafka 生产者/消费者
│   ├── casbin/                   # 权限管理
│   ├── cache/                    # Redis 缓存
│   └── otracing/                # 链路追踪
├── pkg/                          # 公共包
│   └── feishu/                   # 飞书 SDK 封装
├── deployments/                  # 部署配置
│   ├── docker/                   # Docker Compose
│   └── k8s/                      # Kubernetes
├── docs/                         # 文档
├── Makefile
└── go.mod
```

## 快速开始

### 前置要求

- Go 1.21+
- Docker & Docker Compose
- Make

### 1. 初始化项目

```bash
# 克隆项目
cd FeishuOfficeSuite

# 安装依赖
make init

# 生成 Protobuf 代码
make proto
```

### 2. 配置飞书应用

参考 [配置说明文档](docs/CONFIGURATION.md) 申请飞书应用并配置凭证。

### 3. 本地开发

```bash
# 启动基础设施 (MySQL, Redis, Kafka, etc.)
docker-compose -f deployments/docker/docker-compose.yaml up -d

# 运行服务
make run
```

服务启动后:
- HTTP API: http://localhost:8000
- gRPC: http://localhost:9000
- Prometheus: http://localhost:9091
- Grafana: http://localhost:3000 (admin/admin123)
- Jaeger: http://localhost:16686

### 4. 构建部署

```bash
# 构建 Docker 镜像
make docker

# 推送到镜像仓库
make docker-push
```

## Docker Compose 部署

```bash
# 启动所有服务
docker-compose -f deployments/docker/docker-compose.yaml up -d

# 查看服务状态
docker-compose -f deployments/docker/docker-compose.yaml ps

# 查看日志
docker-compose -f deployments/docker/docker-compose.yaml logs -f app

# 停止服务
docker-compose -f deployments/docker/docker-compose.yaml down
```

### 环境变量

创建 `.env` 文件:

```env
FEISHU_APP_ID=cli_xxxxxxxxxxxxxx
FEISHU_APP_SECRET=your_app_secret
```

## Kubernetes 部署

### 前提条件

- Kubernetes 1.21+
- kubectl 配置完成

### 部署步骤

```bash
# 创建命名空间
kubectl apply -f deployments/k8s/namespace.yaml

# 部署配置
kubectl apply -f deployments/k8s/configmap.yaml
kubectl apply -f deployments/k8s/secret.yaml

# 部署有状态服务 (MySQL, Redis)
kubectl apply -f deployments/k8s/statefulset.yaml

# 部署应用
kubectl apply -f deployments/k8s/deployment.yaml

# 部署服务
kubectl apply -f deployments/k8s/service.yaml
```

## API 文档

### 用户服务 (FeishuUser)

```protobuf
service FeishuUser {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc GetUserByEmail(GetUserByEmailRequest) returns (GetUserByEmailResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (google.protobuf.Empty);
}
```

### 消息服务 (FeishuMessage)

```protobuf
service FeishuMessage {
  rpc SendMessage(SendMessageRequest) returns (SendMessageResponse);
  rpc SendTextMessage(SendTextMessageRequest) returns (SendMessageResponse);
  rpc GetMessage(GetMessageRequest) returns (GetMessageResponse);
  rpc ListMessages(ListMessagesRequest) returns (ListMessagesResponse);
  rpc ReplyMessage(ReplyMessageRequest) returns (SendMessageResponse);
}
```

### 日历服务 (FeishuCalendar)

```protobuf
service FeishuCalendar {
  rpc GetCalendar(GetCalendarRequest) returns (GetCalendarResponse);
  rpc ListCalendars(ListCalendarsRequest) returns (ListCalendarsResponse);
  rpc CreateEvent(CreateEventRequest) returns (CreateEventResponse);
  rpc GetEvent(GetEventRequest) returns (GetEventResponse);
  rpc ListEvents(ListEventsRequest) returns (ListEventsResponse);
  rpc UpdateEvent(UpdateEventRequest) returns (google.protobuf.Empty);
  rpc DeleteEvent(DeleteEventRequest) returns (google.protobuf.Empty);
}
```

### 审批服务 (FeishuApproval)

```protobuf
service FeishuApproval {
  rpc GetApproval(GetApprovalRequest) returns (GetApprovalResponse);
  rpc ListApprovals(ListApprovalsRequest) returns (ListApprovalsResponse);
  rpc CreateApproval(CreateApprovalRequest) returns (CreateApprovalResponse);
  rpc GetApprovalInstance(GetApprovalInstanceRequest) returns (GetApprovalInstanceResponse);
  rpc ListApprovalInstances(ListApprovalInstancesRequest) returns (ListApprovalInstancesResponse);
  rpc ApproveInstance(ApproveInstanceRequest) returns (google.protobuf.Empty);
  rpc RejectInstance(RejectInstanceRequest) returns (google.protobuf.Empty);
}
```

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make init` | 初始化项目，安装依赖 |
| `make proto` | 生成 Protobuf 代码 |
| `make ent` | 生成 Ent 代码 |
| `make build` | 构建应用 |
| `make run` | 运行服务 (需先启动 Docker) |
| `make stop` | 停止 Docker Compose |
| `make test` | 运行测试 |
| `make lint` | 运行代码检查 |
| `make fmt` | 代码格式化 |
| `make docker` | 构建 Docker 镜像 |
| `make k8s-deploy` | 部署到 Kubernetes |
| `make k8s-delete` | 从 Kubernetes 删除 |

## 监控与告警

### Prometheus 指标

- `feishu_http_requests_total` - HTTP 请求总数
- `feishu_http_request_duration_seconds` - HTTP 请求延迟
- `feishu_grpc_requests_total` - gRPC 请求总数
- `feishu_grpc_request_duration_seconds` - gRPC 请求延迟

### Grafana Dashboard

导入 `deployments/docker/grafana-dashboard.json` 查看预置面板。

## 开发指南

### 添加新的 API

1. 在 `api/feishu/v1/` 下创建/编辑 `.proto` 文件
2. 运行 `make proto` 生成代码
3. 在 `internal/handler/` 实现 Handler
4. 在 `main.go` 注册服务

### 添加新的数据模型

1. 在 `internal/data/ent/schema/` 创建 Ent schema
2. 运行 `make ent` 生成代码
3. 在 `internal/domain/entity/` 定义实体

### 添加新的异步任务

1. 在 `internal/job/executor.go` 定义任务类型
2. 实现 `ExecuteTask` 方法
3. 使用 `job.Enqueue()` 提交任务

## License

MIT