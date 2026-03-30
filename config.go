package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type userConfig struct {
	ScoreServer string       `json:"score_server"`
	Username    string       `json:"username"`
	Prompts     []userPrompt `json:"prompts"`
}

type userPrompt struct {
	Content  string `json:"content"`
	Citation string `json:"citation"`
}

func configPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "gtyper", "config.json"), nil
}

func generateUsername() string {
	b := make([]byte, 4)
	if _, err := rand.Read(b); err != nil {
		return "gtyper"
	}
	return hex.EncodeToString(b)
}

func loadUserConfig() *userConfig {
	path, err := configPath()
	if err != nil {
		return nil
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		cfg := &userConfig{Username: generateUsername()}
		if b, merr := json.MarshalIndent(cfg, "", "  "); merr == nil {
			_ = os.MkdirAll(filepath.Dir(path), 0o755)
			_ = os.WriteFile(path, b, 0o644)
		}
		return cfg
	}
	if err != nil {
		return nil
	}

	var cfg userConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	return &cfg
}
