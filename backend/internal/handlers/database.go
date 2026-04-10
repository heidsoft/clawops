package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"openclaw-deploy/internal/models"
	"openclaw-deploy/pkg/aliyun"
	"openclaw-deploy/pkg/database"
)

type DatabaseHandler struct {
	db *database.Database
}

func NewDatabaseHandler(db *database.Database) *DatabaseHandler {
	return &DatabaseHandler{db: db}
}

// GetDatabaseDeployments 获取数据库部署列表
func (h *DatabaseHandler) GetDatabaseDeployments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	deployments, total, err := models.GetDatabaseDeployments(h.db.GetDB(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      deployments,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateDatabaseDeployment 创建数据库部署
func (h *DatabaseHandler) CreateDatabaseDeployment(c *gin.Context) {
	var req CreateDatabaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取数据库配置
	dbConfig := getDatabaseConfig(req.DatabaseType, req.Version, req.Plan)
	password := generatePassword()

	// 创建 ECS 实例
	instanceName := "db-" + req.Name + "-" + req.DatabaseType
	instance, err := aliyun.CreateInstance(aliyun.CreateInstanceArgs{
		InstanceType:    dbConfig.InstanceType,
		ImageID:         "ubuntu_20_04_x64_20G_alibase_20210120.vhd",
		SecurityGroupID: aliyun.GetSecurityGroupID(),
		InstanceName:    instanceName,
		Bandwidth:       dbConfig.Bandwidth,
		DiskSize:        dbConfig.DiskSize,
		Password:        password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instance: " + err.Error()})
		return
	}

	// 创建数据库记录
	deployment := &models.DatabaseDeployment{
		ID:           uuid.New().String(),
		UserID:       req.UserID,
		Name:         req.Name,
		DatabaseType: req.DatabaseType,
		Version:      req.Version,
		InstanceID:   instance.InstanceId,
		Host:         "待分配",
		Port:         dbConfig.Port,
		Username:     dbConfig.DefaultUsername,
		Password:     password,
		Status:       "deploying",
		DiskSize:     dbConfig.DiskSize,
		MemorySize:   dbConfig.MemorySize,
		Region:       aliyun.GetRegionID(),
	}

	if err := models.CreateDatabaseDeployment(h.db.GetDB(), deployment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    deployment,
		"message": "Database deployment created successfully",
	})
}

// GetDatabaseDeployment 获取数据库部署详情
func (h *DatabaseHandler) GetDatabaseDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDatabaseDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database deployment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deployment})
}

// DeleteDatabaseDeployment 删除数据库部署
func (h *DatabaseHandler) DeleteDatabaseDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDatabaseDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database deployment not found"})
		return
	}

	// 删除 ECS 实例
	if err := aliyun.DeleteInstance(deployment.InstanceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除数据库记录
	models.DeleteDatabaseDeployment(h.db.GetDB(), id)

	c.JSON(http.StatusOK, gin.H{"message": "Database deployment deleted"})
}

// GetDatabaseBackups 获取数据库备份列表
func (h *DatabaseHandler) GetDatabaseBackups(c *gin.Context) {
	id := c.Param("id")

	// 模拟备份数据
	backups := []map[string]interface{}{
		{
			"id":         uuid.New().String(),
			"backup_id": "backup-" + time.Now().Format("20060102150405"),
			"size":       "128MB",
			"created_at": time.Now().Add(-24 * time.Hour).Format("2006-01-02 15:04:05"),
			"status":     "completed",
		},
		{
			"id":         uuid.New().String(),
			"backup_id": "backup-" + time.Now().Add(-48*time.Hour).Format("20060102150405"),
			"size":       "125MB",
			"created_at": time.Now().Add(-48 * time.Hour).Format("2006-01-02 15:04:05"),
			"status":     "completed",
		},
	}

	c.JSON(http.StatusOK, gin.H{"data": backups})
}

// CreateDatabaseBackup 创建数据库备份
func (h *DatabaseHandler) CreateDatabaseBackup(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDatabaseDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database deployment not found"})
		return
	}

	// 模拟备份创建
	backup := map[string]interface{}{
		"id":         uuid.New().String(),
		"backup_id": "backup-" + time.Now().Format("20060102150405"),
		"database_id": deployment.ID,
		"size":       "0MB",
		"created_at": time.Now().Format("2006-01-02 15:04:05"),
		"status":     "creating",
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    backup,
		"message": "Backup creation started",
	})
}

type CreateDatabaseRequest struct {
	UserID       string `json:"user_id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	DatabaseType string `json:"database_type" binding:"required"`
	Version      string `json:"version"`
	Plan         string `json:"plan"`
}

type DBConfig struct {
	InstanceType    string
	Bandwidth       int
	DiskSize        int
	MemorySize      int
	Port            int
	DefaultUsername string
}

func getDatabaseConfig(dbType, version, plan string) DBConfig {
	configs := map[string]map[string]DBConfig{
		"mysql": {
			"small":    {InstanceType: "ecs.n4.small", Bandwidth: 1, DiskSize: 40, MemorySize: 4096, Port: 3306, DefaultUsername: "root"},
			"medium":   {InstanceType: "ecs.n4.large", Bandwidth: 3, DiskSize: 100, MemorySize: 8192, Port: 3306, DefaultUsername: "root"},
			"large":    {InstanceType: "ecs.n4.xlarge", Bandwidth: 5, DiskSize: 200, MemorySize: 16384, Port: 3306, DefaultUsername: "root"},
		},
		"postgresql": {
			"small":    {InstanceType: "ecs.n4.small", Bandwidth: 1, DiskSize: 40, MemorySize: 4096, Port: 5432, DefaultUsername: "postgres"},
			"medium":   {InstanceType: "ecs.n4.large", Bandwidth: 3, DiskSize: 100, MemorySize: 8192, Port: 5432, DefaultUsername: "postgres"},
			"large":    {InstanceType: "ecs.n4.xlarge", Bandwidth: 5, DiskSize: 200, MemorySize: 16384, Port: 5432, DefaultUsername: "postgres"},
		},
	}

	// 默认 small 套餐
	if plan == "" {
		plan = "small"
	}

	dbTypeConfigs, ok := configs[dbType]
	if !ok {
		return configs["mysql"]["small"]
	}

	config, ok := dbTypeConfigs[plan]
	if !ok {
		return dbTypeConfigs["small"]
	}

	// 设置版本
	if version == "" {
		if dbType == "mysql" {
			version = "8.0"
		} else {
			version = "14"
		}
	}

	return config
}
