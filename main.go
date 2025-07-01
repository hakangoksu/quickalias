package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"quickalias/internal/alias"
	"quickalias/internal/config"
	"quickalias/internal/shell"
	"quickalias/internal/ui" // ui paketini import et
)

const (
	VERSION           = "1.0.0"
	USER_CONFIG_DIR   = ".config/quickalias"
	GLOBAL_CONFIG_DIR = "/etc/quickalias"
)

// QuickAlias is the main struct that encapsulates the application's state and methods.
type QuickAlias struct {
	UserConfigPath   string
	GlobalConfigPath string
	UserAliases      []alias.Alias // alias.Alias struct'ını kullan
	GlobalAliases    []alias.Alias // alias.Alias struct'ını kullan
	Config           config.Config // config.Config struct'ını kullan
	PersistManager   *alias.PersistManager
}

// GetShellType implements the shell.QuickAliasConfig interface.
func (qa *QuickAlias) GetShellType() string {
	return qa.Config.ShellType
}

// SetShellType implements the shell.QuickAliasConfig interface.
func (qa *QuickAlias) SetShellType(shellType string) {
	qa.Config.ShellType = shellType
}

// IsInitialized implements the shell.QuickAliasConfig interface.
func (qa *QuickAlias) IsInitialized() bool {
	return qa.Config.Initialized
}

// SaveConfig implements the shell.QuickAliasConfig interface.
// This is needed for shell.AddShellIntegration, but actual saving logic is in internal/config.
func (qa *QuickAlias) SaveConfig() error {
	return config.SaveConfig(qa.UserConfigPath, &qa.Config, ui.Msg.ErrorWritingConfigFile)
}

func main() {
	// Check if any arguments are provided. If not, show usage.
	if len(os.Args) < 2 {
		ui.ShowUsage() // ui paketinden ShowUsage'ı çağır
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Handle `set` and `unset` commands with automatic sudo retry.
	if (command == "set" || command == "unset") && os.Geteuid() != 0 {
		fmt.Printf("%s%s%s\n", ui.ColorRed+ui.ColorBold, ui.Msg.AccessDeniedGlobalAlias, ui.ColorReset)
		fmt.Printf("%s💡 %s%s\n", ui.ColorCyan, ui.Msg.AttemptingAsAdmin, ui.ColorReset)

		exe, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ %s%s\n", ui.ColorRed, fmt.Errorf(ui.Msg.ErrorInitializingQA, err), ui.ColorReset)
			os.Exit(1)
		}

		sudoArgs := []string{exe}
		sudoArgs = append(sudoArgs, os.Args[1:]...)

		cmd := exec.Command("sudo", sudoArgs...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s❌ Hata: %v%s\n", ui.ColorRed, err, ui.ColorReset)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Initialize QuickAlias instance.
	qa, err := NewQuickAlias()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s❌ %s%s\n", ui.ColorRed, fmt.Errorf(ui.Msg.ErrorInitializingQA, err), ui.ColorReset)
		os.Exit(1)
	}

	skipInitCheck := []string{"setup", "init", "version", "help", "--help", "-h"}
	needsInit := true
	for _, cmd := range skipInitCheck {
		if command == cmd {
			needsInit = false
			break
		}
	}

	if needsInit && !qa.Config.Initialized {
		fmt.Printf("%s⚠️  %s%s\n", ui.ColorYellow+ui.ColorBold, ui.Msg.QuickAliasNotSetup, ui.ColorReset)
		fmt.Printf("%s💡 %s%s\n", ui.ColorCyan, fmt.Sprintf(ui.Msg.RunSetupTip, ui.ColorBold, ui.ColorReset), ui.ColorReset)
		os.Exit(1)
	}

	// Handle different commands based on user input.
	switch command {
	case "add":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.AddAliasUsage, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.AddAlias(args[0], strings.Join(args[1:], " "), "user")
	case "set":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.SetAliasUsage, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.AddAlias(args[0], strings.Join(args[1:], " "), "global")
	case "remove":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.RemoveAliasUsage, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.RemoveAlias(args[0], "user")
	case "unset":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.UnsetAliasUsage, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.RemoveAlias(args[0], "global")
	case "list":
		keyword := ""
		if len(args) > 0 {
			keyword = args[0]
		}
		err = qa.ListAliases(keyword)
	case "search":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.SearchAliasUsage, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.SearchAliases(args[0])
	case "control", "status":
		err = qa.ShowStatus()
	case "setup":
		err = qa.Setup()
	case "init":
		err = qa.Init()
	case "config":
		if len(args) < 1 {
			fmt.Fprintf(os.Stderr, "%s%s%s\n", ui.ColorRed, ui.Msg.ConfigSubcommandRequired, ui.ColorReset)
			os.Exit(1)
		}
		err = qa.HandleConfig(args)
	case "version":
		fmt.Printf("%sQuickAlias (qq) versiyon %s%s%s\n", ui.ColorGreen+ui.ColorBold, VERSION, ui.ColorReset, ui.ColorReset)
	case "help", "--help", "-h":
		ui.ShowUsage()
	default:
		fmt.Fprintf(os.Stderr, "%s%s: %s%s\n", ui.ColorRed, ui.Msg.UnknownCommand, command, ui.ColorReset)
		ui.ShowUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s❌ Hata: %v%s\n", ui.ColorRed, err, ui.ColorReset)
		os.Exit(1)
	}
}

// NewQuickAlias creates and initializes a new QuickAlias instance.
// It sets up configuration paths and loads existing aliases and configurations.
func NewQuickAlias() (*QuickAlias, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf(ui.Msg.ErrorGettingCurrentUser, err)
	}

	userConfigPath := filepath.Join(currentUser.HomeDir, USER_CONFIG_DIR)
	globalConfigPath := GLOBAL_CONFIG_DIR

	qa := &QuickAlias{
		UserConfigPath:   userConfigPath,
		GlobalConfigPath: globalConfigPath,
		UserAliases:      []alias.Alias{},
		GlobalAliases:    []alias.Alias{},
		Config: config.Config{ // config paketinden Config struct'ı
			Version:     VERSION,
			ShellType:   "",
			Initialized: false,
			Settings:    make(map[string]string),
		},
	}

	// Create user config directory if it doesn't exist.
	if err := os.MkdirAll(userConfigPath, 0755); err != nil {
		return nil, fmt.Errorf(ui.Msg.ErrorCreatingUserConfigDir, err)
	}

	// Create backup directory within user config if it doesn't exist.
	if err := os.MkdirAll(filepath.Join(userConfigPath, alias.BACKUP_DIR), 0755); err != nil { // alias.BACKUP_DIR kullan
		return nil, fmt.Errorf(ui.Msg.ErrorCreatingBackupDir, err)
	}

	// Initialize PersistManager
	qa.PersistManager = alias.NewPersistManager(qa.UserConfigPath, qa.GlobalConfigPath, &qa.UserAliases, &qa.GlobalAliases)

	// Load existing aliases and config.
	qa.PersistManager.LoadAliases()                  // PersistManager üzerinden çağır
	config.LoadConfig(qa.UserConfigPath, &qa.Config) // config paketinden çağır

	return qa, nil
}

