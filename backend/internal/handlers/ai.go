package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"openclaw-deploy/internal/models"
	"openclaw-deploy/pkg/database"
)

type AIHandler struct {
	db *database.Database
}

func NewAIHandler(db *database.Database) *AIHandler {
	return &AIHandler{db: db}
}

// 意图定义
type Intent struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Action      string   `json:"action"`
	Confidence  float64  `json:"confidence"`
}

// 数字员工技能
var aiSkills = []Intent{
	{
		Name:        "query_deployment",
		Description: "查询部署实例状态",
		Keywords:    []string{"查看", "查询", "状态", "部署", "实例", "列表", "有哪些"},
		Action:      "list_deployments",
	},
	{
		Name:        "create_deployment",
		Description: "创建新的部署实例",
		Keywords:    []string{"创建", "新建", "部署", "开通", "添加"},
		Action:      "create_deployment",
	},
	{
		Name:        "query_database",
		Description: "查询数据库部署",
		Keywords:    []string{"数据库", "DB", "mysql", "postgresql"},
		Action:      "list_databases",
	},
	{
		Name:        "create_database",
		Description: "创建数据库实例",
		Keywords:    []string{"创建数据库", "新建数据库", "开通数据库"},
		Action:      "create_database",
	},
	{
		Name:        "query_docker",
		Description: "查询 Docker 容器",
		Keywords:    []string{"容器", "docker", "镜像"},
		Action:      "list_docker",
	},
	{
		Name:        "create_docker",
		Description: "创建 Docker 容器",
		Keywords:    []string{"创建容器", "新建容器", "部署容器"},
		Action:      "create_docker",
	},
	{
		Name:        "system_status",
		Description: "查看系统状态",
		Keywords:    []string{"系统", "状态", "健康", "监控", "运行"},
		Action:      "system_status",
	},
	{
		Name:        "greeting",
		Description: "问候",
		Keywords:    []string{"你好", "hi", "hello", "嗨", "在吗"},
		Action:      "greeting",
	},
	{
		Name:        "help",
		Description: "获取帮助",
		Keywords:    []string{"帮助", "help", "怎么用", "功能"},
		Action:      "help",
	},
}

// ChatRequest 聊天请求
type ChatRequest struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`
	Message   string `json:"message" binding:"required"`
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID        string      `json:"id"`
	SessionID string      `json:"session_id"`
	Role      string      `json:"role"`
	Content   string      `json:"content"`
	Intent    string      `json:"intent"`
	Action    string      `json:"action"`
	Data      interface{} `json:"data,omitempty"`
}

// Chat 核心聊天接口
func (h *AIHandler) Chat(c *gin.Context) {
	var req ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 创建或获取会话
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = uuid.New().String()
	}

	userMessage := &models.AIMessage{
		ID:        uuid.New().String(),
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Message,
	}

	// 意图识别
	intent := h.recognizeIntent(req.Message)
	userMessage.Intent = intent.Name
	userMessage.Action = intent.Action

	// 保存用户消息
	models.CreateAIMessage(h.db.GetDB(), userMessage)

	// 生成响应
	response := h.generateResponse(req.Message, intent, req.UserID)

	// 保存助手消息
	assistantMessage := &models.AIMessage{
		ID:        response.ID,
		SessionID: sessionID,
		Role:      "assistant",
		Content:   response.Content,
		Intent:    intent.Name,
		Action:    intent.Action,
		Result:    string(response.Data),
	}
	models.CreateAIMessage(h.db.GetDB(), assistantMessage)

	// 更新会话
	h.updateSession(sessionID, req.UserID, req.Message, response.Content)

	response.SessionID = sessionID
	c.JSON(http.StatusOK, response)
}

// recognizeIntent 意图识别
func (h *AIHandler) recognizeIntent(message string) Intent {
	message = strings.ToLower(message)
	bestMatch := Intent{Name: "unknown", Action: "unknown", Confidence: 0}

	for _, skill := range aiSkills {
		score := 0
		for _, keyword := range skill.Keywords {
			if strings.Contains(message, keyword) {
				score++
			}
		}
		// 计算置信度
		confidence := float64(score) / float64(len(skill.Keywords))
		if confidence > bestMatch.Confidence && score > 0 {
			bestMatch = skill
			bestMatch.Confidence = confidence
		}
	}

	// 如果没有匹配，检查是否包含关键动作词
	if bestMatch.Confidence == 0 {
		if strings.Contains(message, "停") || strings.Contains(message, "删除") || strings.Contains(message, "销毁") {
			bestMatch = Intent{Name: "delete", Action: "delete", Confidence: 0.5}
		} else if strings.Contains(message, "启动") || strings.Contains(message, "开启") {
			bestMatch = Intent{Name: "start", Action: "start", Confidence: 0.5}
		} else if strings.Contains(message, "停止") {
			bestMatch = Intent{Name: "stop", Action: "stop", Confidence: 0.5}
		} else if strings.Contains(message, "监控") || strings.Contains(message, "指标") {
			bestMatch = Intent{Name: "monitor", Action: "monitor", Confidence: 0.5}
		}
	}

	return bestMatch
}

// generateResponse 生成响应
func (h *AIHandler) generateResponse(message string, intent Intent, userID string) ChatResponse {
	response := ChatResponse{
		ID:    uuid.New().String(),
		Role:  "assistant",
		Intent: intent.Name,
		Action: intent.Action,
	}

	switch intent.Action {
	case "greeting":
		response.Content = fmt.Sprintf("你好！我是 ClawOps 数字员工 🤖\n\n我可以帮你：\n• 查询和管理部署实例\n• 创建数据库（MySQL/PostgreSQL）\n• 管理 Docker 容器\n• 查看系统状态和监控\n\n有什么可以帮你的吗？")

	case "help":
		response.Content = `📖 **ClawOps 数字员工使用指南**

