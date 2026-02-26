package controller

import (
	"strconv"

	"github.com/weibaohui/k8m/internal/dao"
	"github.com/weibaohui/k8m/pkg/comm/utils/amis"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/service"
	"github.com/weibaohui/k8m/pkg/response"
	"k8s.io/klog/v2"
)

// ListRegistries 获取Harbor仓库配置列表
func ListRegistries(c *response.Context) {
	klog.V(6).Infof("获取Harbor仓库配置列表")

	params := dao.BuildParams(c)
	m := &models.HarborRegistry{}
	items, total, err := dao.GenericQuery(params, m)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	// 隐藏密码
	for i := range items {
		items[i].Password = "******"
	}

	amis.WriteJsonListWithTotal(c, total, items)
}

// CreateRegistry 创建Harbor仓库配置
func CreateRegistry(c *response.Context) {
	var req models.HarborRegistry
	if err := c.ShouldBindJSON(&req); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	klog.V(6).Infof("创建Harbor仓库配置，名称=%s", req.Name)

	// 测试连接
	client := service.NewHarborClient(&req)
	if err := client.TestConnection(); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	params := dao.BuildParams(c)
	if err := dao.GenericSave(params, &req); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonData(c, req)
}

// UpdateRegistry 更新Harbor仓库配置
func UpdateRegistry(c *response.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var req models.HarborRegistry
	if err = c.ShouldBindJSON(&req); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	klog.V(6).Infof("更新Harbor仓库配置，ID=%d", id64)
	req.ID = uint(id64)

	// 如果密码是******，则不更新密码
	if req.Password == "******" {
		var old models.HarborRegistry
		if err := dao.DB().First(&old, id64).Error; err == nil {
			req.Password = old.Password
		}
	}

	// 测试连接
	client := service.NewHarborClient(&req)
	if err := client.TestConnection(); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	params := dao.BuildParams(c)
	if err = dao.GenericSave(params, &req); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonOK(c)
}

// DeleteRegistry 删除Harbor仓库配置
func DeleteRegistry(c *response.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	klog.V(6).Infof("删除Harbor仓库配置，ID=%d", id64)

	params := dao.BuildParams(c)
	if err = dao.GenericDelete(params, &models.HarborRegistry{}, []int64{int64(id64)}); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonOK(c)
}

// TestRegistryConnection 测试Harbor仓库连接
func TestRegistryConnection(c *response.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		amis.WriteJsonError(c, fmt.Errorf("仓库ID不能为空"))
		return
	}
	
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var registry models.HarborRegistry
	if err := dao.DB().First(&registry, id64).Error; err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	client := service.NewHarborClient(&registry)
	if err := client.TestConnection(); err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonOK(c)
}
