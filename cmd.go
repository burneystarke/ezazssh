package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func osCommand(cmd string) exec.Cmd {
	switch runtime.GOOS {
	case "windows":
		return *exec.Command("cmd", "/c", cmd)
	default:
		return *exec.Command("sh", "-c", cmd)
	}
}

func CallClear() {
	var clear string = "clear"
	if runtime.GOOS == "windows" {
		clear = "cls"
	}
	cmd := osCommand(clear)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// SaveConfig saves a key-value pair to the config file
func SaveConfig(key, value string) error {
	path := config.ConfigFile
	os.MkdirAll(config.ConfigPath, os.ModePerm)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	// Read existing lines
	scanner := bufio.NewScanner(file)
	lines := []string{}
	keyFound := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+"=") {
			lines = append(lines, fmt.Sprintf("%s=%s", key, value))
			keyFound = true
		} else {
			lines = append(lines, line)
		}
	}

	if !keyFound {
		lines = append(lines, fmt.Sprintf("%s=%s", key, value))
	}

	// Write updated lines back to file
	file.Seek(0, 0)
	file.Truncate(0)
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write config file: %w", err)
		}
	}
	return nil
}

// LoadConfig loads the value for a given key from the config file
func LoadConfig(key string) string {
	path := config.ConfigFile
	file, _ := os.Open(path)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, key+"=") {
			return strings.TrimPrefix(line, key+"=")
		}
	}
	return ""
}
