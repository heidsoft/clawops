package models

import (
	"time"

	"gorm.io/gorm"
)

// AI 对话消息
type AIMessage struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	SessionID  string    `gorm:"type:varchar(36);index" json:"session_id"`
	Role       string    `gorm:"type:varchar(20)" json:"role"` // user/assistant/system
	Content    string    `gorm:"type:text" json:"content"`
	Intent     string    `gorm:"type:varchar(50)" json:"intent"` // 意图分类
	Action     string    `gorm:"type:varchar(100)" json:"action"` // 执行的动作
	Result     string    `gorm:"type:text" json:"result"` // 执行结果
	Metadata   string    `gorm:"type:text" json:"metadata"` // JSON 额外数据
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AIMessage) TableName() string {
	return "ai_messages"
}

func CreateAIMessage(db *gorm.DB, msg *AIMessage) error {
	return db.Create(msg).Error
}

func GetAIMessages(db *gorm.DB, sessionID string, limit int) ([]AIMessage, error) {
	var messages []AIMessage
	err := db.Where("session_id = ?", sessionID).Order("created_at DESC").Limit(limit).Find(&messages).Error
	return messages, err
}

// AI 技能定义
type AISkill struct {
	ID          string   `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Name        string   `gorm:"type:varchar(100)" json:"name"`
	Description string   `gorm:"type:text" json:"description"`
	Keywords    string   `gorm:"type:text"` // 逗号分隔的关键词
	ActionType  string   `gorm:"type:varchar(50)" json:"action_type"`
	Parameters  string   `gorm:"type:text"` // JSON 格式参数定义
	Enabled     bool     `gorm:"type:boolean" json:"enabled"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (AISkill) TableName() string {
	return "ai_skills"
}

// AI 会话
type AISession struct {
	ID         string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID     string    `gorm:"type:varchar(36);index" json:"user_id"`
	Title      string    `gorm:"type:varchar(200)" json:"title"`
	Status     string    `gorm:"type:varchar(20)" json:"status"` // active/archived
	LastMessage string   `gorm:"type:text" json:"last_message"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (AISession) TableName() string {
	return "ai_sessions"
}

func CreateAISession(db *gorm.DB, session *AISession) error {
	return db.Create(session).Error
}

func GetAISessions(db *gorm.DB, userID string) ([]AISession, error) {
	var sessions []AISession
	err := db.Where("user_id = ? AND status = ?", userID, "active").Order("updated_at DESC").Find(&sessions).Error
	return sessions, err
}

func UpdateAISession(db *gorm.DB, session *AISession) error {
	return db.Save(session).Error
}
