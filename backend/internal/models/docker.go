package models

import (
	"time"

	"gorm.io/gorm"
)

type DockerDeployment struct {
	ID           string `gorm:"primaryKey;type:varchar(36)" json:"id"`
	UserID       string `gorm:"type:varchar(36);index" json:"user_id"`
	Name         string `gorm:"type:varchar(100)" json:"name"`
	InstanceID   string `gorm:"type:varchar(50)" json:"instance_id"`
	Image        string `gorm:"type:varchar(200)" json:"image"`
	ContainerID  string `gorm:"type:varchar(100)" json:"container_id"`
	Host         string `gorm:"type:varchar(50)" json:"host"`
	Port         int    `gorm:"type:int" json:"port"`
	ExternalPort int    `gorm:"type:int" json:"external_port"`
	Status       string `gorm:"type:varchar(20)" json:"status"`
	CPU          int    `gorm:"type:int" json:"cpu"`          // cores
	Memory       int    `gorm:"type:int" json:"memory"`       // MB
	DiskSize     int    `gorm:"type:int" json:"disk_size"`    // GB
	Region       string `gorm:"type:varchar(50)" json:"region"`
	Command      string `gorm:"type:varchar(500)" json:"command"`
	Environment  string `gorm:"type:text" json:"environment"` // JSON string
	Volumes      string    `gorm:"type:text" json:"volumes"`  // JSON array string
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DockerDeployment) TableName() string {
	return "docker_deployments"
}

func CreateDockerDeployment(db *gorm.DB, deployment *DockerDeployment) error {
	return db.Create(deployment).Error
}

func GetDockerDeployments(db *gorm.DB, page, pageSize int) ([]DockerDeployment, int64, error) {
	var deployments []DockerDeployment
	var total int64

	offset := (page - 1) * pageSize

	if err := db.Model(&DockerDeployment{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&deployments).Error
	return deployments, total, err
}

func GetDockerDeployment(db *gorm.DB, id string) (*DockerDeployment, error) {
	var deployment DockerDeployment
	err := db.First(&deployment, "id = ?", id).Error
	return &deployment, err
}

func UpdateDockerDeployment(db *gorm.DB, deployment *DockerDeployment) error {
	return db.Save(deployment).Error
}

func DeleteDockerDeployment(db *gorm.DB, id string) error {
	return db.Delete(&DockerDeployment{}, "id = ?", id).Error
}
