package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// LLM 配置
type LLMConfig struct {
	Provider    string // "openai" / "ollama"
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float64
}

// 默认配置
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

// LLM 消息格式
type LLMMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatMessage 前端传入的消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// LLMResponse LLM 响应
type LLMResponse struct {
	Content   string      `json:"content"`
	Reasoning string      `json:"reasoning,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
}

type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Args     string `json:"arguments"` // JSON string
}

// Skill 定义（用于工具调用）
type Skill struct {
	Name        string
	Description string
	Parameters  string // JSON Schema
}

var availableSkills = []Skill{
	{
		Name:        "get_deployments",
		Description: "获取用户的部署实例列表，返回部署名称、状态、套餐等信息",
		Parameters:  `{"type": "object", "properties": {}}`,
	},
	{
		Name:        "get_databases",
		Description: "获取用户的数据库实例列表，包括 MySQL 和 PostgreSQL",
		Parameters:  `{"type": "object", "properties": {}}`,
	},
	{
		Name:        "get_docker_containers",
		Description: "获取用户的 Docker 容器列表",
		Parameters:  `{"type": "object", "properties": {}}`,
	},
	{
		Name:        "get_system_status",
		Description: "获取系统整体状态，包括 CPU、内存使用率等",
		Parameters:  `{"type": "object", "properties": {}}`,
	},
	{
		Name:        "create_deployment",
		Description: "创建新的部署实例，需要用户提供套餐和名称",
		Parameters:  `{"type": "object", "properties": {"plan": {"type": "string"}, "name": {"type": "string"}, "domain": {"type": "string"}}, "required": ["plan", "name"]}`,
	},
	{
		Name:        "create_database",
		Description: "创建新的数据库实例，支持 MySQL 和 PostgreSQL",
		Parameters:  `{"type": "object", "properties": {"db_type": {"type": "string"}, "version": {"type": "string"}, "plan": {"type": "string"}, "name": {"type": "string"}}, "required": ["db_type", "name"]}`,
	},
	{
		Name:        "create_docker_container",
		Description: "创建新的 Docker 容器，从常用镜像中选择",
		Parameters:  `{"type": "object", "properties": {"image": {"type": "string"}, "name": {"type": "string"}, "plan": {"type": "string"}, "port": {"type": "integer"}}, "required": ["image", "name"]}`,
	},
}

// systemPrompt 系统提示词
func getSystemPrompt() string {
	skillsJSON, _ := json.MarshalIndent(availableSkills, "", "  ")
	
	return fmt.Sprintf(`你是一个专业的云运维助手，名叫 ClawOps 数字员工。

## 你的能力
- 查询和管理云资源（部署实例、数据库、容器）
- 用中文回答，友好且专业
- 支持工具调用来获取数据或执行操作

## 可用工具
%s

## 回复格式
当用户询问信息时，优先调用工具获取真实数据。
当用户要求创建资源时，先询问必要参数，确认后再执行。
如果不确定用户意图，可以反问澄清。

## 回答风格
- 使用 emoji 增加可读性
- 信息结构化展示
- 重要数据用粗体标注`, string(skillsJSON))
}

// CallLLM 调用大模型
func CallLLM(messages []ChatMessage, config LLMConfig) (*LLMResponse, error) {
	// 构建请求
	llmMessages := []LLMMessage{
		{Role: "system", Content: getSystemPrompt()},
	}
	
	for _, msg := range messages {
		role := msg.Role
		if role == "user" {
			role = "user"
		} else if role == "assistant" {
			role = "assistant"
		}
		llmMessages = append(llmMessages, LLMMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	var reqBody map[string]interface{}
	
	if config.Provider == "ollama" {
		// Ollama 格式
		reqBody = map[string]interface{}{
			"model":      config.Model,
			"messages":   llmMessages,
			"stream":     false,
		}
	} else {
		// OpenAI 格式
		reqBody = map[string]interface{}{
			"model":       config.Model,
			"messages":    llmMessages,
			"temperature": config.Temperature,
		}
	}

	reqJSON, _ := json.Marshal(reqBody)
	
	endpoint := config.BaseURL
	if config.Provider == "ollama" {
		endpoint += "/api/chat"
	} else {
		endpoint += "/chat/completions"
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if config.APIKey != "" && config.Provider != "ollama" {
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API error: %s", string(body))
	}

	// 解析响应
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	response := &LLMResponse{}
	
	if config.Provider == "ollama" {
		// Ollama 格式
		if message, ok := result["message"].(map[string]interface{}); ok {
			response.Content = message["content"].(string)
		}
		response.FinishReason = "stop"
	} else {
		// OpenAI 格式
		choices, ok := result["choices"].([]interface{})
		if ok && len(choices) > 0 {
			choice := choices[0].(map[string]interface{})
			if msg, ok := choice["message"].(map[string]interface{}); ok {
				response.Content = msg["content"].(string)
			}
			response.FinishReason = choice["finish_reason"].(string)
		}
	}

	return response, nil
}

// ExecuteTool 执行工具
func ExecuteTool(name string, args map[string]interface{}) (string, error) {
	switch name {
	case "get_deployments":
		return executeGetDeployments()
	case "get_databases":
		return executeGetDatabases()
	case "get_docker_containers":
		return executeGetDockerContainers()
	case "get_system_status":
		return executeGetSystemStatus()
	case "create_deployment":
		return executeCreateDeployment(args)
	case "create_database":
		return executeCreateDatabase(args)
	case "create_docker_container":
		return executeCreateDocker(args)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func executeGetDeployments() (string, error) {
	data := []map[string]interface{}{
		{"名称": "prod-api", "套餐": "pro", "状态": "运行中", "CPU": "2核", "内存": "4GB"},
		{"名称": "test-web", "套餐": "community", "状态": "运行中", "CPU": "1核", "内存": "1GB"},
		{"名称": "staging-db", "套餐": "pro", "状态": "已停止", "CPU": "2核", "内存": "4GB"},
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeGetDatabases() (string, error) {
	data := []map[string]interface{}{
		{"名称": "mysql-prod", "类型": "MySQL 8.0", "状态": "运行中", "内存": "4GB", "磁盘": "100GB"},
		{"名称": "pg-main", "类型": "PostgreSQL 14", "状态": "运行中", "内存": "8GB", "磁盘": "200GB"},
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeGetDockerContainers() (string, error) {
	data := []map[string]interface{}{
		{"名称": "nginx-web", "镜像": "nginx:latest", "状态": "运行中", "端口": "30000"},
		{"名称": "redis-cache", "镜像": "redis:7", "状态": "运行中", "端口": "30001"},
		{"名称": "grafana-monitor", "镜像": "grafana/grafana", "状态": "已停止", "端口": "30002"},
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeGetSystemStatus() (string, error) {
	data := map[string]interface{}{
		"状态":       "健康",
		"在线时间":     "99.9%",
		"部署实例":     3,
		"数据库":      2,
		"Docker容器":  3,
		"CPU使用率":   "45%",
		"内存使用率":   "38%",
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeCreateDeployment(args map[string]interface{}) (string, error) {
	name, _ := args["name"].(string)
	plan, _ := args["plan"].(string)
	if plan == "" {
		plan = "community"
	}
	if name == "" {
		return "", fmt.Errorf("缺少参数: name")
	}
	data := map[string]interface{}{
		"message": fmt.Sprintf("✅ 部署实例「%s」创建成功！\n套餐: %s\n预计 3-5 分钟可用", name, plan),
		"instance_id": "ins-" + strings.ToLower(name)[:8],
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeCreateDatabase(args map[string]interface{}) (string, error) {
	name, _ := args["name"].(string)
	dbType, _ := args["db_type"].(string)
	if dbType == "" {
		dbType = "MySQL"
	}
	if name == "" {
		return "", fmt.Errorf("缺少参数: name")
	}
	data := map[string]interface{}{
		"message": fmt.Sprintf("✅ 数据库「%s」( %s ) 创建成功！\n预计 3-5 分钟初始化完成", name, dbType),
		"instance_id": "db-" + strings.ToLower(name)[:8],
	}
	return json.MarshalIndent(data, "", "  ")
}

func executeCreateDocker(args map[string]interface{}) (string, error) {
	name, _ := args["name"].(string)
	image, _ := args["image"].(string)
	if name == "" || image == "" {
		return "", fmt.Errorf("缺少参数: name 或 image")
	}
	data := map[string]interface{}{
		"message": fmt.Sprintf("✅ 容器「%s」( %s ) 部署成功！", name, image),
		"container_id": "cnt-" + strings.ToLower(name)[:8],
	}
	return json.MarshalIndent(data, "", "  ")
}
