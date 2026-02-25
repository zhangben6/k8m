package route

import (
	"github.com/go-chi/chi/v5"
	"github.com/weibaohui/k8m/pkg/plugins/modules"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/controller"
	"github.com/weibaohui/k8m/pkg/response"
)

// RegisterManagementRoutes 注册Harbor插件的管理相关路由
func RegisterManagementRoutes(mrg chi.Router) {
	prefix := "/plugins/" + modules.PluginNameHarbor

	// Harbor仓库配置管理API
	mrg.Get(prefix+"/registries", response.Adapter(controller.ListRegistries))
	mrg.Post(prefix+"/registries", response.Adapter(controller.CreateRegistry))
	mrg.Put(prefix+"/registries/{id}", response.Adapter(controller.UpdateRegistry))
	mrg.Delete(prefix+"/registries/{id}", response.Adapter(controller.DeleteRegistry))
	mrg.Post(prefix+"/registries/{id}/test", response.Adapter(controller.TestRegistryConnection))
}
