# Harbor 镜像仓库插件

## 功能介绍

Harbor 插件为 k8m 提供了完整的 Harbor 镜像仓库管理功能，支持：

- 多 Harbor 仓库配置管理
- 项目列表查看
- 镜像仓库浏览
- 镜像制品（Artifact）管理
- 镜像标签查看
- 镜像删除操作

## 安装与启用

### 1. 插件已自动注册

Harbor 插件已在系统启动时自动注册，无需手动注册。

### 2. 启用插件

在 k8m 管理界面中：

1. 进入 **系统管理** → **插件管理**
2. 找到 **Harbor镜像仓库** 插件
3. 点击 **安装** 按钮
4. 点击 **启用** 按钮
5. 点击 **启动** 按钮

插件启动后，左侧菜单会出现 **Harbor仓库** 菜单项。

## 使用指南

### 1. 配置 Harbor 仓库

首先需要添加 Harbor 仓库配置：

1. 点击左侧菜单 **Harbor仓库** → **仓库管理**
2. 点击 **新增仓库** 按钮
3. 填写以下信息：
   - **仓库名称**：自定义名称，用于标识
   - **Harbor地址**：Harbor 服务器地址，如 `https://harbor.example.com`
   - **用户名**：Harbor 登录用户名
   - **密码**：Harbor 登录密码
   - **跳过TLS验证**：如果使用自签名证书，需要开启
   - **设为默认仓库**：是否设为默认仓库
   - **描述**：可选的描述信息

4. 点击 **提交** 保存

### 2. 测试连接

添加仓库后，可以测试连接是否正常：

1. 在仓库列表中找到对应的仓库
2. 点击 **测试连接** 按钮
3. 如果连接成功，会显示成功提示

### 3. 查看项目列表

1. 点击左侧菜单 **Harbor仓库** → **项目列表**
2. 在下拉框中选择要查看的 Harbor 仓库
3. 系统会自动加载该仓库下的所有项目
4. 可以查看每个项目的仓库数量、Chart 数量等信息

### 4. 浏览镜像仓库

1. 点击左侧菜单 **Harbor仓库** → **镜像仓库**
2. 选择 Harbor 仓库
3. 输入项目名称
4. 系统会显示该项目下的所有镜像仓库
5. 点击 **查看镜像** 可以查看仓库中的所有镜像制品

### 5. 管理镜像制品

在镜像列表中，可以：

- 查看镜像的 Digest（摘要）
- 查看镜像的所有标签（Tags）
- 查看镜像大小
- 查看推送时间
- 删除不需要的镜像（需要确认）

## API 接口

### 仓库管理 API

```
GET    /mgm/plugins/harbor/registries          # 获取仓库列表
POST   /mgm/plugins/harbor/registries          # 创建仓库
PUT    /mgm/plugins/harbor/registries/{id}     # 更新仓库
DELETE /mgm/plugins/harbor/registries/{id}     # 删除仓库
POST   /mgm/plugins/harbor/registries/{id}/test # 测试连接
```

### Harbor 数据 API

```
GET    /k8s/cluster/{cluster}/plugins/harbor/projects      # 获取项目列表
GET    /k8s/cluster/{cluster}/plugins/harbor/repositories  # 获取仓库列表
GET    /k8s/cluster/{cluster}/plugins/harbor/artifacts     # 获取镜像制品列表
DELETE /k8s/cluster/{cluster}/plugins/harbor/artifacts     # 删除镜像制品
```

## 数据库表结构

插件会创建以下数据库表：

### harbor_registries

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键 |
| name | string | 仓库名称（唯一） |
| url | string | Harbor 地址 |
| username | string | 用户名 |
| password | string | 密码 |
| insecure | bool | 是否跳过 TLS 验证 |
| description | string | 描述 |
| is_default | bool | 是否为默认仓库 |
| created_at | time | 创建时间 |
| updated_at | time | 更新时间 |
| deleted_at | time | 删除时间（软删除） |

## 权限说明

- **仓库管理**：需要登录用户权限
- **项目和镜像查看**：需要登录用户权限
- **镜像删除**：需要登录用户权限，且 Harbor 用户需要有相应权限

## 注意事项

1. **密码安全**：密码存储在数据库中，建议使用专门的 Harbor 用户，并限制其权限
2. **TLS 证书**：生产环境建议使用有效的 TLS 证书，避免开启"跳过 TLS 验证"
3. **删除操作**：删除镜像是不可逆操作，请谨慎操作
4. **API 限制**：Harbor API 有速率限制，频繁调用可能会被限流

## 故障排查

### 连接失败

1. 检查 Harbor 地址是否正确
2. 检查用户名和密码是否正确
3. 检查网络连接是否正常
4. 如果使用自签名证书，确保开启了"跳过 TLS 验证"

### 无法查看项目

1. 确保 Harbor 用户有查看项目的权限
2. 检查项目名称是否正确
3. 查看 k8m 后端日志获取详细错误信息

### 无法删除镜像

1. 确保 Harbor 用户有删除镜像的权限
2. 检查镜像是否被标记为不可变（Immutable）
3. 查看 Harbor 服务器日志

## 开发说明

### 目录结构

```
harbor/
├── README.md                    # 本文档
├── metadata.go                  # 插件元数据
├── lifecycle.go                 # 生命周期管理
├── models/
│   └── models.go               # 数据模型
├── service/
│   └── harbor_client.go        # Harbor API 客户端
├── controller/
│   ├── registry.go             # 仓库管理控制器
│   └── project.go              # 项目和镜像控制器
├── route/
│   ├── cluster_api.go          # 集群相关路由
│   ├── mgm_api.go              # 管理相关路由
│   └── admin_api.go            # 管理员路由
└── frontend/
    ├── registries.json         # 仓库管理页面
    ├── projects.json           # 项目列表页面
    └── repositories.json       # 镜像仓库页面
```

### 扩展功能

可以继续添加以下功能：

1. 镜像扫描结果查看
2. 镜像漏洞报告
3. 镜像签名验证
4. Webhook 配置
5. 复制规则管理
6. 垃圾回收配置

## 版本历史

- **v1.0.0** (2025-02-25)
  - 初始版本
  - 支持基本的 Harbor 仓库管理
  - 支持项目、仓库、镜像查看
  - 支持镜像删除操作
