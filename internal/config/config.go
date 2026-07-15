package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds server and provider credentials.
type Config struct {
	Addr string `yaml:"addr"`

	ZhipuAPIKey    string `yaml:"zhipu_api_key"`
	StepFunAPIKey  string `yaml:"stepfun_api_key"`
	DashScopeAPIKey string `yaml:"dashscope_api_key"`

	VolcengineAppID      string `yaml:"volcengine_app_id"`
	VolcengineAccessKey  string `yaml:"volcengine_access_key"`
	VolcengineResourceID string `yaml:"volcengine_resource_id"`
	VolcengineAppKey     string `yaml:"volcengine_app_key"`

	BailianWorkspaceID string `yaml:"bailian_workspace_id"`
	BailianAppID       string `yaml:"bailian_app_id"`
	BailianUserID      string `yaml:"bailian_user_id"`
}

// Load reads configuration from environment variables.
func Load() (*Config, error) {
	cfg := &Config{
		Addr:                 envOr("VOICE_REALTIME_ADDR", ":8080"),
		ZhipuAPIKey:          os.Getenv("ZHIPU_API_KEY"),
		StepFunAPIKey:        os.Getenv("STEPFUN_API_KEY"),
		DashScopeAPIKey:      os.Getenv("DASHSCOPE_API_KEY"),
		VolcengineAppID:      os.Getenv("VOLCENGINE_APP_ID"),
		VolcengineAccessKey:  os.Getenv("VOLCENGINE_ACCESS_KEY"),
		VolcengineResourceID: envOr("VOLCENGINE_RESOURCE_ID", "volc.speech.dialog"),
		VolcengineAppKey:     envOr("VOLCENGINE_APP_KEY", "PlgvMymc7f3tQnJ6"),
		BailianWorkspaceID:   os.Getenv("BAILIAN_WORKSPACE_ID"),
		BailianAppID:         os.Getenv("BAILIAN_APP_ID"),
		BailianUserID:        envOr("BAILIAN_USER_ID", "voice-realtime-user"),
	}
	return cfg, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Port returns the listen port parsed from Addr.
func (c *Config) Port() int {
	if c.Addr == "" || c.Addr[0] != ':' {
		return 8080
	}
	p, err := strconv.Atoi(c.Addr[1:])
	if err != nil {
		return 8080
	}
	return p
}

// ValidateProvider checks that required credentials exist for a provider.
func (c *Config) ValidateProvider(name string) error {
	switch name {
	case "zhipu":
		if c.ZhipuAPIKey == "" {
			return fmt.Errorf("ZHIPU_API_KEY is required for provider zhipu")
		}
	case "stepfun":
		if c.StepFunAPIKey == "" {
			return fmt.Errorf("STEPFUN_API_KEY is required for provider stepfun")
		}
	case "volcengine":
		if c.VolcengineAppID == "" || c.VolcengineAccessKey == "" {
			return fmt.Errorf("VOLCENGINE_APP_ID and VOLCENGINE_ACCESS_KEY are required for provider volcengine")
		}
	case "bailian":
		if c.DashScopeAPIKey == "" {
			return fmt.Errorf("DASHSCOPE_API_KEY is required for provider bailian")
		}
		if c.BailianWorkspaceID == "" || c.BailianAppID == "" {
			return fmt.Errorf("BAILIAN_WORKSPACE_ID and BAILIAN_APP_ID are required for provider bailian")
		}
	default:
		return fmt.Errorf("unknown provider %q", name)
	}
	return nil
}