**常用命令：**
- "查看我的部署实例" - 列出所有部署
- "创建一个 MySQL 数据库" - 新建数据库
- "部署一个 Nginx 容器" - 创建 Docker 容器
- "系统状态怎么样" - 查看监控状态

**技能：**
我能理解自然语言，自动执行相应的云资源操作。`

	case "list_deployments":
		data := h.getDeployments()
		response.Data = data
		response.Content = fmt.Sprintf("📋 **你的部署实例**\n\n共 %d 个部署：\n%s\n\n输入"创建部署"可以开通新实例。", len(data), formatDeployments(data))

	case "list_databases":
		data := h.getDatabases()
		response.Data = data
		response.Content = fmt.Sprintf("🗄️ **数据库实例**\n\n共 %d 个数据库：\n%s\n\n输入"创建数据库"可以开通新数据库。", len(data), formatDatabases(data))

	case "list_docker":
		data := h.getDocker()
		response.Data = data
		response.Content = fmt.Sprintf("🐳 **Docker 容器**\n\n共 %d 个容器：\n%s\n\n输入"创建容器"可以部署新容器。", len(data), formatDocker(data))

	case "system_status":
		data := h.getSystemStatus()
		response.Data = data
		response.Content = fmt.Sprintf("📊 **系统状态**\n\n状态：✅ 健康\n在线时间：99.9%%\n\n**资源概览：**\n- 部署实例：%d\n- 数据库：%d\n- Docker 容器：%d", data["deployments"], data["databases"], data["docker"])

	case "create_deployment":
		// 模拟创建
		response.Content = "🎉 收到！我来帮你创建一个部署实例...\n\n请告诉我：\n1. **套餐类型**：community（免费）/ pro（专业版）/ enterprise（企业版）\n2. **实例名称**：给你的实例起个名字\n3. **域名**（可选）：自定义域名\n\n例如："创建一个 pro 版的部署，叫 my-app""

	case "create_database":
		response.Content = "🗄️ 收到！创建数据库\n\n请告诉我：\n1. **数据库类型**：MySQL / PostgreSQL\n2. **版本**：如 MySQL 8.0\n3. **套餐**：small（1核2G）/ medium（2核4G）/ large（4核8G）\n\n例如："创建一个 MySQL 8.0，中型套餐""

	case "create_docker":
		response.Content = "🐳 收到！部署 Docker 容器\n\n请告诉我：\n1. **镜像**：如 nginx / redis / postgres\n2. **容器名称**：给你的容器起个名字\n3. **套餐**：small / medium / large\n\n例如："部署一个 nginx 容器，叫 web-server""

	case "unknown":
		response.Content = fmt.Sprintf("🤔 我理解了你的意思，但需要更具体一点。\n\n你说的："%s"\n\n我可以帮你管理部署、数据库和容器。请试试：\n• "查看部署实例"\n• "创建一个 MySQL 数据库"\n• "部署 Nginx 容器"", message)
	}

	return response
}

// 获取部署数据（模拟）
func (h *AIHandler) getDeployments() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": uuid.New().String(), "name": "prod-api", "status": "running", "plan": "pro", "cpu": 2, "memory": "4GB"},
		{"id": uuid.New().String(), "name": "test-web", "status": "running", "plan": "community", "cpu": 1, "memory": "1GB"},
	}
}

// 获取数据库数据（模拟）
func (h *AIHandler) getDatabases() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": uuid.New().String(), "name": "mysql-prod", "type": "MySQL 8.0", "status": "running", "memory": "4GB"},
		{"id": uuid.New().String(), "name": "pg-main", "type": "PostgreSQL 14", "status": "running", "memory": "8GB"},
	}
}

// 获取 Docker 数据（模拟）
func (h *AIHandler) getDocker() []map[string]interface{} {
	return []map[string]interface{}{
		{"id": uuid.New().String(), "name": "nginx-web", "image": "nginx:latest", "status": "running", "port": 30000},
		{"id": uuid.New().String(), "name": "redis-cache", "image": "redis:7", "status": "running", "port": 30001},
	}
}

// 获取系统状态
func (h *AIHandler) getSystemStatus() map[string]interface{} {
	return map[string]interface{}{
		"status":       "healthy",
		"uptime":       "99.9%",
		"deployments":  2,
		"databases":    2,
		"docker":       2,
		"cpu_usage":    45.5,
		"memory_usage": 38.2,
	}
}

// 更新会话
func (h *AIHandler) updateSession(sessionID, userID, lastUserMsg, lastAssistantMsg string) {
	session := &models.AISession{
		ID:          sessionID,
		UserID:      userID,
		Title:       truncate(lastUserMsg, 50),
		Status:      "active",
		LastMessage: lastAssistantMsg,
	}
	models.CreateAISession(h.db.GetDB(), session)
}

// GetSessions 获取会话列表
func (h *AIHandler) GetSessions(c *gin.Context) {
	userID := c.DefaultQuery("user_id", "default")
	sessions, _ := models.GetAISessions(h.db.GetDB(), userID)
	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// GetMessages 获取历史消息
func (h *AIHandler) GetMessages(c *gin.Context) {
	sessionID := c.Param("session_id")
	messages, _ := models.GetAIMessages(h.db.GetDB(), sessionID, 50)

	// 反转顺序，按时间正序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	c.JSON(http.StatusOK, gin.H{"data": messages})
}

// ExecuteAction 执行动作
func (h *AIHandler) ExecuteAction(c *gin.Context) {
	var req struct {
		SessionID string      `json:"session_id"`
		Action    string      `json:"action" binding:"required"`
		Params    interface{} `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := map[string]interface{}{
		"action":   req.Action,
		"status":   "success",
		"message":  "操作执行成功",
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	switch req.Action {
	case "create_deployment":
		result["message"] = "部署实例创建成功！实例ID: " + uuid.New().String()[:8]
	case "create_database":
		result["message"] = "数据库创建成功！预计 3-5 分钟可用"
	case "create_docker":
		result["message"] = "容器部署成功！"
	case "start":
		result["message"] = "启动成功"
	case "stop":
		result["message"] = "停止成功"
	case "delete":
		result["message"] = "删除成功"
	}

	c.JSON(http.StatusOK, result)
}

// 格式化函数
func formatDeployments(deployments []map[string]interface{}) string {
	var sb strings.Builder
	for i, d := range deployments {
		status := "🟢"
		if d["status"] != "running" {
			status = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%d. %s %s (%s)\n", i+1, status, d["name"], d["plan"]))
	}
	return sb.String()
}

func formatDatabases(databases []map[string]interface{}) string {
	var sb strings.Builder
	for i, d := range databases {
		status := "🟢"
		if d["status"] != "running" {
			status = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%d. %s %s - %s\n", i+1, status, d["name"], d["type"]))
	}
	return sb.String()
}

func formatDocker(docker []map[string]interface{}) string {
	var sb strings.Builder
	for i, d := range docker {
		status := "🟢"
		if d["status"] != "running" {
			status = "🔴"
		}
		sb.WriteString(fmt.Sprintf("%d. %s %s : %s\n", i+1, status, d["name"], d["image"]))
	}
	return sb.String()
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen]) + "..."
	}
	return s
}

// 辅助：解析 JSON
func parseJSON(s string) map[string]interface{} {
	var result map[string]interface{}
	json.Unmarshal([]byte(s), &result)
	return result
}
