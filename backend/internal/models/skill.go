package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Skill 定义
type Skill struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	Version     string    `gorm:"type:varchar(20)" json:"version"`         // semver
	Description string    `gorm:"type:text" json:"description"`             // 技能描述
	Author      string    `gorm:"type:varchar(100)" json:"author"`         // 作者
	AuthorURL   string    `gorm:"type:varchar(255)" json:"author_url"`     // 作者链接
	Category    string    `gorm:"type:varchar(50);index" json:"category"`  // 分类
	Tags        string    `gorm:"type:varchar(255)" json:"tags"`           // 标签，逗号分隔
	Readme      string    `gorm:"type:text" json:"readme"`                 // 详细说明
	SourceURL   string    `gorm:"type:varchar(255)" json:"source_url"`     // 源码地址
	HomepageURL string    `gorm:"type:varchar(255)" json:"homepage_url"`   // 主页
	Icon        string    `gorm:"type:varchar(100)" json:"icon"`           // 图标 emoji
	Stars       int       `gorm:"type:int;default:0" json:"stars"`       // 点赞数
	Installs    int       `gorm:"type:int;default:0" json:"installs"`     // 安装数
	IsOfficial  bool      `gorm:"type:boolean;default:false" json:"is_official"`
	IsBuiltin   bool      `gorm:"type:boolean;default:false" json:"is_builtin"` // 内置技能
	IsPremium   bool      `gorm:"type:boolean;default:false" json:"is_premium"`  // 付费技能
	Enabled     bool      `gorm:"type:boolean;default:true" json:"enabled"`       // 是否启用
	Manifest    string    `gorm:"type:text" json:"manifest"`               // SKILL.md 内容
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Skill) TableName() string {
	return "skills"
}

// 用户安装的 Skill
type UserSkill struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID    string    `gorm:"type:varchar(36);index" json:"user_id"`
	SkillID   string    `gorm:"type:varchar(36);index" json:"skill_id"`
	SkillName string    `gorm:"type:varchar(100)" json:"skill_name"`
	Version   string    `gorm:"type:varchar(20)" json:"version"`    // 安装时的版本
	Config    string    `gorm:"type:text" json:"config"`           // JSON 配置
	Enabled   bool      `gorm:"type:boolean;default:true" json:"enabled"`
	Status    string    `gorm:"type:varchar(20)" json:"status"`   // active/error/need_update
	InstalledAt time.Time `gorm:"autoCreateTime" json:"installed_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (UserSkill) TableName() string {
	return "user_skills"
}

// Skill 分类
type SkillCategory struct {
	ID          string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Icon        string `gorm:"type:varchar(50)" json:"icon"`
	SortOrder   int    `gorm:"type:int;default:0" json:"sort_order"`
	Count       int    `gorm:"-" json:"count"` // 非数据库字段
}

func (SkillCategory) TableName() string {
	return "skill_categories"
}

// 初始化内置 Skills
func InitBuiltinSkills(db *gorm.DB) error {
	builtinSkills := []Skill{
		{
			ID:          "skill-deploy",
			Name:        "deploy",
			Version:     "1.0.0",
			Description: "自动化部署技能 - 支持 Docker、K8s、VM 部署到阿里云、AWS 等平台",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "devops,deployment,docker,kubernetes,aliyun",
			Icon:        "🚀",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Deploy Skill

## Overview
自动化部署服务，支持多种部署方式。

## When to Use
- 需要部署新服务
- 更新现有服务版本
- 回滚到上一个版本

## Supported Platforms
- Docker Compose
- Kubernetes
- VM (SSH)
`,
		},
		{
			ID:          "skill-monitor",
			Name:        "monitor",
			Version:     "1.0.0",
			Description: "监控告警技能 - 监控 CPU、内存、磁盘、网络，异常时告警",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "devops,monitoring,alerting",
			Icon:        "📊",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Monitor Skill

## Overview
实时监控系统资源和应用状态，支持多种告警规则。

## When to Use
- 查看系统状态
- 配置告警规则
- 分析性能问题
`,
		},
		{
			ID:          "skill-backup",
			Name:        "backup",
			Version:     "1.0.0",
			Description: "备份恢复技能 - 自动备份数据库和文件，支持定时备份和一键恢复",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "devops,backup,recovery",
			Icon:        "💾",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Backup Skill

## Overview
自动化备份和恢复系统。

## When to Use
- 备份数据库
- 恢复数据
- 查看备份历史
`,
		},
		{
			ID:          "skill-log",
			Name:        "log",
			Version:     "1.0.0",
			Description: "日志查询技能 - 收集和分析应用日志，支持关键词搜索",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "devops,logging,troubleshooting",
			Icon:        "📋",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Log Skill

## Overview
日志收集和查询系统。

## When to Use
- 查看应用日志
- 搜索错误信息
- 分析问题根因
`,
		},
		{
			ID:          "skill-k8s-deploy",
			Name:        "kubernetes-deploy",
			Version:     "1.0.0",
			Description: "Kubernetes 部署技能 - 专门用于 K8s 环境的部署和运维",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "kubernetes,k8s,devops,deployment",
			Icon:        "☸️",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Kubernetes Deploy Skill

## Overview
Kubernetes 专用的部署和管理技能。

## When to Use
- 部署到 Kubernetes
- 管理 Pod/Service/Deployment
- 扩缩容
`,
		},
		{
			ID:          "skill-mysql-backup",
			Name:        "mysql-backup",
			Version:     "1.0.0",
			Description: "MySQL 数据库备份技能 - 自动备份 MySQL，支持增量备份",
			Author:      "ClawOps Team",
			Category:    "database",
			Tags:        "mysql,database,backup",
			Icon:        "🐬",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# MySQL Backup Skill

## Overview
MySQL 数据库的专业备份工具。

## When to Use
- 备份 MySQL 数据库
- 恢复 MySQL 数据
- 检查备份完整性
`,
		},
		{
			ID:          "skill-redis-monitor",
			Name:        "redis-monitor",
			Version:     "1.0.0",
			Description: "Redis 监控技能 - 监控 Redis 内存、连接数、命中率等指标",
			Author:      "ClawOps Team",
			Category:    "database",
			Tags:        "redis,monitoring,cache",
			Icon:        "🔴",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Redis Monitor Skill

## Overview
Redis 缓存服务的监控和管理。

## When to Use
- 监控 Redis 状态
- 查看内存使用
- 分析缓存命中率
`,
		},
		{
			ID:          "skill-incident",
			Name:        "incident",
			Version:     "1.0.0",
			Description: "故障响应技能 - 自动化故障定位和应急响应",
			Author:      "ClawOps Team",
			Category:    "devops",
			Tags:        "incident,on-call,devops",
			Icon:        "🚨",
			IsOfficial:  true,
			IsBuiltin:   true,
			Manifest: `# Incident Skill

## Overview
故障检测和应急响应系统。

## When to Use
- 故障告警响应
- 问题定位分析
- 协调故障处理
`,
		},
	}

	// 初始化分类
	categories := []SkillCategory{
		{ID: "cat-devops", Name: "devops", Description: "运维自动化", Icon: "🛠️", SortOrder: 1},
		{ID: "cat-database", Name: "database", Description: "数据库管理", Icon: "🗄️", SortOrder: 2},
		{ID: "cat-security", Name: "security", Description: "安全合规", Icon: "🔒", SortOrder: 3},
		{ID: "cat-network", Name: "network", Description: "网络管理", Icon: "🌐", SortOrder: 4},
		{ID: "cat-development", Name: "development", Description: "开发工具", Icon: "💻", SortOrder: 5},
	}

	for _, cat := range categories {
		var existing SkillCategory
		if err := db.FirstOrCreate(&existing, SkillCategory{ID: cat.ID}).Error; err != nil {
			return err
		}
	}

	// 初始化内置 Skills
	for _, skill := range builtinSkills {
		var existing Skill
		if err := db.FirstOrCreate(&existing, Skill{ID: skill.ID}).Error; err != nil {
			return err
		}
	}

	return nil
}

