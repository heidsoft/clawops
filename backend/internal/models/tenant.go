package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 租户
type Tenant struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string    `gorm:"type:varchar(100)" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(20);default:active" json:"status"`
	Plan        string    `gorm:"type:varchar(50);default:free" json:"plan"`
	Settings    string    `gorm:"type:text" json:"settings"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Tenant) TableName() string {
	return "tenants"
}

// 用户
type User struct {
	ID           string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID     string     `gorm:"type:varchar(36);index" json:"tenant_id"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex" json:"username"`
	Email        string     `gorm:"type:varchar(255);uniqueIndex" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255)" json:"-"`
	Nickname     string     `gorm:"type:varchar(100)" json:"nickname"`
	Phone        string     `gorm:"type:varchar(20)" json:"phone"`
	Avatar       string     `gorm:"type:varchar(255)" json:"avatar"`
	Role         string     `gorm:"type:varchar(50);default:user" json:"role"`
	Status       string     `gorm:"type:varchar(20);default:active" json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	LastLoginIP  string     `gorm:"type:varchar(50)" json:"last_login_ip"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// 角色
type Role struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID    string    `gorm:"type:varchar(36);index" json:"tenant_id"`
	Name        string    `gorm:"type:varchar(50)" json:"name"`
	Code        string    `gorm:"type:varchar(50)" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Permissions string    `gorm:"type:text" json:"permissions"`
	IsSystem    bool      `gorm:"type:boolean;default:false" json:"is_system"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Role) TableName() string {
	return "roles"
}

// 审计日志
type AuditLog struct {
	ID           string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID     string    `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID       string    `gorm:"type:varchar(36);index" json:"user_id"`
	Username     string    `gorm:"type:varchar(100)" json:"username"`
	Action       string    `gorm:"type:varchar(50);index" json:"action"`
	Resource     string    `gorm:"type:varchar(100);index" json:"resource"`
	ResourceID   string    `gorm:"type:varchar(100)" json:"resource_id"`
	Detail       string    `gorm:"type:text" json:"detail"`
	IP           string    `gorm:"type:varchar(50)" json:"ip"`
	UserAgent    string    `gorm:"type:varchar(500)" json:"user_agent"`
	Status       string    `gorm:"type:varchar(20)" json:"status"`
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	CreatedAt    time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// API Token
type APIToken struct {
	ID         string     `gorm:"primaryKey;type:varchar(36)" json:"id"`
	TenantID   string     `gorm:"type:varchar(36);index" json:"tenant_id"`
	UserID     string     `gorm:"type:varchar(36);index" json:"user_id"`
	Name       string     `gorm:"type:varchar(100)" json:"name"`
	Token      string     `gorm:"type:varchar(64);uniqueIndex" json:"token"`
	SecretHash string     `gorm:"type:varchar(255)" json:"-"`
	Scopes     string     `gorm:"type:varchar(255)" json:"scopes"`
	ExpiresAt  *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	Status     string     `gorm:"type:varchar(20);default:active" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (APIToken) TableName() string {
	return "api_tokens"
}

// 初始化系统数据
func InitSystemData(db *gorm.DB) error {
	// 初始化默认租户
	defaultTenant := Tenant{
		ID:          "tenant-default",
		Name:        "默认租户",
		Code:        "default",
		Description: "系统默认租户",
		Status:      "active",
		Plan:        "pro",
	}
	db.FirstOrCreate(&defaultTenant, Tenant{Code: "default"})

	// 初始化系统角色
	systemRoles := []Role{
		{
			ID:          "role-super-admin",
			TenantID:    "tenant-default",
			Name:        "超级管理员",
			Code:        "super_admin",
			Description: "系统超级管理员，拥有所有权限",
			Permissions: `["*"]`,
			IsSystem:    true,
		},
		{
			ID:          "role-admin",
			TenantID:    "tenant-default",
			Name:        "管理员",
			Code:        "admin",
			Description: "租户管理员，拥有租户内所有权限",
			Permissions: `["deploy:*", "database:*", "docker:*", "monitor:*", "user:*", "setting:*"]`,
			IsSystem:    true,
		},
		{
			ID:          "role-manager",
			TenantID:    "tenant-default",
			Name:        "运维经理",
			Code:        "manager",
			Description: "负责运维工作，可管理部署和监控",
			Permissions: `["deploy:*", "database:*", "docker:*", "monitor:read", "monitor:write"]`,
			IsSystem:    true,
		},
		{
			ID:          "role-user",
			TenantID:    "tenant-default",
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，只有查看权限",
			Permissions: `["deploy:read", "database:read", "docker:read", "monitor:read"]`,
			IsSystem:    true,
		},
		{
			ID:          "role-viewer",
			TenantID:    "tenant-default",
			Name:        "访客",
			Code:        "viewer",
			Description: "只读访客，仅能查看概览",
			Permissions: `["*:read"]`,
			IsSystem:    true,
		},
	}

	for _, role := range systemRoles {
		db.FirstOrCreate(&role, Role{ID: role.ID})
	}

	// 初始化默认管理员用户 (密码: admin123)
	defaultAdmin := User{
		ID:           "user-admin",
		TenantID:     "tenant-default",
		Username:     "admin",
		Email:        "admin@clawops.cn",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy",
		Nickname:     "管理员",
		Role:         "super_admin",
		Status:       "active",
	}
	db.FirstOrCreate(&defaultAdmin, User{Username: "admin"})

	return nil
}

// 记录审计日志
func CreateAuditLog(db *gorm.DB, userID, username, action, resource, resourceID, detail, ip, userAgent, status string) error {
	tenantID := "tenant-default"
	log := &AuditLog{
		ID:        uuid.New().String(),
		TenantID:  tenantID,
		UserID:    userID,
		Username:  username,
		Action:    action,
		Resource:  resource,
		ResourceID: resourceID,
		Detail:    detail,
		IP:        ip,
		UserAgent: userAgent,
		Status:    status,
	}
	return db.Create(log).Error
}

// 获取审计日志
func GetAuditLogs(db *gorm.DB, tenantID, userID, action, resource string, startTime, endTime time.Time, page, pageSize int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := db.Model(&AuditLog{}).Where("tenant_id = ?", tenantID)

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if resource != "" {
		query = query.Where("resource = ?", resource)
	}
	if !startTime.IsZero() {
		query = query.Where("created_at >= ?", startTime)
	}
	if !endTime.IsZero() {
		query = query.Where("created_at <= ?", endTime)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// 检查用户权限
func CheckPermission(userPermissions []string, required string) bool {
	for _, perm := range userPermissions {
		if perm == "*" || perm == required {
			return true
		}
		if perm == required[:len(perm)-1]+"*" {
			return true
		}
	}
	return false
}

// 获取用户角色信息
func GetUserWithRole(db *gorm.DB, userID string) (*User, *Role, error) {
	var user User
	if err := db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, nil, err
	}

	var role Role
	if err := db.Where("code = ? AND (tenant_id = ? OR is_system = ?)", user.Role, user.TenantID, true).First(&role).Error; err != nil {
		return &user, nil, nil
	}

	return &user, &role, nil
}
