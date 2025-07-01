package alias

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort" // Dosyaları sıralamak için eklendi
	"strings"
	"time"
)

const (
	// ALIASES_FILE is the name of the aliases file.
	ALIASES_FILE = "aliases.json"
	// BACKUP_DIR is the directory name for backups.
	BACKUP_DIR = "backups"
	// MAX_BACKUPS is the maximum number of backup files to keep.
	MAX_BACKUPS = 5
)

// Alias represents a single alias entry with its name, command, creation date, and level.
type Alias struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Created string `json:"created"`
	Level   string `json:"level"` // "user" or "global"
}

// GetGlobalConfigDir returns the path to the global configuration directory.
// This is typically a system-wide directory like /etc/quickalias on Linux.
// It must be public (starts with a capital letter) to be accessible from other packages.
func GetGlobalConfigDir() (string, error) {
	// Bu kısım işletim sistemine ve uygulamanızın global kurulum stratejisine göre değişir.
	// Linux için /etc/quickalias yaygın bir yoldur.
	// Windows veya macOS için farklı yollar gerekebilir.
	// LÜTFEN: Global aliaslarınızın gerçekte nerede saklandığınıza göre buradaki yolu ayarlayın!
	globalConfigDir := "/etc/quickalias" // <-- BU SATIRI PROJENİZE GÖRE AYARLAYIN!

	// Dizin yoksa otomatik oluşturma (MkdirAll) mantığı genellikle burada olmaz
	// çünkü bu dizin sistem genelidir ve sudo yetkisi gerektirebilir.
	// Sadece yolu döndürüyoruz, dizinin varlığı ve izinleri tryRemoveFile içinde ele alınacak.

	return globalConfigDir, nil
}

// PersistManager handles loading, saving, backing up, importing, and exporting aliases.
type PersistManager struct {
	UserConfigPath   string
	GlobalConfigPath string
	UserAliases      *[]Alias // Pointer to QuickAlias's UserAliases
	GlobalAliases    *[]Alias // Pointer to QuickAlias's GlobalAliases
}

// NewPersistManager creates a new PersistManager instance.
func NewPersistManager(userConfigPath, globalConfigPath string, userAliases, globalAliases *[]Alias) *PersistManager {
	return &PersistManager{
		UserConfigPath:   userConfigPath,
		GlobalConfigPath: globalConfigPath,
		UserAliases:      userAliases,
		GlobalAliases:    globalAliases,
	}
}

// LoadAliases reads user and global alias files into the provided alias slices.
func (pm *PersistManager) LoadAliases() {
	// Load user aliases.
	userAliasPath := filepath.Join(pm.UserConfigPath, ALIASES_FILE)
	if data, err := os.ReadFile(userAliasPath); err == nil {
		json.Unmarshal(data, pm.UserAliases)
	}

	// Load global aliases.
	globalAliasPath := filepath.Join(pm.GlobalConfigPath, ALIASES_FILE)
	if data, err := os.ReadFile(globalAliasPath); err == nil {
		json.Unmarshal(data, pm.GlobalAliases)
	}
}

// SaveAliases writes the current aliases (user or global) to their respective JSON files.
func (pm *PersistManager) SaveAliases(level, errMsgProcess, errMsgWrite, errMsgCreateDir string) error {
	var aliases []Alias
	var configPath string

	if level == "user" {
		aliases = *pm.UserAliases
		configPath = filepath.Join(pm.UserConfigPath, ALIASES_FILE)
	} else { // level == "global"
		aliases = *pm.GlobalAliases
		configPath = filepath.Join(pm.GlobalConfigPath, ALIASES_FILE)

		// Create global config directory if it doesn't exist (required for global aliases).
		if err := os.MkdirAll(pm.GlobalConfigPath, 0755); err != nil {
			return fmt.Errorf(errMsgCreateDir, err)
		}
	}

	data, err := json.MarshalIndent(aliases, "", "  ") // Use 2 spaces for indentation
	if err != nil {
		return fmt.Errorf(errMsgProcess, err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf(errMsgWrite, err)
	}

	return nil
}

// CreateBackup creates a timestamped backup of the current aliases.
func (pm *PersistManager) CreateBackup(level, errMsgProcess, errMsgWrite string) error {
	timestamp := time.Now().Format("20060102_1504")
	backupName := fmt.Sprintf("backup_%s.json", timestamp)
	backupPath := filepath.Join(pm.UserConfigPath, BACKUP_DIR, backupName)

	var aliases []Alias
	if level == "user" {
		aliases = *pm.UserAliases
	} else { // level == "global"
		aliases = *pm.GlobalAliases
	}

	data, err := json.MarshalIndent(aliases, "", "  ") // Use 2 spaces for indentation
	if err != nil {
		return fmt.Errorf(errMsgProcess, err)
	}

	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf(errMsgWrite, err)
	}

	// Clean old backups to maintain a limited number of backups.
	pm.CleanOldBackups()

	return nil
}

