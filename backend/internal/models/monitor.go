package models

import (
	"time"

	"gorm.io/gorm"
)

// 监控指标数据
type MetricData struct {
	ID          string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	InstanceID  string    `gorm:"type:varchar(36);index" json:"instance_id"` // 关联实例
	MetricType  string    `gorm:"type:varchar(50)" json:"metric_type"`      // cpu/memory/disk/network
	MetricName  string    `gorm:"type:varchar(100)" json:"metric_name"`     // 具体指标名
	Value       float64   `json:"value"`                                     // 指标值
	Unit        string    `gorm:"type:varchar(20)" json:"unit"`             // %/GB/MBps
	Timestamp   time.Time `gorm:"index" json:"timestamp"`                   // 采集时间
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (MetricData) TableName() string {
	return "metric_data"
}

// 告警规则
type AlertRule struct {
	ID          string  `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string  `gorm:"type:varchar(100)" json:"name"`           // 规则名称
	Description string  `gorm:"type:text" json:"description"`           // 规则描述
	MetricType  string  `gorm:"type:varchar(50)" json:"metric_type"`    // cpu/memory/disk/network
	Operator    string  `gorm:"type:varchar(10)" json:"operator"`       // >/</>=/<=
	Threshold   float64 `json:"threshold"`                               // 阈值
	Duration    int     `json:"duration"`                                // 持续时间(秒)，超过则告警
	Severity    string  `gorm:"type:varchar(20)" json:"severity"`       // critical/warning/info
	Enabled     bool    `gorm:"type:boolean;default:true" json:"enabled"`
	NotifyTypes string  `gorm:"type:varchar(100)" json:"notify_types"`   // dingtalk/email/webhook,逗号分隔
	InstanceID  string  `gorm:"type:varchar(36)" json:"instance_id"`     // 空表示全局规则
	UserID      string  `gorm:"type:varchar(36)" json:"user_id"`         // 拥有者
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AlertRule) TableName() string {
	return "alert_rules"
}

// 告警记录
type Alert struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	RuleID     string    `gorm:"type:varchar(36);index" json:"rule_id"`
	RuleName   string    `gorm:"type:varchar(100)" json:"rule_name"`
	InstanceID string    `gorm:"type:varchar(36);index" json:"instance_id"`
	InstanceName string  `gorm:"type:varchar(100)" json:"instance_name"`
	MetricType string    `gorm:"type:varchar(50)" json:"metric_type"`
	Severity   string    `gorm:"type:varchar(20)" json:"severity"`
	Message    string    `gorm:"type:text" json:"message"`
	Status     string    `gorm:"type:varchar(20);index" json:"status"` // firing/resolved/acknowledged
	TriggeredAt time.Time `gorm:"index" json:"triggered_at"`
	ResolvedAt  *time.Time `json:"resolved_at"`
	AcknowledgedAt *time.Time `json:"acknowledged_at"`
	AcknowledgedBy string `gorm:"type:varchar(100)" json:"acknowledged_by"`
	UserID     string    `gorm:"type:varchar(36)" json:"user_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Alert) TableName() string {
	return "alerts"
}

// 创建默认告警规则
func CreateDefaultAlertRules(db *gorm.DB, userID string) error {
	defaultRules := []AlertRule{
		{
			ID:          "rule-cpu-high",
			Name:        "CPU 使用率过高",
			Description: "当 CPU 使用率超过 80% 时告警",
			MetricType:  "cpu",
			Operator:    ">",
			Threshold:   80,
			Duration:    60,
			Severity:    "warning",
			NotifyTypes: "dingtalk",
			UserID:      userID,
		},
		{
			ID:          "rule-cpu-critical",
			Name:        "CPU 使用率严重",
			Description: "当 CPU 使用率超过 95% 时告警",
			MetricType:  "cpu",
			Operator:    ">",
			Threshold:   95,
			Duration:    30,
			Severity:    "critical",
			NotifyTypes: "dingtalk",
			UserID:      userID,
		},
		{
			ID:          "rule-memory-high",
			Name:        "内存使用率过高",
			Description: "当内存使用率超过 85% 时告警",
			MetricType:  "memory",
			Operator:    ">",
			Threshold:   85,
			Duration:    120,
			Severity:    "warning",
			NotifyTypes: "dingtalk",
			UserID:      userID,
		},
		{
			ID:          "rule-disk-high",
			Name:        "磁盘使用率过高",
			Description: "当磁盘使用率超过 90% 时告警",
			MetricType:  "disk",
			Operator:    ">",
			Threshold:   90,
			Duration:    300,
			Severity:    "warning",
			NotifyTypes: "dingtalk",
			UserID:      userID,
		},
	}

	for _, rule := range defaultRules {
		// 检查是否已存在
		var existing AlertRule
		if err := db.Where("id = ? AND user_id = ?", rule.ID, userID).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&rule).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// 获取最近的指标数据
func GetRecentMetrics(db *gorm.DB, instanceID string, metricType string, minutes int) ([]MetricData, error) {
	var metrics []MetricData
	since := time.Now().Add(-time.Duration(minutes) * time.Minute)
	err := db.Where("instance_id = ? AND metric_type = ? AND timestamp > ?", 
		instanceID, metricType, since).
		Order("timestamp ASC").
		Find(&metrics).Error
	return metrics, err
}

// 获取活跃告警
func GetActiveAlerts(db *gorm.DB, userID string) ([]Alert, error) {
	var alerts []Alert
	err := db.Where("user_id = ? AND status IN ?", userID, []string{"firing", "acknowledged"}).
		Order("triggered_at DESC").
		Find(&alerts).Error
	return alerts, err
}

// 获取最近的告警历史
func GetAlertHistory(db *gorm.DB, userID string, limit int) ([]Alert, error) {
	var alerts []Alert
	err := db.Where("user_id = ?", userID).
		Order("triggered_at DESC").
		Limit(limit).
		Find(&alerts).Error
	return alerts, err
}