// SaveAliases writes the current aliases (user or global) to their respective JSON files.
// This is now a wrapper for PersistManager.SaveAliases
func (qa *QuickAlias) SaveAliases(level string) error {
	return qa.PersistManager.SaveAliases(level, ui.Msg.ErrorProcessingAliasData, ui.Msg.ErrorWritingAliasFile, ui.Msg.ErrorCreatingUserConfigDir)
}

// AddAlias adds a new alias or updates an existing one at the specified level.
// It performs permission checks for global aliases and handles conflicts.
func (qa *QuickAlias) AddAlias(name, command, level string) error {
	// Check for existing alias and prompt for overwrite.
	existingAlias, existingLevel := alias.GetAlias(name, qa.UserAliases, qa.GlobalAliases) // alias paketinden GetAlias
	if existingAlias != nil && existingLevel != "" {
		fmt.Printf("%s⚠️  %s%s", ui.ColorYellow, fmt.Sprintf(ui.Msg.WarningAliasExists, name, existingLevel), ui.ColorReset)
		var response string
		fmt.Scanln(&response) // Read user input for confirmation.
		if strings.ToLower(response) != "e" && strings.ToLower(response) != "evet" && strings.ToLower(response) != "y" {
			fmt.Printf("%s❌ %s%s\n", ui.ColorRed, ui.Msg.OperationCancelled, ui.ColorReset)
			return nil
		}
	}

	// Create a backup before making changes.
	qa.PersistManager.CreateBackup(level, ui.Msg.ErrorProcessingBackupData, ui.Msg.ErrorWritingBackupFile)

	// Create the new alias struct.
	newAlias := alias.Alias{ // alias.Alias struct'ı
		Name:    name,
		Command: command,
		Created: time.Now().Format("2006-01-02 15:04:05"),
		Level:   level,
	}

	// Add or update the alias in the appropriate slice using alias package functions.
	if level == "user" {
		qa.UserAliases = alias.RemoveAlias(name, qa.UserAliases) // alias.RemoveAlias kullan
		qa.UserAliases = append(qa.UserAliases, newAlias)
	} else { // level == "global"
		qa.GlobalAliases = alias.RemoveAlias(name, qa.GlobalAliases) // alias.RemoveAlias kullan
		qa.GlobalAliases = append(qa.GlobalAliases, newAlias)
	}

	// Save the updated aliases to file.
	if err := qa.SaveAliases(level); err != nil {
		return err
	}

	fmt.Printf("%s✅ %s%s\n", ui.ColorGreen, fmt.Sprintf(ui.Msg.AliasAddedSuccess, name, level), ui.ColorReset)

	fmt.Printf("\n%s💡 %s %s%s%s\n", ui.ColorCyan, ui.Msg.RestartTerminalHint, ui.ColorBold, ui.Msg.RestartTerminalCmdHint, ui.ColorReset)
	return nil
}

