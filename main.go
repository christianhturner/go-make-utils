package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	templateFile = "local-config.template.json"
	configFile   = "local-config.json"
	envFile      = ".env.mk"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: config-tool ensure [project-root]")
		os.Exit(1)
	}

	// Replace with switch if more commands are added
	if os.Args[1] == "ensure" {
		projectRoot := os.Args[2]
		if err := ensureConfig(projectRoot); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("invalid command: \"%s\"", os.Args[1])
		os.Exit(1)
	}
}

func ensureConfig(projectRoot string) error {
	templatePath := filepath.Join(projectRoot, templateFile)
	configPath := filepath.Join(projectRoot, configFile)
	envPath := filepath.Join(projectRoot, envFile)

	template, err := loadJSON(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template: %w", err)
	}

	config, configExists := loadJSONOrEmpty(configPath)

	if !configExists {
		fmt.Println("🔧 First time setup - please provide values for configuration:")
		config = make(map[string]any)
		for key := range template {
			config[key] = promptForValue(key, template[key])
		}
	} else {
		needsUpdate := false

		for key := range template {
			if _, exists := config[key]; !exists {
				fmt.Printf("🆕 New configuration key detected: %s\n", key)
				config[key] = promptForValue(key, template[key])
				needsUpdate = true
			}
		}

		// remove keys when deleted from template
		for key := range config {
			if _, exists := template[key]; !exists {
				fmt.Printf("🗑️Removing obsolete key :%s\n", key)
				delete(config, key)
				needsUpdate = true
			}
		}

		if needsUpdate {
			fmt.Println("✅ Configuration updated")
		}
	}

	if err := saveJSON(configFile, config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	if err := generateEnvMakeFile(config, envPath); err != nil {
		return fmt.Errorf("failed to generate env makefile: %w", err)
	}

	return nil
}

func loadJSON(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func loadJSONOrEmpty(path string) (map[string]any, bool) {
	result, err := loadJSON(path)
	if err != nil {
		return make(map[string]any), false
	}
	return result, true
}

func saveJSON(path string, data map[string]any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0o644)
}

func promptForValue(key string, templateValue any) string {
	reader := bufio.NewReader(os.Stdin)

	hint := ""
	if templateValue != nil {
		hint = fmt.Sprintf(" [%v]", templateValue)
	}

	fmt.Printf("  %s%s: ", key, hint)

	value, _ := reader.ReadString('\n')
	value = strings.TrimSpace(value)

	if value == "" && templateValue != nil {
		if strVal, ok := templateValue.(string); ok {
			return strVal
		}
	}

	return value
}

func generateEnvMakeFile(config map[string]any, envPath string) error {
	var lines []string
	lines = append(lines, "# Auto-generated  local-config.json - do not edit manually")

	for key, value := range config {
		envKey := strings.ToUpper(key)

		var envValue string
		switch v := value.(type) {
		case string:
			envValue = v
		default:
			envValue = fmt.Sprintf("%v", v)
		}

		envValue = strings.ReplaceAll(envValue, "$", "$$")

		lines = append(lines, fmt.Sprintf("export %s := %s", envKey, envValue))
	}

	content := strings.Join(lines, "\n") + "\n"

	return os.WriteFile(envPath, []byte(content), 0o644)
}
