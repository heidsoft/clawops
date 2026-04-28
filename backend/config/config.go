package config

import (
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Aliyun   AliyunConfig
	JWT      JWTConfig
	LLM     LLMConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	DSN string
}

type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	SecurityGroupID string
}

type JWTConfig struct {
	Secret string
	Expire int
}

type LLMConfig struct {
	Provider    string // openai / ollama
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float64
}

func Load() *Config {
	llmProvider := getEnv("LLM_PROVIDER", "openai")
	llmBaseURL := getEnv("LLM_BASE_URL", "")
	if llmBaseURL == "" {
		if llmProvider == "ollama" {
			llmBaseURL = "http://localhost:11434"
		} else {
			llmBaseURL = "https://api.openai.com/v1"
		}
	}
	
	llmModel := getEnv("LLM_MODEL", "")
	if llmModel == "" {
		if llmProvider == "ollama" {
			llmModel = "llama3.2"
		} else {
			llmModel = "gpt-4o-mini"
		}
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Mode: getEnv("GIN_MODE", "release"),
		},
		Database: DatabaseConfig{
			DSN: getEnv("DB_DSN", "host=localhost user=postgres password=password dbname=openclaw_deploy port=5432 sslmode=disable"),
		},
		Aliyun: AliyunConfig{
			AccessKeyID:     os.Getenv("ALIYUN_ACCESS_KEY_ID"),
			AccessKeySecret: os.Getenv("ALIYUN_ACCESS_KEY_SECRET"),
			RegionID:        getEnv("ALIYUN_REGION_ID", "cn-shanghai"),
			SecurityGroupID: os.Getenv("SECURITY_GROUP_ID"),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "clawops-secret-key"),
			Expire: 24, // hours
		},
		LLM: LLMConfig{
			Provider:    llmProvider,
			APIKey:     os.Getenv("OPENAI_API_KEY"),
			BaseURL:    llmBaseURL,
			Model:     llmModel,
			Temperature: 0.7,
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var globalConfig *Config

func Init() {
	globalConfig = Load()
}

func Get() *Config {
	return globalConfig
}