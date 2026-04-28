package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LLM 配置
type LLMConfig struct {
	Provider    string
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float64
}

func getLLMConfig() LLMConfig {
	provider := os.Getenv("LLM_PROVIDER")
	if provider == "" {
		provider = "openai"
	}

	baseURL := os.Getenv("LLM_BASE_URL")
	if baseURL == "" {
		if provider == "ollama" {
			baseURL = "http://localhost:11434"
		} else {
			baseURL = "https://api.openai.com/v1"
		}
	}

	model := os.Getenv("LLM_MODEL")
	if model == "" {
		if provider == "ollama" {
			model = "llama3.2"
		} else {
			model = "gpt-4o-mini"
		}
	}

	return LLMConfig{
		Provider:    provider,
		APIKey:      os.Getenv("OPENAI_API_KEY"),
		BaseURL:     baseURL,
		Model:       model,
		Temperature: 0.7,
	}
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func CallLLM(messages []ChatMessage, config LLMConfig) (string, error) {
	llmMessages := []ChatMessage{
		{Role: "system", Content: `你是一个专业的云运维助手，名叫 ClawOps 数字员工。

## 你的能力
- 查询和管理云资源（部署实例、数据库、容器）
- 用中文回答，友好且专业
- 支持工具调用来获取数据

## 可用工具
当用户询问信息时，返回真实的资源数据。

## 回答风格
- 使用 emoji 增加可读性
- 信息结构化展示
- 重要数据用粗体标注`},
	}

	for _, msg := range messages {
		llmMessages = append(llmMessages, ChatMessage{Role: msg.Role, Content: msg.Content})
	}

	reqBody := map[string]interface{}{
		"model":       config.Model,
		"messages":    llmMessages,
		"temperature": config.Temperature,
	}

	reqJSON, _ := json.Marshal(reqBody)

	endpoint := config.BaseURL + "/chat/completions"

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API error: %s", string(body))
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	choices, ok := result["choices"].([]interface{})
	if ok && len(choices) > 0 {
		choice := choices[0].(map[string]interface{})
		if msg, ok := choice["message"].(map[string]interface{}); ok {
			return msg["content"].(string), nil
		}
	}

	return "", fmt.Errorf("no response content")
}

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
			// 获取模拟告警数据
			c.JSON(http.StatusOK, gin.H{
				"alerts": []gin.H{
					{
						"id": uuid.New().String(),
						"rule_name": "CPU 使用率过高",
						"instance_name": "prod-001",
						"severity": "warning",
						"status": "firing",
						"message": "CPU 使用率超过 80%，当前 85%",
						"triggered_at": time.Now().Add(-30 * time.Minute),
					},
					{
						"id": uuid.New().String(),
						"rule_name": "内存使用率过高",
						"instance_name": "prod-002",
						"severity": "warning",
						"status": "acknowledged",
						"message": "内存使用率超过 85%，当前 87%",
						"triggered_at": time.Now().Add(-2 * time.Hour),
						"acknowledged_at": time.Now().Add(-1 * time.Hour),
					},
				},
				"total": 2,
			})
		})

		api.GET("/monitor/overview", func(c *gin.Context) {
			// 监控概览
			c.JSON(http.StatusOK, gin.H{
				"overview": []gin.H{
					{"instance_id": "i-001", "instance_name": "prod-001", "status": "running", "cpu": 45.5, "memory": 62.3, "disk": 45.2},
					{"instance_id": "i-002", "instance_name": "prod-002", "status": "running", "cpu": 38.2, "memory": 55.8, "disk": 38.1},
					{"instance_id": "i-003", "instance_name": "test-001", "status": "deploying", "cpu": 12.0, "memory": 25.0, "disk": 15.0},
				},
				"timestamp": time.Now(),
			})
		})

		api.GET("/monitor/rules", func(c *gin.Context) {
			// 告警规则列表
			c.JSON(http.StatusOK, gin.H{
				"rules": []gin.H{
					{"id": "rule-1", "name": "CPU 使用率过高", "metric_type": "cpu", "threshold": 80, "severity": "warning", "enabled": true},
					{"id": "rule-2", "name": "CPU 使用率严重", "metric_type": "cpu", "threshold": 95, "severity": "critical", "enabled": true},
					{"id": "rule-3", "name": "内存使用率过高", "metric_type": "memory", "threshold": 85, "severity": "warning", "enabled": true},
					{"id": "rule-4", "name": "磁盘使用率过高", "metric_type": "disk", "threshold": 90, "severity": "warning", "enabled": false},
				},
				"total": 4,
			})
		})

		api.GET("/monitor/stats", func(c *gin.Context) {
			// 告警统计
			c.JSON(http.StatusOK, gin.H{
				"total": 156,
				"firing": 2,
				"acknowledged": 5,
				"resolved": 149,
				"daily_stats": map[string]int{
					"04-22": 12,
					"04-23": 18,
					"04-24": 8,
					"04-25": 15,
					"04-26": 22,
					"04-27": 11,
					"04-28": 5,
				},
				"timestamp": time.Now(),
			})
		})

		api.GET("/monitor/alerts/:id", func(c *gin.Context) {
			// 告警详情
			c.JSON(http.StatusOK, gin.H{
				"id": c.Param("id"),
				"rule_name": "CPU 使用率过高",
				"instance_name": "prod-001",
				"severity": "warning",
				"status": "firing",
				"message": "CPU 使用率超过 80%，当前 85%",
				"triggered_at": time.Now().Add(-30 * time.Minute),
				"metric_data": gin.H{
					"cpu_usage": 85.5,
					"memory_usage": 62.3,
				},
			})
		})

		// Skills 市场 APIs
		api.GET("/skills", func(c *gin.Context) {
			// 模拟 Skills 列表数据
			category := c.Query("category")
			search := c.Query("search")
			page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
			pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

			skills := []gin.H{
				{"id": "skill-deploy", "name": "deploy", "version": "1.0.0", "description": "自动化部署技能", "author": "ClawOps Team", "category": "devops", "icon": "🚀", "stars": 128, "installs": 1024, "is_official": true},
				{"id": "skill-monitor", "name": "monitor", "version": "1.0.0", "description": "监控告警技能", "author": "ClawOps Team", "category": "devops", "icon": "📊", "stars": 96, "installs": 856, "is_official": true},
				{"id": "skill-backup", "name": "backup", "version": "1.0.0", "description": "备份恢复技能", "author": "ClawOps Team", "category": "devops", "icon": "💾", "stars": 72, "installs": 543, "is_official": true},
				{"id": "skill-log", "name": "log", "version": "1.0.0", "description": "日志查询技能", "author": "ClawOps Team", "category": "devops", "icon": "📋", "stars": 65, "installs": 432, "is_official": true},
				{"id": "skill-incident", "name": "incident", "version": "1.0.0", "description": "故障响应技能", "author": "ClawOps Team", "category": "devops", "icon": "🚨", "stars": 58, "installs": 321, "is_official": true},
				{"id": "skill-k8s", "name": "kubernetes-deploy", "version": "1.0.0", "description": "K8s 部署技能", "author": "ClawOps Team", "category": "devops", "icon": "☸️", "stars": 89, "installs": 678, "is_official": true},
				{"id": "skill-mysql", "name": "mysql-backup", "version": "1.0.0", "description": "MySQL 备份技能", "author": "ClawOps Team", "category": "database", "icon": "🐬", "stars": 54, "installs": 287, "is_official": true},
				{"id": "skill-redis", "name": "redis-monitor", "version": "1.0.0", "description": "Redis 监控技能", "author": "ClawOps Team", "category": "database", "icon": "🔴", "stars": 47, "installs": 234, "is_official": true},
				{"id": "skill-nginx", "name": "nginx-config", "version": "1.0.0", "description": "Nginx 配置技能", "author": "社区贡献", "category": "devops", "icon": "🌟", "stars": 43, "installs": 156, "is_official": false},
				{"id": "skill-docker", "name": "docker-optimize", "version": "1.0.0", "description": "Docker 优化技能", "author": "社区贡献", "category": "devops", "icon": "🐳", "stars": 38, "installs": 123, "is_official": false},
			}

			// 过滤
			if category != "" {
				filtered := []gin.H{}
				for _, s := range skills {
					if s["category"] == category {
						filtered = append(filtered, s)
					}
				}
				skills = filtered
			}

			if search != "" {
				filtered := []gin.H{}
				for _, s := range skills {
					name := s["name"].(string)
					desc := s["description"].(string)
					if strings.Contains(strings.ToLower(name), strings.ToLower(search)) || strings.Contains(strings.ToLower(desc), strings.ToLower(search)) {
						filtered = append(filtered, s)
					}
				}
				skills = filtered
			}

			categories := []gin.H{
				{"id": "cat-devops", "name": "devops", "description": "运维自动化", "icon": "🛠️", "count": 6},
				{"id": "cat-database", "name": "database", "description": "数据库管理", "icon": "🗄️", "count": 2},
				{"id": "cat-security", "name": "security", "description": "安全合规", "icon": "🔒", "count": 0},
				{"id": "cat-network", "name": "network", "description": "网络管理", "icon": "🌐", "count": 0},
				{"id": "cat-development", "name": "development", "description": "开发工具", "icon": "💻", "count": 0},
			}

			c.JSON(http.StatusOK, gin.H{
				"skills":       skills,
				"total":        len(skills),
				"page":         page,
				"page_size":    pageSize,
				"total_pages":  1,
				"categories":   categories,
			})
		})

		api.GET("/skills/:id", func(c *gin.Context) {
			id := c.Param("id")
			skill := gin.H{
				"id":          id,
				"name":        "deploy",
				"version":     "1.0.0",
				"description": "自动化部署技能 - 支持 Docker、K8s、VM 部署",
				"author":      "ClawOps Team",
				"category":    "devops",
				"tags":        "devops,deployment,docker,kubernetes",
				"icon":        "🚀",
				"stars":       128,
				"installs":    1024,
				"is_official": true,
				"readme":      "# Deploy Skill\n\n自动化部署服务，支持多种部署方式。\n\n## 支持平台\n- Docker Compose\n- Kubernetes\n- VM (SSH)\n\n## 使用示例\n```\n部署 nginx 到生产环境\n```",
			}
			c.JSON(http.StatusOK, gin.H{"skill": skill})
		})

		api.GET("/my-skills", func(c *gin.Context) {
			// 用户已安装的 Skills
			c.JSON(http.StatusOK, gin.H{
				"user_skills": []gin.H{
					{"id": "us-1", "skill_id": "skill-deploy", "skill_name": "deploy", "version": "1.0.0", "status": "active", "enabled": true},
					{"id": "us-2", "skill_id": "skill-monitor", "skill_name": "monitor", "version": "1.0.0", "status": "active", "enabled": true},
					{"id": "us-3", "skill_id": "skill-backup", "skill_name": "backup", "version": "1.0.0", "status": "active", "enabled": true},
				},
				"total": 3,
			})
		})

		// 企业版 APIs - 用户管理
		api.POST("/auth/login", func(c *gin.Context) {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
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
		})

		api.GET("/users/me", func(c *gin.Context) {
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
		})

		api.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"users": []gin.H{
					{"id": "user-admin", "username": "admin", "nickname": "管理员", "email": "admin@clawops.cn", "role": "super_admin", "status": "active"},
					{"id": "user-002", "username": "operator", "nickname": "运维人员", "email": "operator@clawops.cn", "role": "manager", "status": "active"},
					{"id": "user-003", "username": "viewer", "nickname": "访客", "email": "viewer@clawops.cn", "role": "viewer", "status": "active"},
				},
				"total": 3,
			})
		})

		api.POST("/users", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"user": gin.H{
					"id":       uuid.New().String(),
					"username": "newuser",
					"role":     "user",
					"status":   "active",
				},
				"message": "用户创建成功",
			})
		})

		api.PUT("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "用户更新成功"})
		})

		api.DELETE("/users/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "用户删除成功"})
		})

		// 审计日志
		api.GET("/audit/logs", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"logs": []gin.H{
					{"id": uuid.New().String(), "username": "admin", "action": "login", "resource": "session", "ip": "192.168.1.100", "status": "success", "created_at": time.Now().Add(-1 * time.Hour)},
					{"id": uuid.New().String(), "username": "admin", "action": "create", "resource": "deployment", "resource_id": "dep-001", "ip": "192.168.1.100", "status": "success", "created_at": time.Now().Add(-2 * time.Hour)},
					{"id": uuid.New().String(), "username": "operator", "action": "update", "resource": "monitor", "ip": "192.168.1.101", "status": "success", "created_at": time.Now().Add(-3 * time.Hour)},
				},
				"total": 3,
			})
		})

		api.GET("/audit/stats", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"total_actions": 1256,
				"today_actions":  45,
				"week_actions":   312,
			})
		})

		// API Token
		api.GET("/api-tokens", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"tokens": []gin.H{
					{"id": uuid.New().String(), "name": "CI/CD 集成", "scopes": "deploy:*", "status": "active", "created_at": time.Now().Add(-30 * 24 * time.Hour)},
				},
				"total": 1,
			})
		})

		api.POST("/api-tokens", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{
				"token": gin.H{
					"id":      uuid.New().String(),
					"name":    "新 Token",
					"token":   uuid.New().String(),
					"secret":  uuid.New().String(),
					"status":  "active",
				},
				"message": "Token 创建成功，请妥善保管",
			})
		})

		api.DELETE("/api-tokens/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Token 已撤销"})
		})

		// 角色
		api.GET("/roles", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"roles": []gin.H{
					{"id": "role-super-admin", "name": "超级管理员", "code": "super_admin", "is_system": true, "user_count": 1},
					{"id": "role-admin", "name": "管理员", "code": "admin", "is_system": true, "user_count": 0},
					{"id": "role-manager", "name": "运维经理", "code": "manager", "is_system": true, "user_count": 1},
					{"id": "role-user", "name": "普通用户", "code": "user", "is_system": true, "user_count": 1},
					{"id": "role-viewer", "name": "访客", "code": "viewer", "is_system": true, "user_count": 1},
				},
			})
		})

		// 租户
		api.GET("/tenant", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"tenant": gin.H{
					"id":   "tenant-default",
					"name": "默认租户",
					"plan": "pro",
					"usage": gin.H{
						"users":       4,
						"deployments": 3,
						"databases":   2,
					},
				},
			})
		})

		// 菜单
		api.GET("/menus", func(c *gin.Context) {
			role := c.Query("role")
			isAdmin := role == "super_admin" || role == "admin"
			menus := []gin.H{
				{"id": "dashboard", "name": "概览", "icon": "fa-home", "path": "/"},
				{"id": "deployments", "name": "部署实例", "icon": "fa-server", "path": "/deployments"},
				{"id": "databases", "name": "数据库", "icon": "fa-database", "path": "/databases"},
				{"id": "docker", "name": "Docker容器", "icon": "fa-box", "path": "/docker"},
				{"id": "ai", "name": "数字员工", "icon": "fa-robot", "path": "/ai"},
				{"id": "skills", "name": "Skill市场", "icon": "fa-plug", "path": "/skills"},
				{"id": "monitoring", "name": "监控告警", "icon": "fa-chart-line", "path": "/monitoring"},
			}
			if isAdmin {
				adminMenus := []gin.H{
					{"id": "users", "name": "用户管理", "icon": "fa-users", "path": "/users"},
					{"id": "audit", "name": "审计日志", "icon": "fa-history", "path": "/audit"},
					{"id": "settings", "name": "系统设置", "icon": "fa-cog", "path": "/settings"},
				}
				menus = append(menus, adminMenus...)
			}
			c.JSON(http.StatusOK, gin.H{"menus": menus})
		})

		// 数字员工 AI 对话 APIs (LLM 版本)
			var req struct {
				SkillID string `json:"skill_id"`
				UserID  string `json:"user_id"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message":     "安装成功",
				"user_skill": gin.H{"id": uuid.New().String(), "skill_id": req.SkillID, "status": "active"},
			})
		})

		api.POST("/skills/:id/uninstall", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "卸载成功"})
		})

		api.POST("/skills/:id/star", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
		})

		// 数字员工 AI 对话 APIs (LLM 版本)
		api.POST("/ai/chat", func(c *gin.Context) {
			var req struct {
				SessionID string `json:"session_id"`
				UserID    string `json:"user_id"`
				Message   string `json:"message" binding:"required"`
				UseLLM    bool   `json:"use_llm"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			sessionID := req.SessionID
			if sessionID == "" {
				sessionID = uuid.New().String()
			}

			// 如果启用 LLM (默认为 true)
			useLLM := req.UseLLM
			if !useLLM && req.UseLLM {
				// 如果前端没有传 use_llm，默认尝试使用 LLM（如果有配置）
				config := getLLMConfig()
				if config.APIKey != "" || config.Provider == "ollama" {
					useLLM = true
				}
			}

			if useLLM {
				config := getLLMConfig()
				messages := []ChatMessage{
					{Role: "user", Content: req.Message},
				}
				resp, err := CallLLM(messages, config)
				if err != nil {
					c.JSON(http.StatusOK, gin.H{
						"id":         uuid.New().String(),
						"session_id": sessionID,
						"role":       "assistant",
						"content":    "抱歉，AI 服务暂时不可用。请稍后再试。\n\n错误：" + err.Error(),
						"intent":     "error",
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"id":         uuid.New().String(),
					"session_id": sessionID,
					"role":       "assistant",
					"content":    resp,
					"intent":     "llm",
				})
				return
			}

			// 规则引擎响应
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
					response = "📋 你的部署实例：\n\n1. 🟢 prod-api (pro) - 运行中\n2. 🟢 test-web (community) - 运行中\n3. 🔴 staging-db (pro) - 已停止\n\n输入「创建部署」可以开通新实例。"
					intent = "list_deployments"
				}
			} else if strings.Contains(msg, "数据库") || strings.Contains(msg, "db") || strings.Contains(msg, "mysql") || strings.Contains(msg, "postgresql") {
				if strings.Contains(msg, "创建") || strings.Contains(msg, "新建") {
					response = "🗄️ 收到！创建数据库\n\n请提供：\n1. 类型：MySQL / PostgreSQL\n2. 版本\n3. 套餐：small / medium / large"
					intent = "create_database"
				} else {
					response = "🗄️ 数据库实例：\n\n1. 🟢 mysql-prod (MySQL 8.0) - 4GB - 运行中\n2. 🟢 pg-main (PostgreSQL 14) - 8GB - 运行中\n\n输入「创建数据库」可以开通新数据库。"
					intent = "list_databases"
				}
			} else if strings.Contains(msg, "docker") || strings.Contains(msg, "容器") {
				if strings.Contains(msg, "创建") || strings.Contains(msg, "新建") {
					response = "🐳 收到！部署 Docker 容器\n\n请提供：\n1. 镜像：nginx / redis / postgres 等\n2. 容器名称\n3. 套餐：small / medium / large"
					intent = "create_docker"
				} else {
					response = "🐳 Docker 容器：\n\n1. 🟢 nginx-web (nginx:latest) - 端口 30000\n2. 🟢 redis-cache (redis:7) - 端口 30001\n\n输入「创建容器」可以部署新容器。"
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
