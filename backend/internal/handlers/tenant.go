package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// 用户登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 简化验证 - 实际应该查数据库
	if req.Username == "admin" && req.Password == "admin123" {
		token := uuid.New().String()
		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"user": gin.H{
				"id":       "user-admin",
				"username": "admin",
				"nickname": "管理员",
				"email":    "admin@clawops.cn",
				"role":     "super_admin",
			},
		})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
}

// 获取当前用户
func GetCurrentUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"id":       "user-admin",
		"username": "admin",
		"nickname": "管理员",
		"email":    "admin@clawops.cn",
		"role":     "super_admin",
		"tenant": gin.H{
			"id":   "tenant-default",
			"name": "默认租户",
			"plan": "pro",
		},
	})
}

// 获取用户列表
func GetUsers(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 模拟用户数据
	users := []gin.H{
		{"id": "user-admin", "username": "admin", "nickname": "管理员", "email": "admin@clawops.cn", "role": "super_admin", "status": "active"},
		{"id": "user-002", "username": "operator", "nickname": "运维人员", "email": "operator@clawops.cn", "role": "manager", "status": "active"},
		{"id": "user-003", "username": "viewer", "nickname": "访客", "email": "viewer@clawops.cn", "role": "viewer", "status": "active"},
	}

	if tenantID != "" {
		filtered := []gin.H{}
		for _, u := range users {
			filtered = append(filtered, u)
		}
		users = filtered
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(users) {
		users = []gin.H{}
	} else if end > len(users) {
		users = users[start:]
	} else {
		users = users[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total":       3,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": 1,
	})
}

// 创建用户
func CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 密码哈希
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := gin.H{
		"id":           uuid.New().String(),
		"username":     req.Username,
		"email":        req.Email,
		"nickname":     req.Nickname,
		"role":         req.Role,
		"status":       "active",
		"password_set": true,
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "用户创建成功",
	})
}

// 更新用户
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	delete(updates, "id")
	delete(updates, "password")

	c.JSON(http.StatusOK, gin.H{
		"user":    gin.H{"id": id}.Merge(updates),
		"message": "用户更新成功",
	})
}

// 删除用户
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
}

// 获取审计日志
func GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	action := c.Query("action")
	resource := c.Query("resource")

	// 模拟审计日志数据
	logs := []gin.H{
		{
			"id":         uuid.New().String(),
			"user_id":    "user-admin",
			"username":   "admin",
			"action":     "login",
			"resource":   "session",
			"detail":     `{"ip": "192.168.1.100"}`,
			"ip":         "192.168.1.100",
			"status":     "success",
			"created_at": time.Now().Add(-1 * time.Hour),
		},
		{
			"id":         uuid.New().String(),
			"user_id":    "user-admin",
			"username":   "admin",
			"action":     "create",
			"resource":   "deployment",
			"resource_id": "dep-001",
			"detail":     `{"name": "prod-api", "plan": "pro"}`,
			"ip":         "192.168.1.100",
			"status":     "success",
			"created_at": time.Now().Add(-2 * time.Hour),
		},
		{
			"id":         uuid.New().String(),
			"user_id":    "user-002",
			"username":   "operator",
			"action":     "update",
			"resource":   "monitor",
			"detail":     `{"rule_id": "rule-1", "threshold": 85}`,
			"ip":         "192.168.1.101",
			"status":     "success",
			"created_at": time.Now().Add(-3 * time.Hour),
		},
		{
			"id":         uuid.New().String(),
			"user_id":    "user-admin",
			"username":   "admin",
			"action":     "delete",
			"resource":   "docker",
			"resource_id": "cnt-003",
			"detail":     `{"name": "test-container"}`,
			"ip":         "192.168.1.100",
			"status":     "success",
			"created_at": time.Now().Add(-5 * time.Hour),
		},
	}

	// 过滤
	filtered := logs
	if action != "" {
		filtered = filterLogs(filtered, "action", action)
	}
	if resource != "" {
		filtered = filterLogs(filtered, "resource", resource)
	}

	total := len(filtered)
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > total {
		filtered = []gin.H{}
	} else if end > total {
		filtered = filtered[start:]
	} else {
		filtered = filtered[start:end]
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":        filtered,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + pageSize - 1) / pageSize,
	})
}

func filterLogs(logs []gin.H, key, value string) []gin.H {
	result := []gin.H{}
	for _, log := range logs {
		if log[key] == value {
			result = append(result, log)
		}
	}
	return result
}

// 获取审计统计
func GetAuditStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"total_actions": 1256,
		"by_action": gin.H{
			"login":    423,
			"create":   312,
			"update":   298,
			"delete":   123,
			"logout":   100,
		},
		"by_resource": gin.H{
			"deployment": 456,
			"database":   234,
			"docker":     198,
			"monitor":    178,
			"user":       190,
		},
		"today_actions": 45,
		"week_actions":  312,
	})
}

// 创建 API Token
func CreateAPIToken(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Scopes    string `json:"scopes"`
		ExpiresIn int    `json:"expires_in"` // 天数，0 表示永不过期
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 生成 token
	tokenBytes := make([]byte, 32)
	rand.Read(tokenBytes)
	token := hex.EncodeToString(tokenBytes)

	secretBytes := make([]byte, 16)
	rand.Read(secretBytes)
	secret := hex.EncodeToString(secretBytes)

	var expiresAt *time.Time
	if req.ExpiresIn > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresIn)
		expiresAt = &t
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": gin.H{
			"id":          uuid.New().String(),
			"name":        req.Name,
			"token":       token,
			"secret":      secret,
			"scopes":      req.Scopes,
			"expires_at":  expiresAt,
			"status":      "active",
			"created_at":  time.Now(),
		},
		"message": "API Token 创建成功，请妥善保管 Token 和 Secret",
	})
}

