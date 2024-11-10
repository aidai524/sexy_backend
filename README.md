# Sexy Backend

基于 Gin 框架的 RESTful API 服务，使用 Supabase 作为数据存储。

## 目录结构

sexy_backend/
├── cmd/                # 主要的应用程序入口
│   └── main.go        # 主程序入口文件
├── internal/          # 私有应用程序和库代码
│   ├── handler/      # HTTP 处理器
│   │   └── user.go   # 用户相关的 HTTP 处理器
│   ├── middleware/   # HTTP 中间件
│   │   └── cors.go   # CORS 中间件
│   ├── model/       # 数据模型
│   │   └── user.go  # 用户模型
│   └── service/     # 业务逻辑层
│       └── user_service.go # 用户服务
├── pkg/              # 可以被外部应用程序使用的库代码
│   └── supabase/    # Supabase 客户端封装
│       └── client.go # Supabase 客户端
├── config/          # 配置文件
│   ├── config.go    # 配置加载逻辑
│   └── config.yaml  # 配置文件
└── README.md        # 项目文档

## 功能特性

- RESTful API 设计
- Supabase 数据库集成
- 用户管理 CRUD 操作
- CORS 支持
- 配置管理
- 错误处理

## 快速开始

### 前置要求

- Go 1.16 或更高版本
- Supabase 账号和项目
- Git

### 安装步骤

1. 克隆项目
```bash
git clone https://github.com/yourusername/sexy_backend.git
cd sexy_backend
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境
```bash
cp config/config.yaml.example config/config.yaml
```

4. 修改配置文件 `config/config.yaml`，填入您的 Supabase 配置：
```yaml
server:
  port: "8080"
supabase:
  url: "your-project-url"     # 从 Supabase 项目设置 > API 获取
  api_key: "your-api-key"     # 从 Supabase 项目设置 > API 获取 service_role key
```

5. 运行服务
```bash
go run cmd/main.go
```

## API 文档

### 用户管理 API

#### 创建用户

**请求**:
```http
POST /api/v1/users
Content-Type: application/json

{
    "email": "user@example.com",
    "username": "johndoe"
}
```

**示例**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe"
  }'
```

**成功响应** (201 Created):
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-03-06T12:00:00Z",
    "updated_at": "2024-03-06T12:00:00Z"
}
```

#### 获取用户

**请求**:
```http
GET /api/v1/users/:id
```

**示例**:
```bash
curl http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**成功响应** (200 OK):
```json
{
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "username": "johndoe",
    "created_at": "2024-03-06T12:00:00Z",
    "updated_at": "2024-03-06T12:00:00Z"
}
```

#### 更新用户

**请求**:
```http
PUT /api/v1/users/:id
Content-Type: application/json

{
    "email": "updated@example.com",
    "username": "janedoe"
}
```

**示例**:
```bash
curl -X PUT http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "updated@example.com",
    "username": "janedoe"
  }'
```

**成功响应** (200 OK):
```json
{
    "message": "User updated",
    "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### 删除用户

**请求**:
```http
DELETE /api/v1/users/:id
```

**示例**:
```bash
curl -X DELETE http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000
```

**成功响应** (200 OK):
```json
{
    "message": "User deleted",
    "id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### 错误响应

所有 API 在发生错误时都会返回统一的错误格式：

```json
{
    "error": "错误信息描述"
}
```

常见的错误状态码：
- 400 Bad Request: 请求参数错误
- 404 Not Found: 资源不存在
- 500 Internal Server Error: 服务器内部错误

## 开发指南

### 添加新的 API 端点

1. 在 `internal/model` 中定义数据模型
2. 在 `internal/service` 中实现业务逻辑
3. 在 `internal/handler` 中创建处理器
4. 在 `cmd/main.go` 中注册路由

### 配置管理

项目使用 Viper 进行配置管理，支持：
- 配置文件 (config.yaml)
- 环境变量
- 默认值

环境变量映射：
- `SUPABASE_URL`: Supabase 项目 URL
- `SUPABASE_KEY`: Supabase service_role key

### 错误处理

- 在 service 层返回具体的错误信息
- 在 handler 层统一处理错误响应
- 使用适当的 HTTP 状态码

## 部署

### 使用 Docker 部署

1. 构建镜像：
```bash
docker build -t sexy_backend .
```

2. 运行容器：
```bash
docker run -p 8080:8080 \
  -e SUPABASE_URL=your-project-url \
  -e SUPABASE_KEY=your-service-role-key \
  sexy_backend
```

### 环境变量配置

在生产环境中，建议使用环境变量进行配置：

```bash
export SUPABASE_URL=your-project-url
export SUPABASE_KEY=your-service-role-key
```

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'Add some amazing feature'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 创建 Pull Request

## 许可证

[MIT License](LICENSE)

## 联系方式

如有问题或建议，请提交 Issue 或 Pull Request。