// RemoveAlias removes an alias from the specified level.
// It handles cases where the alias is not found.
func (qa *QuickAlias) RemoveAlias(name, level string) error {
	// Attempt to remove the alias.
	var found bool
	if level == "user" {
		newAliases := alias.RemoveAlias(name, qa.UserAliases) // alias.RemoveAlias kullan
		found = len(newAliases) < len(qa.UserAliases)
		qa.UserAliases = newAliases
	} else { // level == "global"
		newAliases := alias.RemoveAlias(name, qa.GlobalAliases) // alias.RemoveAlias kullan
		found = len(newAliases) < len(qa.GlobalAliases)
		qa.GlobalAliases = newAliases
	}

	if !found {
		fmt.Printf("%s❌ %s%s\n", ui.ColorRed, fmt.Sprintf(ui.Msg.AliasNotFound, name, level), ui.ColorReset)
		return nil
	}

	// Save the updated aliases to file.
	if err := qa.SaveAliases(level); err != nil {
		return err
	}

	// Check if a global alias with the same name should now become active if a user alias was removed.
	var alternativeAlias *alias.Alias
	for _, a := range qa.GlobalAliases {
		if a.Name == name {
			alternativeAlias = &a
			break
		}
	}

	if alternativeAlias != nil {
		fmt.Printf("%s✅ %s%s\n",
			ui.ColorGreen, fmt.Sprintf(ui.Msg.UserAliasRemovedGlobalActive, name, name, alternativeAlias.Command), ui.ColorReset)
	} else {
		fmt.Printf("%s✅ %s%s\n", ui.ColorGreen, fmt.Sprintf(ui.Msg.AliasRemovedSuccess, name, level), ui.ColorReset)
	}

	fmt.Printf("\n%s💡 %s %s%s%s\n", ui.ColorCyan, ui.Msg.RestartTerminalHint, ui.ColorBold, ui.Msg.RestartTerminalCmdHint, ui.ColorReset)
	return nil
}

// AliasExists checks if an alias with the given name exists at either level.
func (qa *QuickAlias) AliasExists(name string) bool {
	_, level := alias.GetAlias(name, qa.UserAliases, qa.GlobalAliases) // alias.GetAlias kullan
	return level != ""
}

// ListAliases prints all user and global aliases, optionally filtered by a keyword.
func (qa *QuickAlias) ListAliases(keyword string) error {
	alias.ListAliases(qa.UserAliases, qa.GlobalAliases, keyword, ui.Msg.GlobalAliasesHeader, ui.Msg.UserAliasesHeader, ui.Msg.NoGlobalAliases, ui.Msg.NoUserAliases, ui.Msg.TotalAliasesFound, ui.ColorPurple+ui.ColorBold, ui.ColorBlue+ui.ColorBold, ui.ColorGreen, ui.ColorReset, ui.ColorYellow, ui.ColorCyan)
	return nil
}

// SearchAliases searches for aliases containing the given keyword in their name or command.
func (qa *QuickAlias) SearchAliases(keyword string) error {
	alias.SearchAliases(qa.UserAliases, qa.GlobalAliases, keyword, ui.Msg.SearchResults, ui.Msg.GlobalAliasesHeader, ui.Msg.UserAliasesHeader, ui.Msg.NoResultsFound, ui.Msg.TotalResultsFound, ui.ColorCyan+ui.ColorBold, ui.ColorRed, ui.ColorGreen, ui.ColorReset, ui.ColorPurple, ui.ColorBlue, ui.ColorCyan)
	return nil
}

// ShowStatus displays the current status of QuickAlias, including alias counts and conflicts.
func (qa *QuickAlias) ShowStatus() error {
	conflicts := alias.FindConflicts(qa.UserAliases, qa.GlobalAliases) // alias.FindConflicts kullan
	alias.ShowStatus(len(qa.UserAliases), len(qa.GlobalAliases), conflicts, qa.Config.Initialized, ui.Msg.QuickAliasStatus, ui.Msg.UserAliasesCount, ui.Msg.GlobalAliasesCount, ui.Msg.UserGlobalConflicts, ui.Msg.ConflictsHint, ui.Msg.ShellIntegrationStatus, ui.Msg.StatusActive, ui.Msg.StatusNotActive, ui.Msg.ConflictPrecedenceHint, ui.ColorCyan+ui.ColorBold, ui.ColorBlue, ui.ColorPurple, ui.ColorGreen, ui.ColorYellow, ui.ColorReset, ui.ColorWhite)
	return nil
}