// 获取 API Token 列表
func GetAPITokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tokens": []gin.H{
			{
				"id":         uuid.New().String(),
				"name":       "CI/CD 集成",
				"scopes":     "deploy:*,database:read",
				"status":     "active",
				"last_used":  time.Now().Add(-1 * time.Hour),
				"created_at": time.Now().Add(-30 * 24 * time.Hour),
			},
			{
				"id":         uuid.New().String(),
				"name":       "监控数据导出",
				"scopes":     "monitor:read",
				"status":     "active",
				"last_used":  time.Now().Add(-2 * time.Hour),
				"created_at": time.Now().Add(-7 * 24 * time.Hour),
			},
		},
		"total": 2,
	})
}

// 撤销 API Token
func RevokeAPIToken(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Token 已撤销",
	})
}

// 获取角色列表
func GetRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"roles": []gin.H{
			{
				"id":          "role-super-admin",
				"name":        "超级管理员",
				"code":        "super_admin",
				"description": "系统超级管理员，拥有所有权限",
				"permissions": []string{"*"},
				"is_system":   true,
				"user_count":  1,
			},
			{
				"id":          "role-admin",
				"name":        "管理员",
				"code":        "admin",
				"description": "租户管理员，拥有租户内所有权限",
				"permissions": []string{"deploy:*", "database:*", "docker:*", "monitor:*", "user:*", "setting:*"},
				"is_system":   true,
				"user_count":  0,
			},
			{
				"id":          "role-manager",
				"name":        "运维经理",
				"code":        "manager",
				"description": "负责运维工作，可管理部署和监控",
				"permissions": []string{"deploy:*", "database:*", "docker:*", "monitor:read", "monitor:write"},
				"is_system":   true,
				"user_count":  1,
			},
			{
				"id":          "role-user",
				"name":        "普通用户",
				"code":        "user",
				"description": "普通用户，只有查看权限",
				"permissions": []string{"deploy:read", "database:read", "docker:read", "monitor:read"},
				"is_system":   true,
				"user_count":  1,
			},
			{
				"id":          "role-viewer",
				"name":        "访客",
				"code":        "viewer",
				"description": "只读访客，仅能查看概览",
				"permissions": []string{"*:read"},
				"is_system":   true,
				"user_count":  1,
			},
		},
	})
}

// 获取租户信息
func GetTenant(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"tenant": gin.H{
			"id":          "tenant-default",
			"name":        "默认租户",
			"code":        "default",
			"description": "系统默认租户",
			"status":      "active",
			"plan":        "pro",
			"settings": gin.H{
				"max_users":       100,
				"max_deployments": 50,
				"max_databases":   20,
			},
			"usage": gin.H{
				"users":       4,
				"deployments": 3,
				"databases":   2,
				"docker":      3,
			},
			"created_at": time.Now().Add(-30 * 24 * time.Hour),
		},
	})
}

// 获取菜单
func GetMenus(c *gin.Context) {
	role := c.Query("role")

	// 根据角色返回不同菜单
	allMenus := []gin.H{
		{
			"id":    "dashboard",
			"name":  "概览",
			"icon":  "fa-home",
			"path":  "/",
			"action": "read",
		},
		{
			"id":    "deployments",
			"name":  "部署实例",
			"icon":  "fa-server",
			"path":  "/deployments",
			"action": "deploy:read",
		},
		{
			"id":    "databases",
			"name":  "数据库",
			"icon":  "fa-database",
			"path":  "/databases",
			"action": "database:read",
		},
		{
			"id":    "docker",
			"name":  "Docker容器",
			"icon":  "fa-box",
			"path":  "/docker",
			"action": "docker:read",
		},
		{
			"id":    "ai",
			"name":  "数字员工",
			"icon":  "fa-robot",
			"path":  "/ai",
			"action": "ai:read",
		},
		{
			"id":    "skills",
			"name":  "Skill市场",
			"icon":  "fa-plug",
			"path":  "/skills",
			"action": "skills:read",
		},
		{
			"id":    "monitoring",
			"name":  "监控告警",
			"icon":  "fa-chart-line",
			"path":  "/monitoring",
			"action": "monitor:read",
		},
	}

	// 管理员及以上角色可以看到管理菜单
	if role == "super_admin" || role == "admin" {
		adminMenus := []gin.H{
			{
				"id":    "users",
				"name":  "用户管理",
				"icon":  "fa-users",
				"path":  "/users",
				"action": "user:read",
			},
			{
				"id":    "roles",
				"name":  "角色权限",
				"icon":  "fa-shield-alt",
				"path":  "/roles",
				"action": "user:write",
			},
			{
				"id":    "audit",
				"name":  "审计日志",
				"icon":  "fa-history",
				"path":  "/audit",
				"action": "audit:read",
			},
			{
				"id":    "settings",
				"name":  "系统设置",
				"icon":  "fa-cog",
				"path":  "/settings",
				"action": "setting:read",
			},
		}
		allMenus = append(allMenus, adminMenus...)
	}

	c.JSON(http.StatusOK, gin.H{"menus": allMenus})
}
