package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Deployment struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Plan         string                 `json:"plan"`
	InstanceName string                 `json:"instance_name"`
	Domain       string                 `json:"domain"`
	Status       string                 `json:"status"`
	Metrics      map[string]interface{} `json:"metrics"`
	CreatedAt    time.Time              `json:"created_at"`
}

var deployments = []Deployment{
	{ID: uuid.New().String(), UserID: "u001", Plan: "pro", InstanceName: "prod-001", Domain: "user1.openclaw.cn", Status: "running", CreatedAt: time.Now(), Metrics: map[string]interface{}{"cpu_usage": 45.5, "memory_usage": 2.1, "qps": 1234, "response_time": 45}},
	{ID: uuid.New().String(), UserID: "u002", Plan: "pro", InstanceName: "prod-002", Domain: "user2.openclaw.cn", Status: "running", CreatedAt: time.Now(), Metrics: map[string]interface{}{"cpu_usage": 38.2, "memory_usage": 1.8, "qps": 987, "response_time": 52}},
	{ID: uuid.New().String(), UserID: "u003", Plan: "community", InstanceName: "test-001", Domain: "user3.openclaw.cn", Status: "deploying", CreatedAt: time.Now(), Metrics: map[string]interface{}{"cpu_usage": 12.0, "memory_usage": 0.5, "qps": 0, "response_time": 0}},
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin,Content-Type,Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
	})

	api := r.Group("/api/v1")
	{
		api.GET("/deployments", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"data": deployments, "total": len(deployments)})
		})

		api.GET("/deployments/:id", func(c *gin.Context) {
			id := c.Param("id")
			for _, d := range deployments {
				if d.ID == id {
					c.JSON(http.StatusOK, gin.H{"data": d})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		})

		api.POST("/deployments", func(c *gin.Context) {
			var req struct {
				UserID   string `json:"user_id"`
				Plan     string `json:"plan"`
				Domain   string `json:"domain"`
				Username string `json:"username"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			newDeployment := Deployment{
				ID: uuid.New().String(),
				UserID: req.UserID,
				Plan: req.Plan,
				InstanceName: "openclaw-" + req.UserID + "-" + req.Plan,
				Domain: req.Domain,
				Status: "deploying",
				CreatedAt: time.Now(),
				Metrics: map[string]interface{}{"cpu_usage": 0, "memory_usage": 0, "qps": 0, "response_time": 0},
			}
			deployments = append(deployments, newDeployment)
			c.JSON(http.StatusCreated, gin.H{"data": newDeployment, "message": "Deployment created"})
		})

		api.DELETE("/deployments/:id", func(c *gin.Context) {
			id := c.Param("id")
			for i, d := range deployments {
				if d.ID == id {
					deployments = append(deployments[:i], deployments[i+1:]...)
					c.JSON(http.StatusOK, gin.H{"message": "Deployment deleted"})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		})

		api.POST("/deployments/:id/start", func(c *gin.Context) {
			id := c.Param("id")
			for i := range deployments {
				if deployments[i].ID == id {
					deployments[i].Status = "running"
					c.JSON(http.StatusOK, gin.H{"message": "Deployment started"})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		})

		api.POST("/deployments/:id/stop", func(c *gin.Context) {
			id := c.Param("id")
			for i := range deployments {
				if deployments[i].ID == id {
					deployments[i].Status = "stopped"
					c.JSON(http.StatusOK, gin.H{"message": "Deployment stopped"})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		})

		api.GET("/deployments/:id/metrics", func(c *gin.Context) {
			id := c.Param("id")
			for _, d := range deployments {
				if d.ID == id {
					// 模拟实时数据
					metrics := map[string]interface{}{
						"cpu_usage": 40 + float64(time.Now().Second())*0.5,
						"memory_usage": 2.1,
						"disk_usage": 15.3,
						"network_in": 1024,
						"network_out": 512,
						"qps": 1234,
						"response_time": 45,
						"timestamp": time.Now().Unix(),
					}
					c.JSON(http.StatusOK, gin.H{"data": metrics})
					return
				}
			}
			c.JSON(http.StatusNotFound, gin.H{"error": "Deployment not found"})
		})

		api.GET("/deployments/:id/logs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"logs": []string{}})
		})

		// 数据库部署 APIs
		api.GET("/databases", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{ID: uuid.New().String(), Name: "mysql-prod", DatabaseType: "mysql", Version: "8.0", Status: "running", Host: "47.52.xxx.xxx", Port: 3306, MemorySize: 4096, DiskSize: 100},
					{ID: uuid.New().String(), Name: "pg-main", DatabaseType: "postgresql", Version: "14", Status: "running", Host: "47.52.xxx.xxx", Port: 5432, MemorySize: 8192, DiskSize: 200},
				},
				"total": 2,
				"page": 1,
				"page_size": 20,
			})
		})

		api.POST("/databases", func(c *gin.Context) {
			var req struct {
				UserID       string `json:"user_id"`
				Name         string `json:"name"`
				DatabaseType string `json:"database_type"`
				Version      string `json:"version"`
				Plan         string `json:"plan"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			db := gin.H{
				ID:           uuid.New().String(),
				UserID:       req.UserID,
				Name:         req.Name,
				DatabaseType: req.DatabaseType,
				Version:      req.Version,
				Status:       "deploying",
				Host:         "待分配",
				Port:         3306,
				MemorySize:   4096,
				DiskSize:     40,
			}
			c.JSON(http.StatusCreated, gin.H{"data": db, "message": "Database deployment created"})
		})

		api.GET("/databases/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					ID:           c.Param("id"),
					Name:         "mysql-prod",
					DatabaseType: "mysql",
					Version:      "8.0",
					Status:       "running",
					Host:         "47.52.xxx.xxx",
					Port:         3306,
					Username:     "root",
					MemorySize:   4096,
					DiskSize:     100,
				},
			})
		})

		api.DELETE("/databases/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Database deployment deleted"})
		})

		api.GET("/databases/:id/backups", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{ID: uuid.New().String(), BackupID: "backup-20260410080000", Size: "128MB", CreatedAt: "2026-04-10 08:00:00", Status: "completed"},
					{ID: uuid.New().String(), BackupID: "backup-20260409080000", Size: "125MB", CreatedAt: "2026-04-09 08:00:00", Status: "completed"},
				},
			})
		})

		api.POST("/databases/:id/backups", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "Backup created", "data": gin.H{"ID": uuid.New().String(), Status: "creating"}})
		})

		// Docker 部署 APIs
		api.GET("/docker", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{ID: uuid.New().String(), Name: "nginx-web", Image: "nginx:latest", Status: "running", Host: "47.52.xxx.xxx", ExternalPort: 30000, CPU: 1, Memory: 1024},
					{ID: uuid.New().String(), Name: "redis-cache", Image: "redis:7", Status: "running", Host: "47.52.xxx.xxx", ExternalPort: 30001, CPU: 1, Memory: 2048},
					{ID: uuid.New().String(), Name: "grafana-monitor", Image: "grafana/grafana:latest", Status: "stopped", Host: "47.52.xxx.xxx", ExternalPort: 30002, CPU: 2, Memory: 4096},
				},
				"total": 3,
				"page": 1,
				"page_size": 20,
			})
		})

		api.POST("/docker", func(c *gin.Context) {
			var req struct {
				UserID        string            `json:"user_id"`
				Name          string            `json:"name"`
				Image         string            `json:"image"`
				Plan          string            `json:"plan"`
				ContainerPort int               `json:"container_port"`
				Command       string            `json:"command"`
				Environment   map[string]string `json:"environment"`
				Volumes       []string          `json:"volumes"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			docker := gin.H{
				ID:           uuid.New().String(),
				UserID:       req.UserID,
				Name:         req.Name,
				Image:        req.Image,
				Status:       "deploying",
				Host:         "待分配",
				ExternalPort: 30000,
				CPU:          1,
				Memory:       1024,
			}
			c.JSON(http.StatusCreated, gin.H{"data": docker, "message": "Docker deployment created"})
		})

		api.GET("/docker/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					ID:           c.Param("id"),
					Name:         "nginx-web",
					Image:        "nginx:latest",
					Status:       "running",
					Host:         "47.52.xxx.xxx",
					ExternalPort: 30000,
					CPU:          1,
					Memory:       1024,
					ContainerID:  "container-" + uuid.New().String()[:8],
				},
			})
		})

		api.POST("/docker/:id/start", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Docker deployment started"})
		})

		api.POST("/docker/:id/stop", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Docker deployment stopped"})
		})

		api.DELETE("/docker/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Docker deployment deleted"})
		})

		api.GET("/docker/:id/logs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{Time: "2026-04-10 08:00:00", Level: "INFO", Message: "Container started"},
					{Time: "2026-04-10 08:00:01", Level: "INFO", Message: "Application initialized"},
					{Time: "2026-04-10 08:00:02", Level: "INFO", Message: "Listening on port 80"},
				},
			})
		})

		api.GET("/docker/:id/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": gin.H{
					CPU:        45.5,
					Memory:     262144,
					NetworkRX:  10485760,
					NetworkTX:  524288,
					DiskRead:   104857600,
					DiskWrite:  52428800,
				},
			})
		})

		api.GET("/docker/images", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{Name: "nginx", Description: "Web 服务器", Size: "142MB"},
					{Name: "redis", Description: "缓存数据库", Size: "130MB"},
					{Name: "postgres", Description: "PostgreSQL 数据库", Size: "373MB"},
					{Name: "mysql", Description: "MySQL 数据库", Size: "516MB"},
					{Name: "mongo", Description: "MongoDB 数据库", Size: "700MB"},
					{Name: "node", Description: "Node.js 运行时", Size: "1.1GB"},
					{Name: "python", Description: "Python 运行时", Size: "3.5GB"},
					{Name: "grafana", Description: "监控可视化", Size: "325MB"},
					{Name: "prometheus", Description: "监控时序数据库", Size: "188MB"},
					{Name: "minio", Description: "对象存储", Size: "365MB"},
				},
			})
		})

		api.GET("/monitor/system", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "healthy", "uptime": "99.9%"})
		})

		api.GET("/monitor/alerts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"alerts": []gin.H{}})
		})

		// 数字员工 AI 对话 APIs
		api.POST("/ai/chat", func(c *gin.Context) {
			var req struct {
				SessionID string `json:"session_id"`
				UserID    string `json:"user_id"`
				Message   string `json:"message" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			sessionID := req.SessionID
			if sessionID == "" {
				sessionID = uuid.New().String()
			}

			// 模拟 AI 响应
			intent := "general"
			msg := strings.ToLower(req.Message)

			var response string
			if strings.Contains(msg, "你好") || strings.Contains(msg, "hi") || strings.Contains(msg, "hello") {
				response = "你好！我是 ClawOps 数字员工 🤖\n\n我可以帮你：\n• 查询和管理部署实例\n• 创建数据库（MySQL/PostgreSQL）\n• 管理 Docker 容器\n• 查看系统状态和监控\n\n有什么可以帮你的吗？"
				intent = "greeting"
			} else if strings.Contains(msg, "部署") || strings.Contains(msg, "实例") {
				if strings.Contains(msg, "创建") || strings.Contains(msg, "新建") {
					response = "🎉 收到！创建部署实例\n\n请提供：\n1. 套餐：community / pro / enterprise\n2. 实例名称\n3. 域名（可选）"
					intent = "create_deployment"
				} else {
					response = "📋 你的部署实例：\n\n1. 🟢 prod-api (pro) - 运行中\n2. 🟢 test-web (community) - 运行中\n3. 🔴 staging-db (pro) - 已停止\n\n输入"创建部署"可以开通新实例。"
					intent = "list_deployments"
				}
			} else if strings.Contains(msg, "数据库") || strings.Contains(msg, "db") || strings.Contains(msg, "mysql") || strings.Contains(msg, "postgresql") {
				if strings.Contains(msg, "创建") || strings.Contains(msg, "新建") {
					response = "🗄️ 收到！创建数据库\n\n请提供：\n1. 类型：MySQL / PostgreSQL\n2. 版本\n3. 套餐：small / medium / large"
					intent = "create_database"
				} else {
					response = "🗄️ 数据库实例：\n\n1. 🟢 mysql-prod (MySQL 8.0) - 4GB - 运行中\n2. 🟢 pg-main (PostgreSQL 14) - 8GB - 运行中\n\n输入"创建数据库"可以开通新数据库。"
					intent = "list_databases"
				}
			} else if strings.Contains(msg, "docker") || strings.Contains(msg, "容器") {
				if strings.Contains(msg, "创建") || strings.Contains(msg, "新建") {
					response = "🐳 收到！部署 Docker 容器\n\n请提供：\n1. 镜像：nginx / redis / postgres 等\n2. 容器名称\n3. 套餐：small / medium / large"
					intent = "create_docker"
				} else {
					response = "🐳 Docker 容器：\n\n1. 🟢 nginx-web (nginx:latest) - 端口 30000\n2. 🟢 redis-cache (redis:7) - 端口 30001\n\n输入"创建容器"可以部署新容器。"
					intent = "list_docker"
				}
			} else if strings.Contains(msg, "状态") || strings.Contains(msg, "监控") {
				response = "📊 系统状态\n\n状态：✅ 健康\n在线时间：99.9%\n\n资源概览：\n• 部署实例：2\n• 数据库：2\n• Docker 容器：2\n• CPU 使用：45%\n• 内存使用：38%"
				intent = "system_status"
			} else if strings.Contains(msg, "帮助") || strings.Contains(msg, "help") {
				response = "📖 ClawOps 数字员工使用指南\n\n常用命令：\n• \"查看部署实例\" - 列出所有部署\n• \"创建一个 MySQL 数据库\" - 新建数据库\n• \"部署 Nginx 容器\" - 创建 Docker 容器\n• \"系统状态怎么样\" - 查看监控状态"
				intent = "help"
			} else {
				response = fmt.Sprintf("🤔 我理解了你的意思，但需要更具体一点。\n\n你说的：\"%s\"\n\n我可以帮你管理部署、数据库和容器。请试试：\n• \"查看部署实例\"\n• \"创建一个 MySQL 数据库\"\n• \"部署 Nginx 容器\"", req.Message)
				intent = "unknown"
			}

			c.JSON(http.StatusOK, gin.H{
				"id":         uuid.New().String(),
				"session_id": sessionID,
				"role":       "assistant",
				"content":    response,
				"intent":     intent,
			})
		})

		api.GET("/ai/sessions", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{ID: uuid.New().String(), Title: "部署管理对话", LastMessage: "系统状态怎么样", UpdatedAt: time.Now()},
					{ID: uuid.New().String(), Title: "数据库咨询", LastMessage: "创建一个 MySQL", UpdatedAt: time.Now().Add(-3600)},
				},
			})
		})

		api.GET("/ai/messages/:session_id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"data": []gin.H{
					{Role: "user", Content: "你好", CreatedAt: time.Now().Add(-7200)},
					{Role: "assistant", Content: "你好！我是 ClawOps 数字员工", CreatedAt: time.Now().Add(-7199)},
				},
			})
		})

		api.POST("/monitor/alerts/:id/ack", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Alert acknowledged"})
		})

		api.GET("/users/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"username": "admin", "role": "admin", "email": "admin@openclaw.cn"})
		})

		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"users": []gin.H{}})
		})

		api.PUT("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User updated"})
		})

		api.GET("/domains", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"domains": []gin.H{}})
		})

		api.POST("/domains", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "Domain created"})
		})

		api.DELETE("/domains/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Domain deleted"})
		})

		api.GET("/accounts", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"accounts": []gin.H{}})
		})

		api.POST("/accounts", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"message": "Account created"})
		})

		api.PUT("/accounts/:id/reset-password", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Password reset", "password": uuid.New().String()[:16]})
		})
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Starting OpenClaw Deploy API on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
