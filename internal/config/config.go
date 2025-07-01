package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"quickalias/internal/alias" // Alias paketinden persist fonksiyonlarına erişim için
	"quickalias/internal/ui"    // UI paketinden mesajlar ve renkler için
)

const (
	// CONFIG_FILE is the name of the configuration file.
	CONFIG_FILE = "config.json"
	// ALIASES_FILE is the name of the aliases file.
	ALIASES_FILE = "aliases.json"
	// GLOBAL_ALIASES_FILE is the name of the global aliases file.
	GLOBAL_ALIASES_FILE = "aliases.json"
	appName             = "quickalias"
)

// Config holds the application's configuration, including version, shell type, initialization status, and other settings.
type Config struct {
	Version     string            `json:"version"`
	ShellType   string            `json:"shell_type"`
	Initialized bool              `json:"initialized"`
	Settings    map[string]string `json:"settings"`
}

// SaveConfig writes the current application configuration to its JSON file.
func SaveConfig(configPath string, cfg *Config, errMsg string) error {
	fullConfigPath := filepath.Join(configPath, CONFIG_FILE)
	data, err := json.MarshalIndent(cfg, "", "  ") // Use 2 spaces for indentation
	if err != nil {
		return fmt.Errorf(errMsg, err)
	}

	if err := os.WriteFile(fullConfigPath, data, 0644); err != nil {
		return fmt.Errorf(errMsg, err)
	}

	return nil
}

// LoadConfig reads the application configuration file into the provided Config struct.
func LoadConfig(configPath string, cfg *Config) {
	fullConfigPath := filepath.Join(configPath, CONFIG_FILE)
	if data, err := os.ReadFile(fullConfigPath); err == nil {
		json.Unmarshal(data, cfg)
	}
}

// tryRemoveFile attempts to remove a file. If permission is denied, it automatically tries with sudo.
// Returns nil on success (or if file doesn't exist), or an error if deletion fails.
func tryRemoveFile(filePath, promptMsg string) error {
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, nothing to do.
			return nil
		}
		if os.IsPermission(err) {
			// Permission denied, automatically try with sudo without asking
			fmt.Printf("%s%s%s\n", ui.ColorYellow, promptMsg, ui.ColorReset)
			fmt.Printf("%s sudo rm \"%s\"%s\n", ui.ColorYellow, filePath, ui.ColorReset)

			// Directly try with sudo without user confirmation
			cmd := exec.Command("sudo", "-k", "rm", filePath)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf(ui.Msg.ErrorFileSudoRemovalFailed, filePath, err)
			}
			return nil // Sudo removal successful
		}
		// Other errors
		return fmt.Errorf(ui.Msg.ErrorFileRemovalFailed, filePath, err)
	}
	return nil // Original removal successful
}

// ResetConfig resets the application's configuration to its default state.
// It prompts for user confirmation before proceeding.
func ResetConfig(cfg *Config, version string) error {
	fmt.Printf("%s%s%s", ui.ColorYellow, ui.Msg.ResetConfigConfirmation, ui.ColorReset)
	var response string
	fmt.Scanln(&response)
	if strings.ToLower(response) != "e" && strings.ToLower(response) != "evet" && strings.ToLower(response) != "y" {
		fmt.Printf("%s❌ %s%s\n", ui.ColorRed, ui.Msg.ConfigResetCancelled, ui.ColorReset)
		return nil
	}

	// Get common config directory paths.
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf(ui.Msg.ErrorUserConfigDirNotFound, err)
	}
	appConfigDir := filepath.Join(userConfigDir, appName) // ~/.config/quickalias

	// User alias file path
	userAliasesFilePath := filepath.Join(appConfigDir, ALIASES_FILE)
	if err := tryRemoveFile(userAliasesFilePath, ui.Msg.UserAliasFileRemovePrompt); err != nil {
		return fmt.Errorf(ui.Msg.ErrorUserAliasFileReset, err)
	}

	// Global alias file path (using alias.GetGlobalConfigDir() from internal/alias/persist.go's package)
	// You need to ensure GetGlobalConfigDir() is accessible and returns the correct path (e.g. /etc/quickalias)
	globalConfigDirPath, err := alias.GetGlobalConfigDir()
	if err != nil {
		return fmt.Errorf(ui.Msg.ErrorGlobalConfigDirNotFound, err)
	}
	globalAliasesFilePath := filepath.Join(globalConfigDirPath, GLOBAL_ALIASES_FILE)

	if err := tryRemoveFile(globalAliasesFilePath, ui.Msg.GlobalAliasFileRemovePrompt); err != nil {
		return fmt.Errorf(ui.Msg.ErrorGlobalAliasFileReset, err)
	}

	// Main configuration file path
	configFilePath := filepath.Join(appConfigDir, CONFIG_FILE)
	if err := tryRemoveFile(configFilePath, ui.Msg.MainConfigFileRemovePrompt); err != nil {
		return fmt.Errorf(ui.Msg.ErrorMainConfigFileReset, err)
	}

	// Reset Config struct to default values in memory.
	*cfg = Config{
		Version:     version,
		ShellType:   "",
		Initialized: false,
		Settings:    make(map[string]string),
	}

	fmt.Printf("%s✅ %s%s\n", ui.ColorGreen, ui.Msg.ConfigResetSuccess, ui.ColorReset)
	return nil
}

//
