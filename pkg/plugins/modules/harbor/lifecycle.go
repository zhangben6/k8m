package harbor

import (
	"context"
	"time"

	"github.com/weibaohui/k8m/pkg/plugins"
	"github.com/weibaohui/k8m/pkg/plugins/modules/harbor/models"
	"k8s.io/klog/v2"
)

// HarborLifecycle Harbor插件生命周期实现
type HarborLifecycle struct {
	cancelStart context.CancelFunc
}

// Install 安装Harbor插件，初始化数据库表
func (h *HarborLifecycle) Install(ctx plugins.InstallContext) error {
	if err := models.InitDB(); err != nil {
		klog.V(6).Infof("安装Harbor插件失败: %v", err)
		return err
	}
	klog.V(6).Infof("安装Harbor插件成功")
	return nil
}

// Upgrade 升级Harbor插件，执行数据库迁移
func (h *HarborLifecycle) Upgrade(ctx plugins.UpgradeContext) error {
	klog.V(6).Infof("升级Harbor插件：从版本 %s 到版本 %s", ctx.FromVersion(), ctx.ToVersion())
	if err := models.UpgradeDB(ctx.FromVersion(), ctx.ToVersion()); err != nil {
		return err
	}
	return nil
}

// Enable 启用Harbor插件
func (h *HarborLifecycle) Enable(ctx plugins.EnableContext) error {
	klog.V(6).Infof("启用Harbor插件")
	return nil
}

// Disable 禁用Harbor插件
func (h *HarborLifecycle) Disable(ctx plugins.BaseContext) error {
	klog.V(6).Infof("禁用Harbor插件")
	return nil
}

// Uninstall 卸载Harbor插件，根据keepData参数决定是否删除相关的表及数据
func (h *HarborLifecycle) Uninstall(ctx plugins.UninstallContext) error {
	klog.V(6).Infof("卸载Harbor插件")
	if !ctx.KeepData() {
		if err := models.DropDB(); err != nil {
			return err
		}
		klog.V(6).Infof("卸载Harbor插件完成，已删除相关表及数据")
	} else {
		klog.V(6).Infof("卸载Harbor插件完成，保留相关表及数据")
	}
	return nil
}

// Start 启动Harbor插件的后台任务（不可阻塞）
func (h *HarborLifecycle) Start(ctx plugins.BaseContext) error {
	klog.V(6).Infof("启动Harbor插件后台任务")

	startCtx, cancel := context.WithCancel(context.Background())
	h.cancelStart = cancel

	go func(meta plugins.Meta) {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				klog.V(6).Infof("Harbor插件后台任务运行中，插件: %s，版本: %s", meta.Name, meta.Version)
				// 这里可以添加定期同步Harbor数据的逻辑
			case <-startCtx.Done():
				klog.V(6).Infof("Harbor插件后台任务退出")
				return
			}
		}
	}(ctx.Meta())

	return nil
}

// StartCron 启动Harbor插件的定时任务（不可阻塞）
func (h *HarborLifecycle) StartCron(ctx plugins.BaseContext, spec string) error {
	klog.V(6).Infof("执行Harbor插件定时任务，表达式: %s", spec)
	return nil
}

// Stop 停止Harbor插件的后台任务
func (h *HarborLifecycle) Stop(ctx plugins.BaseContext) error {
	klog.V(6).Infof("停止Harbor插件后台任务")

	if h.cancelStart != nil {
		h.cancelStart()
		h.cancelStart = nil
	}

	return nil
}
