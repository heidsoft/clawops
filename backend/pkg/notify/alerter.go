package notify

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

// 告警规则
type AlertRule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	MetricType  string   `json:"metric_type"`  // cpu/memory/disk/network
	Operator    string   `json:"operator"`     // >/</>=/<=
	Threshold   float64  `json:"threshold"`
	Duration    int      `json:"duration"`     // 持续秒数
	Severity    string   `json:"severity"`     // critical/warning/info
	NotifyTypes []string `json:"notify_types"` // dingtalk/email/webhook
	Enabled     bool     `json:"enabled"`
}

// 告警实例
type Alert struct {
	ID           string    `json:"id"`
	RuleID       string    `json:"rule_id"`
	RuleName     string    `json:"rule_name"`
	InstanceID   string    `json:"instance_id"`
	InstanceName string    `json:"instance_name"`
	Severity     string    `json:"severity"`
	Message      string    `json:"message"`
	Status       string    `json:"status"` // firing/acknowledged/resolved
	TriggeredAt  time.Time `json:"triggered_at"`
	ResolvedAt   *time.Time `json:"resolved_at,omitempty"`
}

// 告警通知器
type Alerter struct {
	rules     []AlertRule
	alerts    map[string]*Alert // active alerts by rule+instance key
	dingTalk  *DingTalkClient
}

// 创建告警通知器
func NewAlerter(dingtalkWebhook string) *Alerter {
	return &Alerter{
		rules:    make([]AlertRule, 0),
		alerts:   make(map[string]*Alert),
		dingTalk: NewDingTalkClient(dingtalkWebhook, ""),
	}
}

// 设置告警规则
func (a *Alerter) SetRules(rules []AlertRule) {
	a.rules = rules
}

// 添加规则
func (a *Alerter) AddRule(rule AlertRule) {
	a.rules = append(a.rules, rule)
}

// 检查指标是否触发告警
func (a *Alerter) CheckMetric(instanceID, instanceName, metricType string, value float64) *Alert {
	for i := range a.rules {
		rule := &a.rules[i]
		if !rule.Enabled {
			continue
		}

		// 匹配指标类型
		if rule.MetricType != metricType {
			continue
		}

		// 检查阈值
		triggered := false
		switch rule.Operator {
		case ">":
			triggered = value > rule.Threshold
		case "<":
			triggered = value < rule.Threshold
		case ">=":
			triggered = value >= rule.Threshold
		case "<=":
			triggered = value <= rule.Threshold
		}

		if triggered {
			key := fmt.Sprintf("%s:%s", rule.ID, instanceID)
			
			// 检查是否已有活跃告警
			if existing, ok := a.alerts[key]; ok {
				return nil // 已有活跃告警，不重复触发
			}

			// 创建新告警
			alert := &Alert{
				ID:           uuid.New().String(),
				RuleID:       rule.ID,
				RuleName:     rule.Name,
				InstanceID:   instanceID,
				InstanceName: instanceName,
				Severity:     rule.Severity,
				Message:      fmt.Sprintf("%s 超过阈值 (当前: %.1f%%, 阈值: %.1f%%)", rule.MetricType, value, rule.Threshold),
				Status:       "firing",
				TriggeredAt:  time.Now(),
			}

			a.alerts[key] = alert
			return alert
		}
	}
	return nil
}

// 触发告警通知
func (a *Alerter) FireAlert(alert *Alert) error {
	if alert == nil {
		return nil
	}

	// 查找对应的规则
	var rule *AlertRule
	for i := range a.rules {
		if a.rules[i].ID == alert.RuleID {
			rule = &a.rules[i]
			break
		}
	}

	if rule == nil {
		return fmt.Errorf("rule not found")
	}

	// 构建通知信息
	info := AlertInfo{
		RuleName:     alert.RuleName,
		Severity:     alert.Severity,
		InstanceName: alert.InstanceName,
		InstanceID:   alert.InstanceID,
		MetricType:   rule.MetricType,
		CurrentValue: rule.Threshold + 10, // 模拟值
		Threshold:    rule.Threshold,
		Message:      alert.Message,
		TriggeredAt:  alert.TriggeredAt,
	}

	// 发送通知
	for _, notifyType := range rule.NotifyTypes {
		switch notifyType {
		case "dingtalk":
			if err := a.dingTalk.SendAlertCard(info); err != nil {
				return fmt.Errorf("failed to send dingtalk notification: %w", err)
			}
		case "email":
			// TODO: 实现邮件通知
			fmt.Println("[Alerter] Email notification not implemented yet")
		case "webhook":
			// TODO: 实现 webhook 通知
			fmt.Println("[Alerter] Webhook notification not implemented yet")
		}
	}

	return nil
}

// 确认告警
func (a *Alerter) AcknowledgeAlert(instanceID, ruleID string) error {
	key := fmt.Sprintf("%s:%s", ruleID, instanceID)
	if alert, ok := a.alerts[key]; ok {
		alert.Status = "acknowledged"
	}
	return nil
}

// 解决告警
func (a *Alerter) ResolveAlert(instanceID, ruleID string) error {
	key := fmt.Sprintf("%s:%s", ruleID, instanceID)
	if alert, ok := a.alerts[key]; ok {
		alert.Status = "resolved"
		now := time.Now()
		alert.ResolvedAt = &now
		delete(a.alerts, key)
	}
	return nil
}

// 获取活跃告警列表
func (a *Alerter) GetActiveAlerts() []*Alert {
	alerts := make([]*Alert, 0, len(a.alerts))
	for _, alert := range a.alerts {
		if alert.Status == "firing" || alert.Status == "acknowledged" {
			alerts = append(alerts, alert)
		}
	}
	return alerts
}

// 模拟告警检测（用于测试）
func (a *Alerter) SimulateDetection(instances []struct{ ID, Name string }) {
	metrics := []string{"cpu", "memory", "disk", "network"}

	for _, instance := range instances {
		for _, metric := range metrics {
			// 随机生成指标值
			value := rand.Float64() * 100

			// 检查是否触发告警
			alert := a.CheckMetric(instance.ID, instance.Name, metric, value)
			if alert != nil {
				fmt.Printf("[Alert] 🚨 %s - %s (%.1f%%)\n", instance.Name, metric, value)
				a.FireAlert(alert)
			}
		}
	}
}

// 初始化默认规则
func GetDefaultRules() []AlertRule {
	return []AlertRule{
		{
			ID:          "rule-cpu-high",
			Name:        "CPU 使用率过高",
			MetricType:  "cpu",
			Operator:    ">",
			Threshold:   80,
			Duration:    60,
			Severity:    "warning",
			NotifyTypes: []string{"dingtalk"},
			Enabled:     true,
		},
		{
			ID:          "rule-cpu-critical",
			Name:        "CPU 使用率严重",
			MetricType:  "cpu",
			Operator:    ">",
			Threshold:   95,
			Duration:    30,
			Severity:    "critical",
			NotifyTypes: []string{"dingtalk"},
			Enabled:     true,
		},
		{
			ID:          "rule-memory-high",
			Name:        "内存使用率过高",
			MetricType:  "memory",
			Operator:    ">",
			Threshold:   85,
			Duration:    120,
			Severity:    "warning",
			NotifyTypes: []string{"dingtalk"},
			Enabled:     true,
		},
		{
			ID:          "rule-disk-high",
			Name:        "磁盘使用率过高",
			MetricType:  "disk",
			Operator:    ">",
			Threshold:   90,
			Duration:    300,
			Severity:    "warning",
			NotifyTypes: []string{"dingtalk"},
			Enabled:     true,
		},
	}
}
