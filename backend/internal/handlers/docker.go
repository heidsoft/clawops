package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"openclaw-deploy/internal/models"
	"openclaw-deploy/pkg/aliyun"
	"openclaw-deploy/pkg/database"
)

type DockerHandler struct {
	db *database.Database
}

func NewDockerHandler(db *database.Database) *DockerHandler {
	return &DockerHandler{db: db}
}

// GetDockerDeployments 获取 Docker 部署列表
func (h *DockerHandler) GetDockerDeployments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	deployments, total, err := models.GetDockerDeployments(h.db.GetDB(), page, pageSize)
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

// CreateDockerDeployment 创建 Docker 部署
func (h *DockerHandler) CreateDockerDeployment(c *gin.Context) {
	var req CreateDockerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取配置
	dockerConfig := getDockerConfig(req.Plan)
	password := generatePassword()

	// 创建 ECS 实例
	instanceName := "docker-" + req.Name
	instance, err := aliyun.CreateInstance(aliyun.CreateInstanceArgs{
		InstanceType:    dockerConfig.InstanceType,
		ImageID:         "ubuntu_20_04_x64_20G_alibase_20210120.vhd",
		SecurityGroupID: aliyun.GetSecurityGroupID(),
		InstanceName:    instanceName,
		Bandwidth:       dockerConfig.Bandwidth,
		DiskSize:        dockerConfig.DiskSize,
		Password:        password,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instance: " + err.Error()})
		return
	}

	// 处理环境变量
	var envVars []string
	if req.Environment != nil {
		for k, v := range req.Environment {
			envVars = append(envVars, k+"="+v)
		}
	}
	envJSON, _ := json.Marshal(envVars)

	// 处理挂载卷
	volJSON, _ := json.Marshal(req.Volumes)

	// 创建部署记录
	deployment := &models.DockerDeployment{
		ID:          uuid.New().String(),
		UserID:      req.UserID,
		Name:        req.Name,
		InstanceID:  instance.InstanceId,
		Image:       req.Image,
		Host:        "待分配",
		Port:        req.ContainerPort,
		ExternalPort: dockerConfig.ExternalPort,
		Status:      "deploying",
		CPU:         dockerConfig.CPU,
		Memory:      dockerConfig.Memory,
		DiskSize:    dockerConfig.DiskSize,
		Region:      aliyun.GetRegionID(),
		Command:     req.Command,
		Environment: string(envJSON),
		Volumes:     string(volJSON),
	}

	if err := models.CreateDockerDeployment(h.db.GetDB(), deployment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    deployment,
		"message": "Docker deployment created successfully",
	})
}

// GetDockerDeployment 获取 Docker 部署详情
func (h *DockerHandler) GetDockerDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDockerDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Docker deployment not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": deployment})
}

// StartDockerDeployment 启动容器
func (h *DockerHandler) StartDockerDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDockerDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Docker deployment not found"})
		return
	}

	// 启动 ECS 实例
	if err := aliyun.StartInstance(deployment.InstanceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deployment.Status = "running"
	models.UpdateDockerDeployment(h.db.GetDB(), deployment)

	c.JSON(http.StatusOK, gin.H{"message": "Docker deployment started"})
}

// StopDockerDeployment 停止容器
func (h *DockerHandler) StopDockerDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDockerDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Docker deployment not found"})
		return
	}

	// 停止 ECS 实例
	if err := aliyun.StopInstance(deployment.InstanceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	deployment.Status = "stopped"
	models.UpdateDockerDeployment(h.db.GetDB(), deployment)

	c.JSON(http.StatusOK, gin.H{"message": "Docker deployment stopped"})
}

// DeleteDockerDeployment 删除 Docker 部署
func (h *DockerHandler) DeleteDockerDeployment(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDockerDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Docker deployment not found"})
		return
	}

	// 删除 ECS 实例
	if err := aliyun.DeleteInstance(deployment.InstanceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 删除部署记录
	models.DeleteDockerDeployment(h.db.GetDB(), id)

	c.JSON(http.StatusOK, gin.H{"message": "Docker deployment deleted"})
}

