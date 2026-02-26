# Harbor 镜像仓库插件开发总结

## 项目概述

为 k8m 项目成功开发了一个完整的 Harbor 镜像仓库管理插件，实现了 Harbor 仓库的配置管理、项目浏览、镜像查看和删除等核心功能。

## 插件结构

```
pkg/plugins/modules/harbor/
├── README.md                    # 完整文档
├── QUICKSTART.md                # 快速开始指南
├── EXAMPLE.md                   # 使用示例
├── metadata.go                  # 插件元数据定义
├── lifecycle.go                 # 生命周期管理
├── models/
│   └── models.go               # 数据模型（HarborRegistry）
├── service/
│   ├── harbor_client.go        # Harbor API 客户端
│   └── harbor_client_test.go  # 单元测试
├── controller/
│   ├── registry.go             # 仓库配置管理控制器
│   └── project.go              # 项目和镜像管理控制器
├── route/
│   ├── cluster_api.go          # 集群相关路由
│   ├── mgm_api.go              # 管理相关路由
│   └── admin_api.go            # 管理员路由
└── frontend/
    ├── registries.json         # 仓库管理页面（AMIS）
    ├── projects.json           # 项目列表页面（AMIS）
    └── repositories.json       # 镜像仓库页面（AMIS）
```

## 核心功能

### 1. 仓库配置管理
- ✅ 添加/编辑/删除 Harbor 仓库配置
- ✅ 支持多个 Harbor 实例
- ✅ 连接测试功能
- ✅ 密码安全处理（显示时隐藏）
- ✅ 支持自签名证书（跳过 TLS 验证）
- ✅ 默认仓库设置

### 2. Harbor 数据浏览
- ✅ 项目列表查看
- ✅ 镜像仓库列表
- ✅ 镜像制品（Artifact）详情
- ✅ 镜像标签（Tags）展示
- ✅ 镜像大小、推送时间等信息

### 3. 镜像管理
- ✅ 镜像删除功能
- ✅ 删除确认机制
- ✅ 权限控制

### 4. 前端界面
- ✅ 基于 AMIS 的低代码页面
- ✅ 响应式设计
- ✅ 友好的用户交互
- ✅ 完整的表单验证

## 技术实现

### 后端技术
- **语言**: Go 1.24.9
- **框架**: Chi Router
- **ORM**: GORM
- **HTTP 客户端**: 标准库 net/http
- **日志**: klog/v2

### 前端技术
- **UI 框架**: 百度 AMIS
- **配置方式**: JSON Schema
- **图标**: Font Awesome

### Harbor API
- **版本**: Harbor API v2.0
- **认证**: Basic Auth
- **支持的操作**:
  - 获取项目列表
  - 获取仓库列表
  - 获取镜像制品
  - 删除镜像制品
  - 系统信息查询

## 插件生命周期

```
Uninstalled → Install → Installed → Enable → Enabled → Start → Running
                                                ↓
                                            Disable → Disabled
                                                ↓
                                            Uninstall → Uninstalled
```

### 生命周期方法实现
- ✅ `Install()` - 创建数据库表
- ✅ `Upgrade()` - 数据库迁移
- ✅ `Enable()` - 启用插件
- ✅ `Disable()` - 禁用插件
- ✅ `Uninstall()` - 删除数据（可选保留）
- ✅ `Start()` - 启动后台任务
- ✅ `Stop()` - 停止后台任务

## API 接口

### 管理接口（/mgm）
```
GET    /mgm/plugins/harbor/registries          # 获取仓库列表
POST   /mgm/plugins/harbor/registries          # 创建仓库
PUT    /mgm/plugins/harbor/registries/{id}     # 更新仓库
DELETE /mgm/plugins/harbor/registries/{id}     # 删除仓库
POST   /mgm/plugins/harbor/registries/{id}/test # 测试连接
```

### 集群接口（/k8s/cluster/{cluster}）
```
GET    /plugins/harbor/projects                # 获取项目列表
GET    /plugins/harbor/repositories            # 获取仓库列表
GET    /plugins/harbor/artifacts               # 获取镜像制品
DELETE /plugins/harbor/artifacts               # 删除镜像制品
```

## 数据库设计

