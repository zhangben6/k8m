package controller

import (
	"fmt"
	"strconv"

	"github.com/weibaohui/k8m/internal/dao"
	"github.com/weibaohui/k8m/pkg/comm/utils/amis"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/service"
	"github.com/weibaohui/k8m/pkg/response"
	"k8s.io/klog/v2"
)

// ListProjects 获取Harbor项目列表
func ListProjects(c *response.Context) {
	registryIDStr := c.Query("registry_id")
	if registryIDStr == "" {
		amis.WriteJsonError(c, fmt.Errorf("registry_id参数不能为空"))
		return
	}

	registryID, err := strconv.ParseUint(registryIDStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var registry models.HarborRegistry
	if err := dao.DB().First(&registry, registryID).Error; err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := service.NewHarborClient(&registry)
	projects, err := client.ListProjects(page, pageSize)
	if err != nil {
		klog.Errorf("获取Harbor项目列表失败: %v", err)
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonListWithTotal(c, int64(len(projects)), projects)
}

// ListRepositories 获取Harbor仓库列表
func ListRepositories(c *response.Context) {
	registryIDStr := c.Query("registry_id")
	projectName := c.Query("project_name")

	if registryIDStr == "" || projectName == "" {
		amis.WriteJsonError(c, fmt.Errorf("registry_id和project_name参数不能为空"))
		return
	}

	registryID, err := strconv.ParseUint(registryIDStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var registry models.HarborRegistry
	if err := dao.DB().First(&registry, registryID).Error; err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := service.NewHarborClient(&registry)
	repos, err := client.ListRepositories(projectName, page, pageSize)
	if err != nil {
		klog.Errorf("获取Harbor仓库列表失败: %v", err)
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonListWithTotal(c, int64(len(repos)), repos)
}

// ListArtifacts 获取Harbor镜像制品列表
func ListArtifacts(c *response.Context) {
	registryIDStr := c.Query("registry_id")
	projectName := c.Query("project_name")
	repoName := c.Query("repo_name")

	if registryIDStr == "" || projectName == "" || repoName == "" {
		amis.WriteJsonError(c, fmt.Errorf("registry_id、project_name和repo_name参数不能为空"))
		return
	}

	registryID, err := strconv.ParseUint(registryIDStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var registry models.HarborRegistry
	if err := dao.DB().First(&registry, registryID).Error; err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))

	client := service.NewHarborClient(&registry)
	artifacts, err := client.ListArtifacts(projectName, repoName, page, pageSize)
	if err != nil {
		klog.Errorf("获取Harbor镜像制品列表失败: %v", err)
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonListWithTotal(c, int64(len(artifacts)), artifacts)
}

// DeleteArtifact 删除Harbor镜像制品
func DeleteArtifact(c *response.Context) {
	registryIDStr := c.Query("registry_id")
	projectName := c.Query("project_name")
	repoName := c.Query("repo_name")
	digest := c.Query("digest")

	if registryIDStr == "" || projectName == "" || repoName == "" || digest == "" {
		amis.WriteJsonError(c, fmt.Errorf("参数不完整"))
		return
	}

	registryID, err := strconv.ParseUint(registryIDStr, 10, 64)
	if err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	var registry models.HarborRegistry
	if err := dao.DB().First(&registry, registryID).Error; err != nil {
		amis.WriteJsonError(c, err)
		return
	}

	client := service.NewHarborClient(&registry)
	if err := client.DeleteArtifact(projectName, repoName, digest); err != nil {
		klog.Errorf("删除Harbor镜像制品失败: %v", err)
		amis.WriteJsonError(c, err)
		return
	}

	amis.WriteJsonOK(c)
}
