# Harbor 插件使用指南

## 一、Harbor 插件的设计逻辑

### 1.1 Harbor 与 K8s 集群的关系

**重要说明**：Harbor 插件是一个**独立的镜像仓库管理工具**，它与 K8s 集群没有强绑定关系。

```
┌─────────────────┐
│   K8M 系统      │
│  ┌───────────┐  │
│  │ Harbor插件 │  │ ──连接──> Harbor 服务器 1 (https://harbor1.example.com)
│  └───────────┘  │ ──连接──> Harbor 服务器 2 (https://harbor2.example.com)
│                 │ ──连接──> Harbor 服务器 3 (https://harbor3.example.com)
└─────────────────┘
```

### 1.2 为什么 API 路径包含 cluster 参数？

当前实现中，API 路径使用了 `/k8s/cluster/{cluster}/plugins/harbor/...` 格式，这主要是为了：
- 保持与其他插件的路由一致性
- 预留未来可能的集群级别配置能力
- **但实际上 Harbor 配置是全局的，不区分集群**

## 二、如何使用 Harbor 插件

### 2.1 前置条件

你需要有一个**真实运行的 Harbor 服务**，可以是：

1. **公有云 Harbor 服务**
   - 阿里云容器镜像服务 ACR
   - 腾讯云容器镜像服务 TCR
   - 华为云容器镜像服务 SWR

2. **自建 Harbor 服务**
   ```bash
   # 使用 Docker Compose 快速部署 Harbor
   wget https://github.com/goharbor/harbor/releases/download/v2.10.0/harbor-offline-installer-v2.10.0.tgz
   tar xzvf harbor-offline-installer-v2.10.0.tgz
   cd harbor
   ./install.sh
   ```

3. **Harbor 官方演示环境**（仅用于测试）
   - URL: https://demo.goharbor.io
   - 用户名: admin
   - 密码: Harbor12345

### 2.2 配置步骤

#### 步骤 1：启用插件

1. 登录 K8M 系统
2. 进入 **系统管理** → **插件管理**
3. 找到 **Harbor镜像仓库** 插件
4. 依次点击：**安装** → **启用** → **启动**

#### 步骤 2：添加 Harbor 仓库配置

1. 点击左侧菜单 **Harbor仓库** → **仓库管理**
2. 点击 **新增仓库** 按钮
3. 填写配置信息：

| 字段 | 说明 | 示例 |
|------|------|------|
| 仓库名称 | 自定义名称，用于标识 | 生产环境Harbor |
| Harbor地址 | Harbor 服务器完整 URL | https://harbor.example.com |
| 用户名 | Harbor 登录用户名 | admin |
| 密码 | Harbor 登录密码 | Harbor12345 |
| 跳过TLS验证 | 自签名证书时需要开启 | ☑️ |
| 设为默认仓库 | 是否设为默认 | ☑️ |
| 描述 | 可选的描述信息 | 公司生产环境镜像仓库 |

4. 点击 **提交** 保存
5. 点击 **测试连接** 验证配置是否正确

#### 步骤 3：浏览和管理镜像

**查看项目列表**
1. 点击 **Harbor仓库** → **项目列表**
2. 选择要查看的 Harbor 仓库
3. 查看所有项目及统计信息

**查看镜像仓库**
1. 点击 **Harbor仓库** → **镜像仓库**
2. 选择 Harbor 仓库
3. 输入项目名称（如：library）
4. 查看该项目下的所有镜像

**查看镜像详情**
1. 在镜像列表中点击 **查看镜像**
2. 查看镜像的所有标签（Tags）
3. 查看镜像大小、推送时间、Digest 等信息

**删除镜像**
1. 在镜像详情中找到要删除的镜像
2. 点击 **删除** 按钮
3. 确认删除操作

## 三、常见使用场景

### 3.1 管理多个 Harbor 实例

如果你的公司有多个 Harbor 服务器（开发、测试、生产），可以全部添加到插件中：

```
仓库 1: 开发环境Harbor (https://harbor-dev.example.com)
仓库 2: 测试环境Harbor (https://harbor-test.example.com)
仓库 3: 生产环境Harbor (https://harbor-prod.example.com)
```

在使用时，通过下拉框选择要操作的 Harbor 实例。

### 3.2 清理旧镜像

定期清理不再使用的镜像标签，释放存储空间：

1. 进入 **镜像仓库**
2. 选择项目和镜像
3. 查看所有标签
4. 删除旧的、不再使用的标签

### 3.3 查看镜像信息

在部署应用前，查看镜像的详细信息：
- 镜像大小（评估下载时间）
- 推送时间（确认镜像是否最新）
- 标签列表（选择合适的版本）

## 四、API 接口说明

### 4.1 仓库管理 API（全局级别）

