# 飞书办公套件配置文件说明

## 配置文件位置

- **本地开发**: `configs/config.yaml`
- **Docker 环境**: `deployments/docker/docker-compose.yaml`
- **Kubernetes 环境**: `deployments/k8s/configmap.yaml`

## 统一配置填写位置

所有敏感配置和密钥统一填写在以下位置:

### Docker Compose 环境变量

在 `deployments/docker/docker-compose.yaml` 的 `app` 服务中配置:

```yaml
app:
  environment:
    FEISHU_APP_ID: ${FEISHU_APP_ID:-your_app_id}
    FEISHU_APP_SECRET: ${FEISHU_APP_SECRET:-your_app_secret}
```

或创建 `.env` 文件:

```env
FEISHU_APP_ID=your_app_id
FEISHU_APP_SECRET=your_app_secret
```

### Kubernetes Secret

在 `deployments/k8s/secret.yaml` 中配置:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: feishu-secrets
  namespace: feishu-suite
type: Opaque
stringData:
  FEISHU_APP_ID: "your_app_id"
  FEISHU_APP_SECRET: "your_app_secret"
```

## 飞书密钥申请和配置

### 1. 飞书应用创建

1. 访问 [飞书开放平台](https://open.feishu.cn/app)
2. 登录开发者账号
3. 点击「创建企业自建应用」
4. 填写应用名称和描述
5. 创建后获取 `App ID` 和 `App Secret`

### 2. 配置应用权限

在飞书开放平台的应用详情页:

1. 进入「权限管理」
2. 添加以下权限:
   - `contact:user.employee_id:readonly` - 获取用户基本信息
   - `contact:user.email:readonly` - 获取用户邮箱
   - `contact:user.phone:readonly` - 获取用户手机号
   - `im:message` - 发送消息
   - `im:chat` - 管理群聊
   - `calendar:calendar` - 日历管理
   - `approval:approval` - 审批管理

### 3. 配置事件订阅

1. 进入「事件订阅」
2. 添加以下事件:
   - `im.message.receive_v1` - 接收消息
   - `contact.user.modify_v3` - 用户变更
   - `contact.department.modify_v3` - 部门变更

### 4. 配置应用凭证

将获取的 `App ID` 和 `App Secret` 填入:

```yaml
feishu:
  app_id: "cli_xxxxxxxxxxxxxx"
  app_secret: "your_app_secret"
```

## 完整配置项说明

### server - 服务配置

```yaml
server:
  http:
    network: tcp          # 网络类型
    timeout: 60           # 超时时间(秒)
  grpc:
    network: tcp          # gRPC 网络类型
    timeout: 60           # gRPC 超时时间(秒)
```

### database - MySQL 数据库配置

```yaml
database:
  dsn: "user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
  max_open_conns: 100     # 最大打开连接数
  max_idle_conns: 10     # 最大空闲连接数
  max_lifetime: 3600     # 连接最大生命周期(秒)
```

**Docker Compose 环境变量**:
```yaml
MYSQL_DSN: "feishu:feishu123@tcp(mysql:3306)/feishu_suite?charset=utf8mb4&parseTime=True&loc=Local"
```

### redis - Redis 缓存配置

```yaml
redis:
  addr: "localhost:6379"     # Redis 地址
  password: ""                # Redis 密码
  db: 0                       # 数据库编号
  pool_size: 100              # 连接池大小
```

**Docker Compose 环境变量**:
```yaml
REDIS_ADDR: "redis:6379"
REDIS_PASSWORD: ""
```

### kafka - Kafka 消息队列配置

```yaml
kafka:
  brokers:
    - "localhost:9092"         # Kafka Broker 地址列表
  topic: "feishu-events"       # 默认主题
```

**Docker Compose 环境变量**:
```yaml
KAFKA_BROKERS: "kafka:29092"
```

### asynq - Asynq 异步任务配置

```yaml
asynq:
  redis:
    addr: "localhost:6379"     # Redis 地址(用于 Asynq)
    password: ""               # Redis 密码
    db: 1                      # 使用不同的数据库编号
```

### feishu - 飞书配置

```yaml
feishu:
  app_id: "cli_xxxxxxxxxxxxxx"    # 飞书应用 App ID
  app_secret: "your_app_secret"    # 飞书应用 App Secret
  token: ""                         # 自定义 Token(可选)
```

### open_telemetry - 链路追踪配置

```yaml
open_telemetry:
  endpoint: "localhost:4317"       # OTLP gRPC 接收端点
```

**Docker Compose 环境变量**:
```yaml
OTEL_ENDPOINT: "otel-collector:4317"
```

### metrics - Prometheus 监控配置

```yaml
metrics:
  enabled: true                    # 是否启用监控
  port: 9090                        # 监控指标端口
```

### registry - 服务注册与发现

```yaml
registry:
  consul:
    address: "localhost:8500"       # Consul 地址
    timeout: 5                      # 超时时间(秒)
```

**Docker Compose 环境变量**:
```yaml
CONSUL_ADDRESS: "consul:8500"
```

## 环境变量完整列表

| 变量名 | 描述 | 默认值 | 必填 |
|--------|------|--------|------|
| `FEISHU_APP_ID` | 飞书应用 App ID | - | 是 |
| `FEISHU_APP_SECRET` | 飞书应用 App Secret | - | 是 |
| `MYSQL_DSN` | MySQL 连接字符串 | - | 是 |
| `REDIS_ADDR` | Redis 连接地址 | `redis:6379` | 否 |
| `REDIS_PASSWORD` | Redis 密码 | 空 | 否 |
| `KAFKA_BROKERS` | Kafka Broker 地址 | `kafka:29092` | 否 |
| `OTEL_ENDPOINT` | OpenTelemetry 收集器地址 | `otel-collector:4317` | 否 |
| `CONSUL_ADDRESS` | Consul 服务地址 | `consul:8500` | 否 |

## 配置优先级

1. 命令行参数 (最高优先级)
2. 环境变量
3. 配置文件 `config.yaml`
4. 默认值 (最低优先级)