// 获取 Skills 列表（分页、筛选）
func GetSkills(db *gorm.DB, category, search string, page, pageSize int) ([]Skill, int64, error) {
	var skills []Skill
	var total int64

	query := db.Model(&Skill{}).Where("enabled = ?", true)

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR tags LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// 统计总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Order("stars DESC, installs DESC").Offset(offset).Limit(pageSize).Find(&skills).Error

	return skills, total, err
}

// 获取用户安装的 Skills
func GetUserSkills(db *gorm.DB, userID string) ([]UserSkill, error) {
	var userSkills []UserSkill
	err := db.Where("user_id = ?", userID).Order("installed_at DESC").Find(&userSkills).Error
	return userSkills, err
}

// 安装 Skill
func InstallSkill(db *gorm.DB, userID, skillID string) (*UserSkill, error) {
	var skill Skill
	if err := db.First(&skill, "id = ? AND enabled = ?", skillID, true).Error; err != nil {
		return nil, err
	}

	// 检查是否已安装
	var existing UserSkill
	err := db.Where("user_id = ? AND skill_id = ?", userID, skillID).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("skill already installed")
	}

	userSkill := &UserSkill{
		ID:          uuid.New().String(),
		UserID:      userID,
		SkillID:     skill.ID,
		SkillName:   skill.Name,
		Version:     skill.Version,
		Enabled:     true,
		Status:      "active",
		InstalledAt: time.Now(),
	}

	if err := db.Create(userSkill).Error; err != nil {
		return nil, err
	}

	// 更新安装数
	db.Model(&Skill{}).Where("id = ?", skillID).Update("installs", gorm.Expr("installs + 1"))

	return userSkill, nil
}

// 卸载 Skill
func UninstallSkill(db *gorm.DB, userID, userSkillID string) error {
	var userSkill UserSkill
	if err := db.Where("id = ? AND user_id = ?", userSkillID, userID).First(&userSkill).Error; err != nil {
		return err
	}

	if err := db.Delete(&userSkill).Error; err != nil {
		return err
	}

	// 更新安装数
	db.Model(&Skill{}).Where("id = ?", userSkill.SkillID).Update("installs", gorm.Expr("installs - 1"))

	return nil
}
