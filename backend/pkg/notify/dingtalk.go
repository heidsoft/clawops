package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// DingTalk 告警通知

type DingTalkClient struct {
	Webhook string      // Webhook 地址
	Secret  string      // 加签密钥（可选）
	Client  *http.Client
}

type DingTalkMessage struct {
	MsgType string      `json:"msgtype"`
	Text    *TextContent `json:"text,omitempty"`
	Markdown *MarkdownContent `json:"markdown,omitempty"`
	At      *AtConfig    `json:"at,omitempty"`
}

type TextContent struct {
	Content string `json:"content"`
}

type MarkdownContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type AtConfig struct {
	AtMobiles []string `json:"atMobiles,omitempty"`
	IsAtAll  bool     `json:"isAtAll,omitempty"`
}

// 创建钉钉客户端
func NewDingTalkClient(webhook, secret string) *DingTalkClient {
	return &DingTalkClient{
		Webhook: webhook,
		Secret:  secret,
		Client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// 发送文本消息
func (c *DingTalkClient) SendText(content string, atMobiles []string, isAtAll bool) error {
	if c.Webhook == "" {
		return fmt.Errorf("webhook is empty")
	}

	msg := DingTalkMessage{
		MsgType: "text",
		Text:    &TextContent{Content: content},
	}

	if len(atMobiles) > 0 || isAtAll {
		msg.At = &AtConfig{
			AtMobiles: atMobiles,
			IsAtAll:   isAtAll,
		}
	}

	return c.send(msg)
}

// 发送 Markdown 消息（支持标题和内容）
func (c *DingTalkClient) SendMarkdown(title, content string) error {
	if c.Webhook == "" {
		return fmt.Errorf("webhook is empty")
	}

	msg := DingTalkMessage{
		MsgType: "markdown",
		Markdown: &MarkdownContent{
			Title: title,
			Text:  content,
		},
	}

	return c.send(msg)
}

// 发送告警卡片（自定义样式）
func (c *DingTalkClient) SendAlertCard(alert AlertInfo) error {
	if c.Webhook == "" {
		return fmt.Errorf("webhook is empty")
	}

	// 根据严重程度选择颜色
	color := "red" // critical
	if alert.Severity == "warning" {
		color = "yellow"
	} else if alert.Severity == "info" {
		color = "green"
	}

	content := fmt.Sprintf(`## 🚨 %s

**级别:** %s
**实例:** %s
**时间:** %s

---

### 📊 指标数据

| 指标 | 当前值 | 阈值 |
|------|--------|------|
| %s | %.1f%% | %.1f%% |

---

%s

> 来源: ClawOps 监控告警`, 
		alert.RuleName,
		alert.SeverityText(),
		alert.InstanceName,
		alert.TriggeredAt,
		alert.MetricType,
		alert.CurrentValue,
		alert.Threshold,
		alert.Message,
	)

	return c.SendMarkdown(alert.RuleName, content)
}

// 发送自定义链接卡片
func (c *DingTalkClient) SendLinkCard(title, text, messageURL, picURL string) error {
	if c.Webhook == "" {
		return fmt.Errorf("webhook is empty")
	}

	msg := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      title,
			"text":       text,
			"messageUrl": messageURL,
			"picUrl":     picURL,
		},
	}

	jsonData, _ := json.Marshal(msg)
	resp, err := c.Client.Post(c.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// 发送消息
func (c *DingTalkClient) send(msg DingTalkMessage) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := c.Client.Post(c.Webhook, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.checkResponse(resp)
}

// 检查响应
func (c *DingTalkClient) checkResponse(resp *http.Response) error {
	if resp.StatusCode != 200 {
		return fmt.Errorf(" DingTalk API error: status code %d", resp.StatusCode)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if err, ok := result["errcode"].(float64); ok && err != 0 {
		return fmt.Errorf(" DingTalk API error: %v", result["errmsg"])
	}

	return nil
}

// 告警信息结构
type AlertInfo struct {
	RuleName     string
	Severity     string  // critical / warning / info
	InstanceName string
	InstanceID   string
	MetricType   string
	CurrentValue float64
	Threshold    float64
	Message      string
	TriggeredAt  time.Time
}

// 获取严重程度文字
func (a *AlertInfo) SeverityText() string {
	switch a.Severity {
	case "critical":
		return "🔴 严重"
	case "warning":
		return "🟡 警告"
	case "info":
		return "🔵 通知"
	default:
		return a.Severity
	}
}

// 快捷发送函数
func SendDingTalkAlert(webhook string, alert AlertInfo) error {
	client := NewDingTalkClient(webhook, "")
	return client.SendAlertCard(alert)
}

// 发送简单文本告警
func SendDingTalkText(webhook, content string) error {
	client := NewDingTalkClient(webhook, "")
	return client.SendText(content, nil, false)
}
