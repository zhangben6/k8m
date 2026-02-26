package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/weibaohui/k8m/pkg/plugins/modules"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/controller"
	"github.com/weibaohui/k8m/pkg/response"
)

// RegisterClusterRoutes 注册Harbor插件的集群相关路由
// Harbor 功能与 K8s 集群无关，这里保留空实现
func RegisterClusterRoutes(crg chi.Router) {
	// Harbor 镜像仓库管理不依赖 K8s 集群
	// 所有功能都在 ManagementRouter 中实现
	_ = modules.PluginNameHarbor
	_ = controller.ListProjects
	_ = response.Adapter
}