### harbor_registries 表
| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 仓库名称（唯一索引） |
| url | string | Harbor 地址 |
| username | string | 用户名 |
| password | string | 密码 |
| insecure | bool | 跳过 TLS 验证 |
| description | string | 描述 |
| is_default | bool | 是否默认仓库 |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |
| deleted_at | timestamp | 删除时间（软删除） |

## 菜单结构

```
Harbor仓库 (fa-brands fa-docker)
├── 仓库管理 (fa-solid fa-server)
├── 项目列表 (fa-solid fa-folder)
└── 镜像仓库 (fa-solid fa-box)
```

## 测试

### 单元测试
- ✅ Harbor 客户端创建测试
- ✅ URL 处理测试
- ✅ 配置验证测试

### 编译测试
```bash
go build ./pkg/plugins/modules/harbor/...
# 编译成功 ✅
```

## 文档

### 用户文档
1. **README.md** - 完整的功能介绍和使用指南
2. **QUICKSTART.md** - 5 分钟快速上手指南
3. **EXAMPLE.md** - 详细的使用示例和场景

### 开发文档
- 插件架构说明
- API 接口文档
- 数据库设计
- 扩展指南

## 集成步骤

### 1. 插件注册
已在以下文件中完成注册：
- ✅ `pkg/plugins/modules/list.go` - 添加插件常量
- ✅ `pkg/plugins/modules/registrar/registrar.go` - 注册插件

### 2. 依赖管理
```bash
go mod tidy  # 已执行 ✅
```

### 3. 前端资源
前端 JSON 配置会在构建时自动复制到 `dist/pages/plugins/harbor/`

## 使用流程

### 管理员操作
1. 启用插件：系统管理 → 插件管理 → Harbor镜像仓库 → 安装/启用/启动
2. 配置仓库：Harbor仓库 → 仓库管理 → 新增仓库
3. 测试连接：确保配置正确

### 用户操作
1. 查看项目：Harbor仓库 → 项目列表
2. 浏览镜像：Harbor仓库 → 镜像仓库
3. 管理镜像：查看详情、删除旧镜像

## 安全考虑

1. **密码存储**: 存储在数据库中，显示时自动隐藏
2. **TLS 验证**: 支持跳过验证（用于自签名证书）
3. **权限控制**: 基于 k8m 的用户权限系统
4. **操作确认**: 删除操作需要用户确认

## 性能优化

1. **连接复用**: HTTP 客户端支持连接池
2. **超时控制**: 30 秒请求超时
3. **分页查询**: 支持分页加载大量数据
4. **后台任务**: 5 分钟定时任务（可扩展）

## 扩展建议

### 短期扩展
1. 镜像扫描结果查看
2. 镜像漏洞报告
3. 镜像复制功能
4. 批量删除镜像

### 长期扩展
1. Harbor Webhook 集成
2. 镜像签名验证
3. 复制规则管理
4. 垃圾回收配置
5. 用户和权限管理
6. 审计日志查看

## 已知限制

1. 仅支持 Harbor API v2.0+
2. 密码未加密存储（建议使用专用账户）
3. 不支持 Harbor v1.x API
4. 删除操作不可恢复

## 故障排查

### 常见问题
1. **连接失败**: 检查 URL、用户名、密码、网络
2. **TLS 错误**: 开启"跳过 TLS 验证"
3. **权限不足**: 确保 Harbor 用户有相应权限
4. **无法删除**: 检查镜像是否不可变

### 调试方法
```bash
# 启用详细日志
./k8m -v 6

# 查看插件状态
# 在管理界面查看插件运行状态
```

## 总结

成功为 k8m 开发了一个功能完整、架构清晰的 Harbor 镜像仓库管理插件。插件遵循 k8m 的插件规范，实现了完整的生命周期管理，提供了友好的用户界面和完善的文档。

### 亮点
- ✅ 完整的插件生命周期实现
- ✅ 清晰的代码结构和注释
- ✅ 基于 AMIS 的低代码前端
- ✅ 完善的文档和示例
- ✅ 单元测试覆盖
- ✅ 编译通过，无语法错误

### 下一步
1. 前端构建测试
2. 集成测试
3. 实际 Harbor 环境测试
4. 用户反馈收集
5. 功能迭代优化

---

**开发时间**: 2025-02-25  
**插件版本**: v1.0.0  
**状态**: ✅ 开发完成，编译通过
