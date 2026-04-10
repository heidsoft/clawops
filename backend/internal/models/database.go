package models

import (
	"time"

	"gorm.io/gorm"
)

type DatabaseType string

const (
	MySQL      DatabaseType = "mysql"
	PostgreSQL DatabaseType = "postgresql"
)

type DatabaseDeployment struct {
	ID           string       `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string       `gorm:"type:varchar(36);index" json:"user_id"`
	Name         string       `gorm:"type:varchar(100)" json:"name"`
	DatabaseType DatabaseType `gorm:"type:varchar(20)" json:"database_type"`
	Version      string       `gorm:"type:varchar(20)" json:"version"`
	InstanceID   string       `gorm:"type:varchar(50)" json:"instance_id"`
	Host         string       `gorm:"type:varchar(50)" json:"host"`
	Port         int          `gorm:"type:int" json:"port"`
	Username     string       `gorm:"type:varchar(100)" json:"username"`
	Password     string       `gorm:"type:varchar(200)" json:"-"`
	Status       string       `gorm:"type:varchar(20)" json:"status"`
	DiskSize     int          `gorm:"type:int" json:"disk_size"`
	MemorySize   int          `gorm:"type:int" json:"memory_size"` // MB
	Region       string       `gorm:"type:varchar(50)" json:"region"`
	ConnectionURL string      `gorm:"-" json:"connection_url"`
	CreatedAt    time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time   `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DatabaseDeployment) TableName() string {
	return "database_deployments"
}

func CreateDatabaseDeployment(db *gorm.DB, deployment *DatabaseDeployment) error {
	return db.Create(deployment).Error
}

func GetDatabaseDeployments(db *gorm.DB, page, pageSize int) ([]DatabaseDeployment, int64, error) {
	var deployments []DatabaseDeployment
	var total int64

	offset := (page - 1) * pageSize

	if err := db.Model(&DatabaseDeployment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&deployments).Error
	return deployments, total, err
}

func GetDatabaseDeployment(db *gorm.DB, id string) (*DatabaseDeployment, error) {
	var deployment DatabaseDeployment
	err := db.First(&deployment, "id = ?", id).Error
	return &deployment, err
}

func UpdateDatabaseDeployment(db *gorm.DB, deployment *DatabaseDeployment) error {
	return db.Save(deployment).Error
}

func DeleteDatabaseDeployment(db *gorm.DB, id string) error {
	return db.Delete(&DatabaseDeployment{}, "id = ?", id).Error
}
