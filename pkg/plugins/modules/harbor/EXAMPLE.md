# Harbor 插件使用示例

## 场景 1：添加 Harbor 仓库并查看镜像

### 步骤 1：添加 Harbor 仓库

通过 API 添加 Harbor 仓库：

```bash
curl -X POST http://localhost:3618/mgm/plugins/harbor/registries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "公司Harbor",
    "url": "https://harbor.company.com",
    "username": "admin",
    "password": "Harbor12345",
    "insecure": false,
    "is_default": true,
    "description": "公司内部Harbor镜像仓库"
  }'
```

### 步骤 2：测试连接

```bash
curl -X POST http://localhost:3618/mgm/plugins/harbor/registries/1/test \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 步骤 3：查看项目列表

```bash
curl -X GET "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/projects?registry_id=1&page=1&perPage=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

响应示例：

```json
{
  "status": 0,
  "msg": "ok",
  "data": {
    "items": [
      {
        "project_id": 1,
        "name": "library",
        "repo_count": 10,
        "chart_count": 5,
        "creation_time": "2024-01-01T00:00:00Z"
      }
    ],
    "total": 1
  }
}
```

### 步骤 4：查看镜像仓库

```bash
curl -X GET "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/repositories?registry_id=1&project_name=library&page=1&perPage=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 步骤 5：查看镜像制品

```bash
curl -X GET "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/artifacts?registry_id=1&project_name=library&repo_name=nginx&page=1&perPage=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 场景 2：清理旧镜像

### 查找并删除旧镜像

```bash
# 1. 获取镜像列表
ARTIFACTS=$(curl -s -X GET "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/artifacts?registry_id=1&project_name=library&repo_name=nginx" \
  -H "Authorization: Bearer YOUR_TOKEN")

# 2. 解析并删除特定镜像
curl -X DELETE "http://localhost:3618/k8s/cluster/your-cluster/plugins/harbor/artifacts?registry_id=1&project_name=library&repo_name=nginx&digest=sha256:abc123..." \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 场景 3：在 Kubernetes 中使用 Harbor 镜像

### 创建 Secret

```bash
kubectl create secret docker-registry harbor-secret \
  --docker-server=harbor.company.com \
  --docker-username=admin \
  --docker-password=Harbor12345 \
  --docker-email=admin@company.com
```

### 在 Pod 中使用

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: test-pod
spec:
  containers:
  - name: nginx
    image: harbor.company.com/library/nginx:latest
  imagePullSecrets:
  - name: harbor-secret
```

## 场景 4：批量管理多个 Harbor 仓库

### 添加多个仓库

```bash
# 开发环境 Harbor
curl -X POST http://localhost:3618/mgm/plugins/harbor/registries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "开发环境",
    "url": "https://harbor-dev.company.com",
    "username": "dev-user",
    "password": "DevPassword123"
  }'

# 生产环境 Harbor
curl -X POST http://localhost:3618/mgm/plugins/harbor/registries \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "生产环境",
    "url": "https://harbor-prod.company.com",
    "username": "prod-user",
    "password": "ProdPassword123",
    "is_default": true
  }'
```

### 查看所有仓库

```bash
curl -X GET "http://localhost:3618/mgm/plugins/harbor/registries?page=1&perPage=20" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## 场景 5：使用 Go 代码调用 Harbor 客户端

```go
package main

import (
    "fmt"
    "github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
    "github.com/weibaohui/k8m/pkg/plugins/modules/harbor/service"
)

func main() {
    // 创建 Harbor 仓库配置
    registry := &models.HarborRegistry{
        Name:     "my-harbor",
        URL:      "https://harbor.example.com",
        Username: "admin",
        Password: "Harbor12345",
        Insecure: false,
    }

    // 创建客户端
    client := service.NewHarborClient(registry)

    // 测试连接
    if err := client.TestConnection(); err != nil {
        fmt.Printf("连接失败: %v\n", err)
        return
    }
    fmt.Println("连接成功！")

    // 获取项目列表
    projects, err := client.ListProjects(1, 20)
    if err != nil {
        fmt.Printf("获取项目列表失败: %v\n", err)
        return
    }

    fmt.Printf("找到 %d 个项目:\n", len(projects))
    for _, project := range projects {
        fmt.Printf("- %s (仓库数: %d)\n", project.Name, project.RepoCount)
    }

    // 获取仓库列表
    repos, err := client.ListRepositories("library", 1, 20)
    if err != nil {
        fmt.Printf("获取仓库列表失败: %v\n", err)
        return
    }

    fmt.Printf("\n找到 %d 个仓库:\n", len(repos))
    for _, repo := range repos {
        fmt.Printf("- %s (镜像数: %d, 拉取次数: %d)\n", 
            repo.Name, repo.ArtifactCount, repo.PullCount)
    }
}
```

## 场景 6：前端页面使用

### 在浏览器中访问

1. 登录 k8m 系统
2. 点击左侧菜单 **Harbor仓库**
3. 选择 **仓库管理** 添加 Harbor 配置
4. 选择 **项目列表** 查看项目
5. 选择 **镜像仓库** 浏览和管理镜像

### 页面功能

- **仓库管理页面**：
  - 添加/编辑/删除 Harbor 仓库配置
  - 测试连接
  - 设置默认仓库

- **项目列表页面**：
  - 选择 Harbor 仓库
  - 查看项目列表
  - 查看项目统计信息

- **镜像仓库页面**：
  - 选择 Harbor 仓库和项目
  - 查看镜像仓库列表
  - 查看镜像制品详情
  - 删除不需要的镜像

## 常见问题

### Q: 如何处理自签名证书？

A: 在添加仓库时，开启"跳过TLS验证"选项。

### Q: 密码会明文存储吗？

A: 密码存储在数据库中，建议使用专门的 Harbor 用户，并限制其权限。未来版本会考虑加密存储。

### Q: 可以同时管理多个 Harbor 吗？

A: 可以，插件支持添加多个 Harbor 仓库配置。

### Q: 删除镜像后可以恢复吗？

A: 不可以，删除操作是不可逆的，请谨慎操作。

### Q: 支持 Harbor v1.x 吗？

A: 插件使用 Harbor API v2.0，建议使用 Harbor v2.0 及以上版本。
