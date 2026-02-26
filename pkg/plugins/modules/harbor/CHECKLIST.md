# Harbor 插件部署检查清单

## 开发阶段 ✅

- [x] 创建插件目录结构
- [x] 实现插件元数据（metadata.go）
- [x] 实现生命周期管理（lifecycle.go）
- [x] 创建数据模型（models/models.go）
- [x] 实现 Harbor API 客户端（service/harbor_client.go）
- [x] 实现控制器（controller/）
- [x] 配置路由（route/）
- [x] 创建前端页面（frontend/）
- [x] 编写单元测试
- [x] 编写文档（README、QUICKSTART、EXAMPLE）
- [x] 在 list.go 中添加插件常量
- [x] 在 registrar.go 中注册插件
- [x] 执行 go mod tidy
- [x] 编译测试通过

## 构建阶段 ⏳

- [ ] 构建前端资源
  ```bash
  cd ui
  pnpm install
  pnpm run build
  ```

- [ ] 构建后端
  ```bash
  go build -o k8m main.go
  ```

- [ ] 验证前端资源复制
  ```bash
  ls -la dist/pages/plugins/harbor/
  # 应该看到 registries.json, projects.json, repositories.json
  ```

## 测试阶段 ⏳

### 单元测试
- [ ] 运行单元测试
  ```bash
  go test ./pkg/plugins/modules/harbor/...
  ```

### 集成测试
- [ ] 启动 k8m
  ```bash
  ./k8m -v 6
  ```

- [ ] 登录系统（k8m/k8m）

- [ ] 安装插件
  - [ ] 进入"系统管理" → "插件管理"
  - [ ] 找到"Harbor镜像仓库"插件
  - [ ] 点击"安装"
  - [ ] 点击"启用"
  - [ ] 点击"启动"
  - [ ] 验证状态为"运行中"

- [ ] 检查菜单
  - [ ] 左侧菜单出现"Harbor仓库"
  - [ ] 子菜单包含：仓库管理、项目列表、镜像仓库

- [ ] 测试仓库管理
  - [ ] 添加 Harbor 仓库配置
  - [ ] 测试连接
  - [ ] 编辑仓库配置
  - [ ] 验证密码隐藏

- [ ] 测试项目列表
  - [ ] 选择仓库
  - [ ] 查看项目列表
  - [ ] 验证数据显示正确

- [ ] 测试镜像仓库
  - [ ] 选择仓库和项目
  - [ ] 查看镜像列表
  - [ ] 查看镜像详情
  - [ ] 测试删除功能（谨慎！）

### API 测试
- [ ] 测试仓库管理 API
  ```bash
  # 获取列表
  curl -X GET http://localhost:3618/mgm/plugins/harbor/registries \
    -H "Authorization: Bearer YOUR_TOKEN"
  
  # 创建仓库
  curl -X POST http://localhost:3618/mgm/plugins/harbor/registries \
    -H "Authorization: Bearer YOUR_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"name":"test","url":"https://harbor.example.com","username":"admin","password":"password"}'
  ```

- [ ] 测试 Harbor 数据 API
  ```bash
  # 获取项目
  curl -X GET "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/projects?registry_id=1" \
    -H "Authorization: Bearer YOUR_TOKEN"
  ```

## 数据库检查 ⏳

- [ ] 验证表创建
  ```sql
  SHOW TABLES LIKE 'harbor_registries';
  DESC harbor_registries;
  ```

- [ ] 验证数据插入
  ```sql
  SELECT * FROM harbor_registries;
  ```

## 日志检查 ⏳

- [ ] 检查插件注册日志
  ```
  注册harbor插件成功
  ```

- [ ] 检查插件安装日志
  ```
  安装Harbor插件成功
  ```

- [ ] 检查插件启动日志
  ```
  启动Harbor插件后台任务
  ```

- [ ] 检查 API 调用日志
  ```
  获取Harbor仓库配置列表
  获取Harbor项目列表失败/成功
  ```

## 性能测试 ⏳

- [ ] 测试大量数据加载
  - [ ] 100+ 项目
  - [ ] 1000+ 镜像

- [ ] 测试并发请求
  - [ ] 多用户同时访问
  - [ ] 多个仓库同时查询

- [ ] 测试超时处理
  - [ ] 网络延迟
  - [ ] Harbor 服务慢响应

## 安全测试 ⏳

- [ ] 测试权限控制
  - [ ] 未登录用户访问
  - [ ] 普通用户权限
  - [ ] 管理员权限

- [ ] 测试密码安全
  - [ ] 密码显示隐藏
  - [ ] 密码更新逻辑

- [ ] 测试 TLS 验证
  - [ ] 有效证书
  - [ ] 自签名证书
  - [ ] 跳过验证选项

## 错误处理测试 ⏳

- [ ] 测试网络错误
  - [ ] Harbor 服务不可达
  - [ ] 超时处理

- [ ] 测试认证错误
  - [ ] 错误的用户名/密码
  - [ ] 过期的凭证

- [ ] 测试数据错误
  - [ ] 无效的 URL
  - [ ] 不存在的项目
  - [ ] 不存在的镜像

## 文档检查 ⏳

- [ ] README.md 完整性
- [ ] QUICKSTART.md 可用性
- [ ] EXAMPLE.md 示例正确性
- [ ] API 文档准确性
- [ ] 代码注释完整性

## 部署准备 ⏳

- [ ] 准备发布说明
- [ ] 更新 CHANGELOG
- [ ] 准备演示视频/截图
- [ ] 准备用户培训材料

## 生产环境检查 ⏳

- [ ] 备份数据库
- [ ] 准备回滚方案
- [ ] 监控告警配置
- [ ] 性能基线记录

## 用户反馈 ⏳

- [ ] 收集用户反馈
- [ ] 记录问题和建议
- [ ] 规划下一版本功能

---

## 检查说明

- ✅ 已完成
- ⏳ 待完成
- ❌ 失败/需要修复
- ⚠️ 需要注意

## 快速测试命令

```bash
# 1. 编译
go build -o k8m main.go

# 2. 启动（调试模式）
./k8m -v 6

# 3. 查看日志
tail -f k8m.log

# 4. 测试 API
curl -X GET http://localhost:3618/mgm/plugins/harbor/registries \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 问题记录

| 日期 | 问题 | 状态 | 解决方案 |
|------|------|------|----------|
| 2025-02-25 | 编译错误：字符串不能作为 error | ✅ 已解决 | 使用 fmt.Errorf() |
| 2025-02-25 | 类型断言错误 | ✅ 已解决 | 直接使用 items[i] |
| | | | |

## 联系方式

如有问题，请联系：
- 开发者：[你的名字]
- Email：[你的邮箱]
- GitHub Issue：[项目地址]
