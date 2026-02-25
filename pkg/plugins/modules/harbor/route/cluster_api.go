package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/weibaohui/k8m/pkg/plugins/modules"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/controller"
	"github.com/weibaohui/k8m/pkg/response"
)

// RegisterClusterRoutes 注册Harbor插件的集群相关路由
func RegisterClusterRoutes(crg chi.Router) {
	prefix := "/plugins/" + modules.PluginNameHarbor

	// Harbor项目和镜像相关API
	crg.Get(prefix+"/projects", response.Adapter(controller.ListProjects))
	crg.Get(prefix+"/repositories", response.Adapter(controller.ListRepositories))
	crg.Get(prefix+"/artifacts", response.Adapter(controller.ListArtifacts))
	crg.Delete(prefix+"/artifacts", response.Adapter(controller.DeleteArtifact))
}
