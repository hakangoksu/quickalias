package shell

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// QuickAliasConfig is an interface that defines the methods needed from the QuickAlias
// struct to perform shell-related operations. We're using an interface here to avoid
// a direct dependency on the main package's QuickAlias struct, promoting
// better separation of concerns.
type QuickAliasConfig interface {
	GetShellType() string
	SetShellType(string)
	IsInitialized() bool
	SaveConfig() error
}

// DetectShell attempts to determine the current shell based on the SHELL environment variable.
func DetectShell() string {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "bash") {
		return "bash"
	} else if strings.Contains(shell, "zsh") {
		return "zsh"
	} else if strings.Contains(shell, "fish") {
		return "fish"
	}
	return "" // Return empty string if shell is not recognized.
}

// AddShellIntegration adds a line to the shell's configuration file to source QuickAlias's init script.
// colorGreen, colorYellow, colorReset parametreleri dışarıdan alınacak.
func AddShellIntegration(qaConfig QuickAliasConfig, colorGreen, colorYellow, colorReset string) error {
	currentUser, _ := user.Current() // Get current user's home directory.
	var configFile string
	var integrationLine string

	shellType := qaConfig.GetShellType()
	if shellType == "" {
		shellType = DetectShell() // Detect shell if not already set.
		qaConfig.SetShellType(shellType)
	}

	// Determine the correct configuration file and integration line based on shell type.
	switch shellType {
	case "bash":
		configFile = filepath.Join(currentUser.HomeDir, ".bashrc")
		integrationLine = "eval \"$(qq init)\""
	case "zsh":
		configFile = filepath.Join(currentUser.HomeDir, ".zshrc")
		integrationLine = "eval \"$(qq init)\""
	case "fish":
		configDir := filepath.Join(currentUser.HomeDir, ".config/fish")
		os.MkdirAll(configDir, 0755) // Ensure Fish config directory exists.
		configFile = filepath.Join(configDir, "config.fish")
		integrationLine = "qq init | source" // Fish uses 'source' differently.
	default:
		return fmt.Errorf("Desteklenmeyen kabuk: %s", shellType)
	}

	// Check if the integration line already exists in the config file to prevent duplicates.
	if data, err := os.ReadFile(configFile); err == nil {
		if strings.Contains(string(data), integrationLine) {
			fmt.Printf("%s⚠️ Kabuk entegrasyonu zaten mevcut.%s\n", colorYellow, colorReset)
			return nil // Already integrated, no action needed.
		}
	}

	// Open the shell config file in append mode, creating it if it doesn't exist.
	file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Kabuk yapılandırma dosyasına erişilemiyor: %w", err)
	}
	defer file.Close() // Ensure the file is closed.

	// Write the integration line with a comment.
	_, err = file.WriteString(fmt.Sprintf("\n# QuickAlias entegrasyonu\n%s\n", integrationLine))
	if err != nil {
		return fmt.Errorf("Kabuk yapılandırma dosyasına yazılamıyor: %w", err)
	}

	fmt.Printf("%s✅ Kabuk entegrasyonu eklendi: %s%s\n", colorGreen, configFile, colorReset)
	return nil
}
