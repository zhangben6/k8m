package models

import (
	"github.com/weibaohui/k8m/internal/dao"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

// HarborRegistry Harbor仓库配置
type HarborRegistry struct {
	gorm.Model
	Name        string `json:"name" gorm:"uniqueIndex;not null;comment:仓库名称"`
	URL         string `json:"url" gorm:"not null;comment:Harbor地址"`
	Username    string `json:"username" gorm:"comment:用户名"`
	Password    string `json:"password" gorm:"comment:密码"`
	Insecure    bool   `json:"insecure" gorm:"default:false;comment:是否跳过TLS验证"`
	Description string `json:"description" gorm:"comment:描述"`
	IsDefault   bool   `json:"is_default" gorm:"default:false;comment:是否为默认仓库"`
}

// TableName 指定表名
func (HarborRegistry) TableName() string {
	return "harbor_registries"
}

// InitDB 初始化数据库表
func InitDB() error {
	db := dao.DB()
	if err := db.AutoMigrate(&HarborRegistry{}); err != nil {
		klog.Errorf("Harbor插件数据库初始化失败: %v", err)
		return err
	}
	klog.V(6).Infof("Harbor插件数据库初始化成功")
	return nil
}

// UpgradeDB 升级数据库
func UpgradeDB(fromVersion, toVersion string) error {
	klog.V(6).Infof("Harbor插件数据库升级：从 %s 到 %s", fromVersion, toVersion)
	// 这里可以添加版本升级逻辑
	return nil
}

// DropDB 删除数据库表
func DropDB() error {
	db := dao.DB()
	if err := db.Migrator().DropTable(&HarborRegistry{}); err != nil {
		klog.Errorf("Harbor插件删除数据库表失败: %v", err)
		return err
	}
	klog.V(6).Infof("Harbor插件数据库表删除成功")
	return nil
}
