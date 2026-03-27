package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func loadUserPrompts() []prompt {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil
		}
		configHome = filepath.Join(home, ".config")
	}

	data, err := os.ReadFile(filepath.Join(configHome, "gtyper", "config.json"))
	if err != nil {
		return nil
	}

	var entries []struct {
		Content  string `json:"content"`
		Citation string `json:"citation"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil
	}

	var prompts []prompt
	for _, e := range entries {
		if strings.TrimSpace(e.Content) == "" {
			return nil
		}
		prompts = append(prompts, prompt{text: e.Content, source: e.Citation})
	}
	if len(prompts) == 0 {
		return nil
	}
	return prompts
}
