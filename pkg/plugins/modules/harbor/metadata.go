package harbor

import (
	"github.com/weibaohui/k8m/pkg/plugins"
	"github.com/weibaohui/k8m/pkg/plugins/modules"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/route"
)

var Metadata = plugins.Module{
	Meta: plugins.Meta{
		Name:        modules.PluginNameHarbor,
		Title:       "Harbor镜像仓库",
		Version:     "1.0.0",
		Description: "Harbor镜像仓库管理插件，支持项目、镜像、标签管理",
	},
	Tables: []string{
		"harbor_registries",
	},
	Menus: []plugins.Menu{
		{
			Key:   "plugin_harbor_index",
			Title: "Harbor仓库",
			Icon:  "fa-brands fa-docker",
			Order: 50,
			Children: []plugins.Menu{
				{
					Key:         "plugin_harbor_registries",
					Title:       "仓库管理",
					Icon:        "fa-solid fa-server",
					EventType:   "custom",
					CustomEvent: `() => loadJsonPage("/plugins/harbor/registries")`,
					Order:       1,
				},
				{
					Key:         "plugin_harbor_projects",
					Title:       "项目列表",
					Icon:        "fa-solid fa-folder",
					EventType:   "custom",
					CustomEvent: `() => loadJsonPage("/plugins/harbor/projects")`,
					Order:       2,
				},
				{
					Key:         "plugin_harbor_repositories",
					Title:       "镜像仓库",
					Icon:        "fa-solid fa-box",
					EventType:   "custom",
					CustomEvent: `() => loadJsonPage("/plugins/harbor/repositories")`,
					Order:       3,
				},
			},
		},
	},
	Dependencies: []string{},
	RunAfter:     []string{},

	Lifecycle:         &HarborLifecycle{},
	ClusterRouter:     route.RegisterClusterRoutes,
	ManagementRouter:  route.RegisterManagementRoutes,
	PluginAdminRouter: route.RegisterPluginAdminRoutes,
}
