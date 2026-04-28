package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"openclaw-deploy/internal/models"
)

// 指标数据结构
type MetricResponse struct {
	InstanceID   string                 `json:"instance_id"`
	InstanceName string                 `json:"instance_name"`
	Metrics      map[string]interface{} `json:"metrics"`
	Timestamp    time.Time              `json:"timestamp"`
}

// 获取实例的监控数据
func GetInstanceMetrics(c *gin.Context, db *gorm.DB) {
	instanceID := c.Query("instance_id")
	metricType := c.Query("type") // cpu/memory/disk/network
	minutes := 60                  // 默认最近60分钟

	if instanceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "instance_id is required"})
		return
	}

	metrics := make(map[string]interface{})

	// 如果没有指定类型，返回所有类型
	if metricType == "" || metricType == "cpu" {
		metrics["cpu"] = generateMockCPUData(minutes)
	}
	if metricType == "" || metricType == "memory" {
		metrics["memory"] = generateMockMemoryData(minutes)
	}
	if metricType == "" || metricType == "disk" {
		metrics["disk"] = generateMockDiskData()
	}
	if metricType == "" || metricType == "network" {
		metrics["network"] = generateMockNetworkData(minutes)
	}

	// 获取实例名称
	instanceName := getInstanceName(db, instanceID)

	c.JSON(http.StatusOK, gin.H{
		"instance_id":   instanceID,
		"instance_name": instanceName,
		"metrics":       metrics,
		"timestamp":     time.Now(),
	})
}

// 生成模拟 CPU 数据
func generateMockCPUData(minutes int) map[string]interface{} {
	labels := []string{}
	values := []float64{}
	now := time.Now()

	for i := minutes; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Minute)
		labels = append(labels, t.Format("15:04"))
		// 生成 30-70% 之间的随机值
		values = append(values, 30+rand.Float64()*40)
	}

	return map[string]interface{}{
		"usage":   values[len(values)-1],
		"labels":  labels,
		"values":  values,
		"average": average(values),
		"max":     max(values),
		"min":     min(values),
		"unit":    "%",
	}
}

// 生成模拟内存数据
func generateMockMemoryData(minutes int) map[string]interface{} {
	labels := []string{}
	values := []float64{}
	now := time.Now()

	for i := minutes; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Minute)
		labels = append(labels, t.Format("15:04"))
		// 生成 40-80% 之间的随机值
		values = append(values, 40+rand.Float64()*40)
	}

	return map[string]interface{}{
		"usage":   values[len(values)-1],
		"labels":  labels,
		"values":  values,
		"average": average(values),
		"max":     max(values),
		"min":     min(values),
		"total":   16384, // 16GB
		"used":    values[len(values)-1] / 100 * 16384,
		"unit":    "%",
	}
}

// 生成模拟磁盘数据
func generateMockDiskData() map[string]interface{} {
	return map[string]interface{}{
		"usage":    45.5,
		"total":    500, // GB
		"used":     227.5,
		"free":     272.5,
		"io_read":  125.5,
		"io_write": 85.2,
		"unit":     "%",
	}
}

// 生成模拟网络数据
func generateMockNetworkData(minutes int) map[string]interface{} {
	labels := []string{}
	inValues := []float64{}
	outValues := []float64{}
	now := time.Now()

	for i := minutes; i >= 0; i-- {
		t := now.Add(-time.Duration(i) * time.Minute)
		labels = append(labels, t.Format("15:04"))
		// 生成随机的入站和出站流量
		inValues = append(inValues, 50+rand.Float64()*100)
		outValues = append(outValues, 30+rand.Float64()*80)
	}

	return map[string]interface{}{
		"labels":     labels,
		"in_values":   inValues,
		"out_values": outValues,
		"in_current":  inValues[len(inValues)-1],
		"out_current": outValues[len(outValues)-1],
		"unit":        "Mbps",
	}
}