// CleanOldBackups removes the oldest backup files if the total number exceeds MAX_BACKUPS.
func (pm *PersistManager) CleanOldBackups() {
	backupDir := filepath.Join(pm.UserConfigPath, BACKUP_DIR)
	files, err := filepath.Glob(filepath.Join(backupDir, "backup_*.json"))
	if err != nil {
		return // Silently return if there's an error listing files.
	}

	if len(files) <= MAX_BACKUPS {
		return // No need to clean if number of backups is within limit.
	}

	// Sort files by filename (which includes timestamp), oldest first.
	sort.Strings(files)

	// Remove oldest files to maintain MAX_BACKUPS.
	for i := 0; i < len(files)-MAX_BACKUPS; i++ {
		os.Remove(files[i]) // Ignore errors for cleanup.
	}
}

// ShowBackups lists all available backup files.
func (pm *PersistManager) ShowBackups(colorCyanBold, colorGreen, colorReset, colorYellow, msgAvailable, msgNotFound string) error {
	backupDir := filepath.Join(pm.UserConfigPath, BACKUP_DIR)
	files, err := filepath.Glob(filepath.Join(backupDir, "backup_*.json"))
	if err != nil {
		// Use a generic error message, as the original msg.AvailableBackups was designed for success case.
		return fmt.Errorf("Yedek dosyaları listelenirken hata oluştu: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("%s%s%s\n", colorYellow, msgNotFound, colorReset)
		return nil
	}

	fmt.Printf("%s%s%s\n", colorCyanBold, msgAvailable, colorReset)
	for i, file := range files {
		fmt.Printf("%s%d.%s %s%s%s\n", colorGreen, i+1, colorReset, colorCyanBold, filepath.Base(file), colorReset)
	}

	return nil
}

// ExportConfig exports all aliases (user and global) to a single JSON file.
func (pm *PersistManager) ExportConfig(path, errMsgProcess, errMsgWrite, successMsg, colorGreen, colorBold, colorReset string) error {
	// Combine user and global aliases into one slice for export.
	allAliases := append(*pm.GlobalAliases, *pm.UserAliases...)
	data, err := json.MarshalIndent(allAliases, "", "  ") // Use 2 spaces for indentation.
	if err != nil {
		return fmt.Errorf(errMsgProcess, err)
	}

	// Write the combined alias data to the specified file.
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf(errMsgWrite, err)
	}

	fmt.Printf("%s✅ %s: %s%s%s\n", colorGreen, successMsg, colorBold, path, colorReset)
	return nil
}

// ImportConfig imports aliases from a JSON file, separating them into user and global levels.
// It prompts for user confirmation and creates backups before importing.
func (pm *PersistManager) ImportConfig(path string, colorYellow, colorRed, colorGreen, colorReset, confirmMsg, cancelMsg, readErr, parseErr, importSuccessMsg string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf(readErr, err)
	}

	var aliases []Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		return fmt.Errorf(parseErr, err)
	}

	fmt.Printf("%s"+confirmMsg+"%s", colorYellow, len(aliases), colorReset)
	var response string
	fmt.Scanln(&response) // Read user input for confirmation.
	if strings.ToLower(response) != "e" && strings.ToLower(response) != "evet" && strings.ToLower(response) != "y" {
		fmt.Printf("%s❌ %s%s\n", colorRed, cancelMsg, colorReset)
		return nil
	}

	// Create backups of current aliases before overwriting.
	pm.CreateBackup("user", "", "")   // Error messages can be ignored for backup, as per original logic.
	pm.CreateBackup("global", "", "") // Error messages can be ignored for backup, as per original logic.

	userCount := 0
	globalCount := 0

	// Clear existing aliases before importing new ones.
	*pm.UserAliases = []Alias{}
	*pm.GlobalAliases = []Alias{}

	// Separate imported aliases into user and global categories.
	for _, alias := range aliases {
		if alias.Level == "global" {
			*pm.GlobalAliases = append(*pm.GlobalAliases, alias)
			globalCount++
		} else {
			*pm.UserAliases = append(*pm.UserAliases, alias)
			userCount++
		}
	}

	// Save the newly imported aliases.
	if err := pm.SaveAliases("user", "", "", ""); err != nil { // Error messages can be ignored here for simplicity or passed from main
		return err
	}
	if err := pm.SaveAliases("global", "", "", ""); err != nil { // Error messages can be ignored here for simplicity or passed from main
		return err
	}

	fmt.Printf("%s✅ %s%s\n", colorGreen,
		fmt.Sprintf(importSuccessMsg, len(aliases), userCount, globalCount, colorReset),
		colorReset)
	return nil
}