这些 API 用于管理 Harbor 仓库配置，不依赖于 K8s 集群：

```
GET    /mgm/plugins/harbor/registries          # 获取所有仓库配置
POST   /mgm/plugins/harbor/registries          # 添加新仓库
PUT    /mgm/plugins/harbor/registries/{id}     # 更新仓库配置
DELETE /mgm/plugins/harbor/registries/{id}     # 删除仓库配置
POST   /mgm/plugins/harbor/registries/{id}/test # 测试仓库连接
```

### 4.2 Harbor 数据 API（集群级别路由）

这些 API 用于查询和操作 Harbor 中的数据：

```
GET    /k8s/cluster/{cluster}/plugins/harbor/projects      # 获取项目列表
GET    /k8s/cluster/{cluster}/plugins/harbor/repositories  # 获取镜像仓库列表
GET    /k8s/cluster/{cluster}/plugins/harbor/artifacts     # 获取镜像制品列表
DELETE /k8s/cluster/{cluster}/plugins/harbor/artifacts     # 删除镜像制品
```

**注意**：虽然路径包含 `{cluster}` 参数，但当前实现中 Harbor 配置是全局的。

## 五、故障排查

### 5.1 连接失败

**问题**：点击"测试连接"失败

**排查步骤**：
1. ✅ 检查 Harbor 地址是否正确（包括 http/https 协议）
2. ✅ 检查用户名和密码是否正确
3. ✅ 检查网络连接（K8M 服务器能否访问 Harbor）
4. ✅ 如果使用自签名证书，确保开启"跳过TLS验证"
5. ✅ 查看 K8M 后端日志：`./k8m -v 6`

### 5.2 看不到项目或镜像

**问题**：Harbor 仓库中明明有项目，但插件中看不到

**排查步骤**：
1. ✅ 确保 Harbor 用户有查看项目的权限
2. ✅ 检查是否选择了正确的 Harbor 仓库
3. ✅ 检查项目名称是否正确（区分大小写）
4. ✅ 查看浏览器控制台的网络请求，确认 API 调用是否成功

### 5.3 无法删除镜像

**问题**：点击删除按钮后失败

**排查步骤**：
1. ✅ 确保 Harbor 用户有删除镜像的权限
2. ✅ 检查镜像是否被标记为不可变（Immutable）
3. ✅ 检查镜像是否正在被使用
4. ✅ 查看 Harbor 服务器日志

### 5.4 API 很多但不知道怎么用

**解答**：
- 前端界面已经封装了所有 API 调用
- 你不需要直接调用这些 API
- 通过界面操作即可完成所有功能
- API 主要用于：
  - 前端页面的数据加载
  - 未来可能的自动化脚本
  - 与其他系统的集成

## 六、与 K8s 集群的集成建议

虽然 Harbor 插件本身不依赖 K8s 集群，但你可以这样使用：

### 6.1 为不同集群配置不同的 Harbor

```
开发集群 (dev-cluster) → 使用 Harbor 开发环境
测试集群 (test-cluster) → 使用 Harbor 测试环境
生产集群 (prod-cluster) → 使用 Harbor 生产环境
```

### 6.2 配置 K8s 集群使用 Harbor

在 K8s 集群中创建 ImagePullSecret：

```bash
kubectl create secret docker-registry harbor-secret \
  --docker-server=https://harbor.example.com \
  --docker-username=admin \
  --docker-password=Harbor12345 \
  --docker-email=admin@example.com
```

在 Pod 中使用：

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-app
spec:
  containers:
  - name: app
    image: harbor.example.com/library/nginx:latest
  imagePullSecrets:
  - name: harbor-secret
```

## 七、总结

1. **Harbor 插件是独立的镜像仓库管理工具**，不强依赖 K8s 集群
2. **你需要有真实的 Harbor 服务**才能使用这个插件
3. **通过界面配置 Harbor 连接信息**，然后就可以浏览和管理镜像
4. **API 路径中的 cluster 参数**主要是为了路由一致性，实际配置是全局的
5. **所有功能都可以通过界面完成**，不需要直接调用 API

## 八、下一步建议

如果你想改进这个插件，可以考虑：

1. **真正实现集群级别的 Harbor 配置**
   - 不同集群可以配置不同的 Harbor 仓库
   - 在数据库中增加 cluster_id 字段

2. **自动发现集群中使用的镜像**
   - 扫描集群中所有 Pod 使用的镜像
   - 自动关联到 Harbor 中的镜像

3. **镜像安全扫描集成**
   - 显示镜像的漏洞扫描结果
   - 阻止使用有严重漏洞的镜像

4. **自动清理策略**
   - 配置自动清理规则
   - 定期删除旧的、未使用的镜像

需要我帮你实现这些改进吗？