// 获取所有实例的监控概览
func GetMetricsOverview(c *gin.Context, db *gorm.DB) {
	// 获取所有部署实例
	var deployments []models.Deployment
	db.Find(&deployments)

	overview := make([]gin.H, 0, len(deployments))
	for _, d := range deployments {
		overview = append(overview, gin.H{
			"instance_id":   d.ID,
			"instance_name": d.InstanceName,
			"status":        d.Status,
			"cpu":           30 + rand.Float64()*40,
			"memory":        40 + rand.Float64()*40,
			"disk":          40 + rand.Float64()*20,
			"network_in":    rand.Float64() * 100,
			"network_out":   rand.Float64() * 80,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"overview": overview,
		"timestamp": time.Now(),
	})
}

// 获取告警规则列表
func GetAlertRules(c *gin.Context, db *gorm.DB) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "default"
	}

	var rules []models.AlertRule
	err := db.Where("user_id = ?", userID).Find(&rules).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rules": rules,
		"total": len(rules),
	})
}

// 创建告警规则
func CreateAlertRule(c *gin.Context, db *gorm.DB) {
	var rule models.AlertRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rule.ID = uuid.New().String()
	if rule.UserID == "" {
		rule.UserID = "default"
	}

	if err := db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rule":  rule,
		"message": "规则创建成功",
	})
}

// 更新告警规则
func UpdateAlertRule(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var rule models.AlertRule
	
	if err := db.First(&rule, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "规则不存在"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 不允许更新 ID
	delete(updates, "id")

	if err := db.Model(&rule).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rule":    rule,
		"message": "规则更新成功",
	})
}

// 删除告警规则
func DeleteAlertRule(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	if err := db.Delete(&models.AlertRule{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "规则删除成功",
	})
}

// 获取告警列表
func GetAlerts(c *gin.Context, db *gorm.DB) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "default"
	}

	status := c.Query("status") // firing/resolved/acknowledged

	query := db.Model(&models.Alert{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var alerts []models.Alert
	err := query.Order("triggered_at DESC").Limit(100).Find(&alerts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  len(alerts),
	})
}

// 确认告警
func AcknowledgeAlert(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var alert models.Alert
	
	if err := db.First(&alert, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "告警不存在"})
		return
	}

	now := time.Now()
	alert.Status = "acknowledged"
	alert.AcknowledgedAt = &now
	alert.AcknowledgedBy = c.Query("user") // 简化，实际应从认证获取

	if err := db.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alert":   alert,
		"message": "告警已确认",
	})
}

// 解决告警
func ResolveAlert(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var alert models.Alert
	
	if err := db.First(&alert, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "告警不存在"})
		return
	}

	now := time.Now()
	alert.Status = "resolved"
	alert.ResolvedAt = &now

	if err := db.Save(&alert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alert":   alert,
		"message": "告警已解决",
	})
}

// 获取告警统计
func GetAlertStats(c *gin.Context, db *gorm.DB) {
	userID := c.Query("user_id")
	if userID == "" {
		userID = "default"
	}

	var total int64
	var firing int64
	var acknowledged int64
	var resolved int64

	db.Model(&models.Alert{}).Where("user_id = ?", userID).Count(&total)
	db.Model(&models.Alert{}).Where("user_id = ? AND status = ?", userID, "firing").Count(&firing)
	db.Model(&models.Alert{}).Where("user_id = ? AND status = ?", userID, "acknowledged").Count(&acknowledged)
	db.Model(&models.Alert{}).Where("user_id = ? AND status = ?", userID, "resolved").Count(&resolved)

	// 获取最近7天的告警趋势
	weekAgo := time.Now().AddDate(0, 0, -7)
	var weekAlerts []models.Alert
	db.Where("user_id = ? AND triggered_at > ?", userID, weekAgo).
		Order("triggered_at ASC").
		Find(&weekAlerts)

	// 按天统计
	dailyStats := make(map[string]int)
	for _, alert := range weekAlerts {
		day := alert.TriggeredAt.Format("01-02")
		dailyStats[day]++
	}

	c.JSON(http.StatusOK, gin.H{
		"total":        total,
		"firing":       firing,
		"acknowledged": acknowledged,
		"resolved":     resolved,
		"daily_stats":  dailyStats,
		"timestamp":    time.Now(),
	})
}

// 辅助函数
func getInstanceName(db *gorm.DB, instanceID string) string {
	var deployment models.Deployment
	if err := db.First(&deployment, "id = ?", instanceID).Error; err == nil {
		return deployment.InstanceName
	}
	return instanceID
}

func average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	m := values[0]
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

// 钉钉通知
func SendDingTalkNotification(webhook string, message string) error {
	if webhook == "" {
		return nil
	}

	payload := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}

	jsonData, _ := json.Marshal(payload)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(webhook, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
