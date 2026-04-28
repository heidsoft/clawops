package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"openclaw-deploy/internal/models"
)

// 获取 Skills 列表
func GetSkills(c *gin.Context, db *gorm.DB) {
	category := c.Query("category")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	skills, total, err := models.GetSkills(db, category, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取分类统计
	var categories []models.SkillCategory
	db.Order("sort_order ASC").Find(&categories)

	// 统计每个分类的数量
	for i := range categories {
		var count int64
		db.Model(&models.Skill{}).Where("category = ? AND enabled = ?", categories[i].Name, true).Count(&count)
		categories[i].Count = int(count)
	}

	c.JSON(http.StatusOK, gin.H{
		"skills":     skills,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		"categories": categories,
	})
}

// 获取单个 Skill 详情
func GetSkill(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var skill models.Skill
	if err := db.First(&skill, "id = ? OR name = ?", id, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"skill": skill,
	})
}

// 获取用户已安装的 Skills
func GetUserSkills(c *gin.Context, db *gorm.DB) {
	userID := c.DefaultQuery("user_id", "default")

	userSkills, err := models.GetUserSkills(db, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取完整的 Skill 信息
	type SkillWithDetail struct {
		models.UserSkill
		Skill *models.Skill `json:"skill,omitempty"`
	}

	result := make([]SkillWithDetail, 0, len(userSkills))
	for _, us := range userSkills {
		var skill models.Skill
		db.First(&skill, "id = ?", us.SkillID)
		result = append(result, SkillWithDetail{
			UserSkill: us,
			Skill:     &skill,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"user_skills": result,
		"total":       len(result),
	})
}

// 安装 Skill
func InstallSkill(c *gin.Context, db *gorm.DB) {
	var req struct {
		SkillID string `json:"skill_id" binding:"required"`
		UserID  string `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.UserID == "" {
		req.UserID = "default"
	}

	userSkill, err := models.InstallSkill(db, req.UserID, req.SkillID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_skill": userSkill,
		"message":    "Skill 安装成功",
	})
}

// 卸载 Skill
func UninstallSkill(c *gin.Context, db *gorm.DB) {
	userSkillID := c.Param("id")
	userID := c.DefaultQuery("user_id", "default")

	err := models.UninstallSkill(db, userID, userSkillID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Skill 卸载成功",
	})
}

// 更新用户 Skill 配置
func UpdateUserSkillConfig(c *gin.Context, db *gorm.DB) {
	userSkillID := c.Param("id")
	userID := c.DefaultQuery("user_id", "default")

	var req struct {
		Config  string `json:"config"`
		Enabled bool   `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userSkill models.UserSkill
	if err := db.Where("id = ? AND user_id = ?", userSkillID, userID).First(&userSkill).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User skill not found"})
		return
	}

	updates := gin.H{}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	updates["enabled"] = req.Enabled

	if err := db.Model(&userSkill).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_skill": userSkill,
		"message":    "配置更新成功",
	})
}

// 获取 Skill 分类
func GetSkillCategories(c *gin.Context, db *gorm.DB) {
	var categories []models.SkillCategory
	db.Order("sort_order ASC").Find(&categories)

	// 统计每个分类的数量
	for i := range categories {
		var count int64
		db.Model(&models.Skill{}).Where("category = ? AND enabled = ?", categories[i].Name, true).Count(&count)
		categories[i].Count = int(count)
	}

	c.JSON(http.StatusOK, gin.H{
		"categories": categories,
	})
}

// 点赞 Skill
func StarSkill(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	result := db.Model(&models.Skill{}).Where("id = ?").Update("stars", gorm.Expr("stars + 1"))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "点赞成功",
	})
}

// 初始化内置 Skills
func InitSkills(c *gin.Context, db *gorm.DB) {
	if err := models.InitBuiltinSkills(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "内置 Skills 初始化成功",
	})
}