// GetDockerLogs 获取容器日志
func (h *DockerHandler) GetDockerLogs(c *gin.Context) {
	id := c.Param("id")
	lines, _ := strconv.Atoi(c.DefaultQuery("lines", "100"))

	// 模拟日志数据
	logs := []map[string]string{
		{"time": "2026-04-10 08:00:00", "level": "INFO", "message": "Container started"},
		{"time": "2026-04-10 08:00:01", "level": "INFO", "message": "Application initialized"},
		{"time": "2026-04-10 08:00:02", "level": "INFO", "message": "Listening on port " + c.Param("port")},
	}

	// 只返回请求的行数
	if lines < len(logs) {
		logs = logs[len(logs)-lines:]
	}

	c.JSON(http.StatusOK, gin.H{"data": logs})
}

// GetDockerStats 获取容器状态
func (h *DockerHandler) GetDockerStats(c *gin.Context) {
	id := c.Param("id")

	deployment, err := models.GetDockerDeployment(h.db.GetDB(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Docker deployment not found"})
		return
	}

	// 模拟容器状态
	stats := map[string]interface{}{
		"id":         deployment.ContainerID,
		"cpu":        45.5,
		"memory":     1024 * 256, // 256MB
		"memory_limit": 1024 * deployment.Memory,
		"network_rx": 1024 * 1024 * 10,  // 10MB
		"network_tx": 1024 * 512,       // 512KB
		"disk_read":  1024 * 1024 * 100, // 100MB
		"disk_write": 1024 * 1024 * 50,  // 50MB
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetDockerImages 获取常用镜像列表
func (h *DockerHandler) GetDockerImages(c *gin.Context) {
	images := []map[string]string{
		{"name": "nginx", "description": "Web 服务器", "size": "142MB"},
		{"name": "redis", "description": "缓存数据库", "size": "130MB"},
		{"name": "postgres", "description": "PostgreSQL 数据库", "size": "373MB"},
		{"name": "mysql", "description": "MySQL 数据库", "size": "516MB"},
		{"name": "mongo", "description": "MongoDB 数据库", "size": "700MB"},
		{"name": "node", "description": "Node.js 运行时", "size": "1.1GB"},
		{"name": "python", "description": "Python 运行时", "size": "3.5GB"},
		{"name": "grafana", "description": "监控可视化", "size": "325MB"},
		{"name": "prometheus", "description": "监控时序数据库", "size": "188MB"},
		{"name": "minio", "description": "对象存储", "size": "365MB"},
	}

	c.JSON(http.StatusOK, gin.H{"data": images})
}

type CreateDockerRequest struct {
	UserID         string            `json:"user_id" binding:"required"`
	Name           string            `json:"name" binding:"required"`
	Image          string            `json:"image" binding:"required"`
	Plan           string            `json:"plan"`
	ContainerPort  int                `json:"container_port"`
	Command        string            `json:"command"`
	Environment    map[string]string `json:"environment"`
	Volumes        []string          `json:"volumes"`
}

type DockerConfig struct {
	InstanceType    string
	Bandwidth       int
	DiskSize        int
	CPU             int
	Memory          int
	ExternalPort    int
}

func getDockerConfig(plan string) DockerConfig {
	configs := map[string]DockerConfig{
		"small":    {InstanceType: "ecs.n4.small", Bandwidth: 1, DiskSize: 40, CPU: 1, Memory: 1024, ExternalPort: 30000},
		"medium":   {InstanceType: "ecs.n4.large", Bandwidth: 3, DiskSize: 100, CPU: 2, Memory: 4096, ExternalPort: 30001},
		"large":    {InstanceType: "ecs.n4.xlarge", Bandwidth: 5, DiskSize: 200, CPU: 4, Memory: 8192, ExternalPort: 30002},
	}

	if plan == "" {
		plan = "small"
	}

	config, ok := configs[plan]
	if !ok {
		return configs["small"]
	}
	return config
}