// Setup initializes QuickAlias by detecting the shell and adding shell integration.
func (qa *QuickAlias) Setup() error {
	fmt.Printf("%s%s%s\n", ui.ColorCyan+ui.ColorBold, ui.Msg.SetupStarting, ui.ColorReset)

	detectedShell := shell.DetectShell() // shell paketinden DetectShell
	if detectedShell == "" {
		return fmt.Errorf(ui.Msg.ErrorUnsupportedShell, "unknown")
	}

	qa.Config.ShellType = detectedShell
	fmt.Printf("%s%s %s%s%s\n", ui.ColorGreen, ui.Msg.ShellDetected, ui.ColorBold, detectedShell, ui.ColorReset)

	// Add the necessary integration line to the shell's configuration file using the new package.
	if err := shell.AddShellIntegration(qa, ui.ColorGreen, ui.ColorYellow, ui.ColorReset); err != nil { // shell.AddShellIntegration kullan
		return err
	}

	qa.Config.Initialized = true // Mark as initialized.
	qa.SaveConfig()              // Save the updated configuration.

	fmt.Printf("%s✅ Kurulum tamamlandı!%s\n", ui.ColorGreen+ui.ColorBold, ui.ColorReset)

	// Automatically run init after setup to load aliases.
	fmt.Printf("%s%s%s\n", ui.ColorCyan, ui.Msg.AliasesLoading, ui.ColorReset)
	if err := qa.Init(); err != nil {
		fmt.Printf("%s⚠️  %s%s\n", ui.ColorYellow, fmt.Errorf(ui.Msg.InitFailedWarning, err), ui.ColorReset)
	}

	fmt.Printf("\n%s💡 %s %s%s%s\n", ui.ColorCyan, ui.Msg.RestartTerminalHint, ui.ColorBold, ui.Msg.RestartTerminalCmdHint, ui.ColorReset)
	return nil
}

// Init outputs alias commands suitable for evaluation by the shell.
// This function is typically called by the `eval "$(qq init)"` line in shell config.
func (qa *QuickAlias) Init() error {
	// Output global aliases first.
	for _, a := range qa.GlobalAliases {
		fmt.Printf("alias %s='%s'\n", a.Name, a.Command)
	}

	// Then output user aliases, which will override global aliases if names conflict.
	for _, a := range qa.UserAliases {
		fmt.Printf("alias %s='%s'\n", a.Name, a.Command)
	}

	return nil
}

// HandleConfig manages various configuration-related sub-commands.
func (qa *QuickAlias) HandleConfig(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(ui.Msg.ConfigSubcommandRequired)
	}

	switch args[0] {
	case "reset":
		return config.ResetConfig(&qa.Config, VERSION) // Sadece 2 parametre geçiliyor
	case "backup":
		return qa.PersistManager.ShowBackups(ui.ColorCyan+ui.ColorBold, ui.ColorGreen, ui.ColorReset, ui.ColorYellow, ui.Msg.AvailableBackups, ui.Msg.BackupsNotFound) // PersistManager.ShowBackups kullan
	case "export":
		exportPath := filepath.Join(os.Getenv("HOME"), "quickalias_export.json") // Default export path.
		if len(args) > 1 {
			exportPath = args[1] // User-specified export path.
			if !filepath.IsAbs(exportPath) {
				currentDir, _ := os.Getwd()
				exportPath = filepath.Join(currentDir, exportPath)
			}
		}
		return qa.PersistManager.ExportConfig(exportPath, ui.Msg.ExportDataProcessingError, ui.Msg.ExportFileWriteError, ui.Msg.ExportConfigSuccess, ui.ColorGreen, ui.ColorBold, ui.ColorReset) // PersistManager.ExportConfig kullan
	case "import":
		if len(args) < 2 {
			return fmt.Errorf(ui.Msg.ImportFileReadError, "path not provided")
		}
		return qa.PersistManager.ImportConfig(args[1], ui.ColorYellow, ui.ColorRed, ui.ColorGreen, ui.ColorReset, ui.Msg.ImportConfirmation, ui.Msg.OperationCancelled, ui.Msg.ImportFileReadError, ui.Msg.ImportFileParseError, ui.Msg.ImportSuccess) // PersistManager.ImportConfig kullan
	default:
		return fmt.Errorf(ui.Msg.UnknownConfigSubcommand, args[0])
	}
}
