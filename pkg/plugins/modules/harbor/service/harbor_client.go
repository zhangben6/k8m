package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
	"k8s.io/klog/v2"
)

// HarborClient Harbor API 客户端
type HarborClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewHarborClient 创建Harbor客户端
func NewHarborClient(registry *models.HarborRegistry) *HarborClient {
	client := &HarborClient{
		BaseURL:  strings.TrimSuffix(registry.URL, "/"),
		Username: registry.Username,
		Password: registry.Password,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	if registry.Insecure {
		client.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return client
}

// doRequest 执行HTTP请求
func (c *HarborClient) doRequest(method, path string, body io.Reader) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v2.0%s", c.BaseURL, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Harbor API错误: %d - %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// Project Harbor项目
type Project struct {
	ProjectID    int64  `json:"project_id"`
	Name         string `json:"name"`
	Public       bool   `json:"metadata.public"`
	RepoCount    int    `json:"repo_count"`
	ChartCount   int    `json:"chart_count"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
}

// Repository Harbor仓库
type Repository struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	ProjectID    int64  `json:"project_id"`
	Description  string `json:"description"`
	PullCount    int64  `json:"pull_count"`
	ArtifactCount int   `json:"artifact_count"`
	CreationTime string `json:"creation_time"`
	UpdateTime   string `json:"update_time"`
}

// Artifact 镜像制品
type Artifact struct {
	ID           int64    `json:"id"`
	Type         string   `json:"type"`
	Digest       string   `json:"digest"`
	Tags         []Tag    `json:"tags"`
	PushTime     string   `json:"push_time"`
	PullTime     string   `json:"pull_time"`
	Size         int64    `json:"size"`
}

// Tag 镜像标签
type Tag struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PushTime     string `json:"push_time"`
	PullTime     string `json:"pull_time"`
	Immutable    bool   `json:"immutable"`
}

// ListProjects 获取项目列表
func (c *HarborClient) ListProjects(page, pageSize int) ([]Project, error) {
	path := fmt.Sprintf("/projects?page=%d&page_size=%d", page, pageSize)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		klog.Errorf("获取Harbor项目列表失败: %v", err)
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// ListRepositories 获取仓库列表
func (c *HarborClient) ListRepositories(projectName string, page, pageSize int) ([]Repository, error) {
	path := fmt.Sprintf("/projects/%s/repositories?page=%d&page_size=%d", projectName, page, pageSize)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		klog.Errorf("获取Harbor仓库列表失败: %v", err)
		return nil, err
	}

	var repos []Repository
	if err := json.Unmarshal(data, &repos); err != nil {
		return nil, err
	}

	return repos, nil
}

// ListArtifacts 获取镜像制品列表
func (c *HarborClient) ListArtifacts(projectName, repoName string, page, pageSize int) ([]Artifact, error) {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts?page=%d&page_size=%d&with_tag=true", 
		projectName, repoName, page, pageSize)
	data, err := c.doRequest("GET", path, nil)
	if err != nil {
		klog.Errorf("获取Harbor镜像制品列表失败: %v", err)
		return nil, err
	}

	var artifacts []Artifact
	if err := json.Unmarshal(data, &artifacts); err != nil {
		return nil, err
	}

	return artifacts, nil
}

// DeleteArtifact 删除镜像制品
func (c *HarborClient) DeleteArtifact(projectName, repoName, digest string) error {
	path := fmt.Sprintf("/projects/%s/repositories/%s/artifacts/%s", projectName, repoName, digest)
	_, err := c.doRequest("DELETE", path, nil)
	if err != nil {
		klog.Errorf("删除Harbor镜像制品失败: %v", err)
		return err
	}
	return nil
}

// TestConnection 测试连接
func (c *HarborClient) TestConnection() error {
	_, err := c.doRequest("GET", "/systeminfo", nil)
	return err
}
