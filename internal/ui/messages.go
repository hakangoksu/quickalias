package ui

import (
	"os"
	"strings"
)

// Color codes for better UX
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

// messages holds all translatable strings for the CLI output.
type messages struct {
	AddAliasUsage                  string
	SetAliasUsage                  string
	RemoveAliasUsage               string
	UnsetAliasUsage                string
	SearchAliasUsage               string
	ConfigSubcommandRequired       string
	UnknownConfigSubcommand        string
	QuickAliasNotSetup             string
	RunSetupTip                    string
	UnknownCommand                 string
	ErrorInitializingQA            string
	ErrorCreatingUserConfigDir     string
	ErrorCreatingBackupDir         string
	ErrorGettingCurrentUser        string
	ErrorProcessingAliasData       string
	ErrorWritingAliasFile          string
	ErrorProcessingConfigData      string
	ErrorWritingConfigFile         string
	ErrorProcessingBackupData      string
	ErrorWritingBackupFile         string
	AccessDeniedGlobalAlias        string
	AdminPrivilegesNeeded          string
	RunCommandAsAdmin              string
	WarningAliasExists             string
	OverwritePrompt                string
	OperationCancelled             string
	AliasNotFound                  string
	AliasAddedSuccess              string
	UserAliasRemovedGlobalActive   string
	AliasRemovedSuccess            string
	GlobalAliasesHeader            string
	NoGlobalAliases                string
	UserAliasesHeader              string
	NoUserAliases                  string
	TotalAliasesFound              string
	SearchResults                  string
	NoResultsFound                 string
	TotalResultsFound              string
	QuickAliasStatus               string
	UserAliasesCount               string
	GlobalAliasesCount             string
	UserGlobalConflicts            string
	ConflictsHint                  string
	ShellIntegrationStatus         string
	StatusActive                   string
	StatusNotActive                string
	ConflictPrecedenceHint         string
	SetupStarting                  string
	ShellDetected                  string
	ErrorUnsupportedShell          string
	ShellIntegrationExists         string
	ShellConfigAccessError         string
	ShellConfigWriteError          string
	ShellIntegrationAdded          string
	AliasesLoading                 string
	InitFailedWarning              string
	RestartTerminalHint            string
	RestartTerminalCmdHint         string
	ResetConfigConfirmation        string
	BackupsNotFound                string
	AvailableBackups               string
	ExportDataProcessingError      string
	ExportFileWriteError           string
	ExportConfigSuccess            string
	ImportFileReadError            string
	ImportFileParseError           string
	ImportConfirmation             string
	ImportSuccess                  string
	ImportUserCount                string
	ImportGlobalCount              string
	AttemptingAsAdmin              string
	AliasAppliedToCurrentSession   string
	AliasRemovedFromCurrentSession string
	UsageTitle                     string
	UsageAliasManagement           string
	UsageListingSearching          string
	UsageSystem                    string
	UsageConfiguration             string
	UsageOther                     string
	TipsHeader                     string
	TipRunSetupFirst               string
	TipUserOverridesGlobal         string
	TipUseSudoGlobal               string
	// Config-specific messages
	ErrorUserConfigDirNotFound   string
	ErrorGlobalConfigDirNotFound string
	ErrorUserAliasFileReset      string
	ErrorGlobalAliasFileReset    string
	ErrorMainConfigFileReset     string
	ConfigResetSuccess           string
	ConfigResetCancelled         string
	UserAliasFileRemovePrompt    string
	GlobalAliasFileRemovePrompt  string
	MainConfigFileRemovePrompt   string
	ErrorFileSudoRemovalFailed   string
	ErrorFileRemovalFailed       string
}

var Msg *messages // Global variable to hold the current language messages, changed to capitalized for export.

// loadMessages populates the global 'Msg' variable with strings for the detected locale.
func LoadMessages(locale string) {
	if strings.Contains(strings.ToLower(locale), "tr") {
		Msg = loadTurkishMessages()
	} else {
		Msg = loadEnglishMessages()
	}
}

// init function in this package will detect locale and load messages when imported.
func init() {
	locale := os.Getenv("LC_ALL")
	if locale == "" {
		locale = os.Getenv("LANG")
	}
	LoadMessages(locale)
}

// These functions contain the localized message sets.
func loadTurkishMessages() *messages {
	return &messages{
		AddAliasUsage:                "Kullanım: qq add <alias> \"<komut>\"",
		SetAliasUsage:                "Kullanım: qq set <alias> \"<komut>\" (Global alias ekle)",
		RemoveAliasUsage:             "Kullanım: qq remove <alias>",
		UnsetAliasUsage:              "Kullanım: qq unset <alias> (Global alias kaldır)",
		SearchAliasUsage:             "Kullanım: qq search <anahtar_kelime>",
		ConfigSubcommandRequired:     "Yapılandırma komutu için alt komut gerekli.",
		UnknownConfigSubcommand:      "Bilinmeyen yapılandırma alt komutu: %s",
		QuickAliasNotSetup:           "QuickAlias kurulumu yapılmamış görünüyor!",
		RunSetupTip:                  "Lütfen QuickAlias'ı yapılandırmak için '%sqq setup%s' komutunu çalıştırın.",
		UnknownCommand:               "Bilinmeyen komut",
		ErrorInitializingQA:          "QuickAlias başlatılırken hata oluştu: %v",
		ErrorCreatingUserConfigDir:   "Kullanıcı yapılandırma dizini oluşturulurken hata oluştu: %w",
		ErrorCreatingBackupDir:       "Yedekleme dizini oluşturulurken hata oluştu: %w",
		ErrorGettingCurrentUser:      "Mevcut kullanıcı alınırken hata oluştu: %w",
		ErrorProcessingAliasData:     "Alias verileri işlenirken hata oluştu: %w",
		ErrorWritingAliasFile:        "Alias dosyasına yazılırken hata oluştu: %w",
		ErrorProcessingConfigData:    "Yapılandırma verileri işlenirken hata oluştu: %w",
		ErrorWritingConfigFile:       "Yapılandırma dosyasına yazılırken hata oluştu: %w",
		ErrorProcessingBackupData:    "Yedekleme verileri işlenirken hata oluştu: %w",
		ErrorWritingBackupFile:       "Yedekleme dosyasına yazılırken hata oluştu: %w",
		AccessDeniedGlobalAlias:      "Yetersiz yetki! Global aliasler için yönetici ayrıcalıkları gerekir.",
		AttemptingAsAdmin:            "Komut yönetici olarak deneniyor...",
		WarningAliasExists:           "⚠️ '%s' adında bir alias zaten mevcut (%s seviyesinde). Üzerine yazmak için 'e' veya 'y' yazın:",
		OperationCancelled:           "İşlem iptal edildi.",
		AliasNotFound:                "'%s' alias'ı %s seviyesinde bulunamadı.",
		AliasAddedSuccess:            "'%s' alias'ı (%s seviyesi) başarıyla eklendi/güncellendi.",
		UserAliasRemovedGlobalActive: "'%s' kullanıcı alias'ı kaldırıldı. Şimdi aktif olan global alias '%s' -> '%s'.",
		AliasRemovedSuccess:          "'%s' alias'ı (%s seviyesi) başarıyla kaldırıldı.",
		GlobalAliasesHeader:          "GLOBAL ALIASLAR:",
		NoGlobalAliases:              "Henüz global alias yok.",
		UserAliasesHeader:            "KULLANICI ALIASLARI:",
		NoUserAliases:                "Henüz kullanıcı alias yok.",
		TotalAliasesFound:            "Toplam %d alias bulundu.",
		SearchResults:                "Arama Sonuçları",
		NoResultsFound:               "'%s' için sonuç bulunamadı.",
		TotalResultsFound:            "Toplam %d sonuç bulundu.",
		QuickAliasStatus:             "QUICKALIAS DURUMU",
		UserAliasesCount:             " %sKullanıcı Aliasları:%s %d",
		GlobalAliasesCount:           " %sGlobal Aliaslar:%s %d",
		UserGlobalConflicts:          " %sKullanıcı-Global Çakışmaları:%s %d",
		ConflictsHint:                "(Çakışanlar: %s)",
		ShellIntegrationStatus:       " %sKabuk Entegrasyonu:%s %s",
		StatusActive:                 "Aktif ✅",
		StatusNotActive:              "Aktif Değil ❌",
		ConflictPrecedenceHint:       "💡 Not: Kullanıcı seviyesi alias'lar, global alias'ları geçersiz kılar.",
		SetupStarting:                "QUICKALIAS KURULUMU BAŞLATILIYOR...",
		ShellDetected:                "Algılanan kabuk tipi:",
		ErrorUnsupportedShell:        "Desteklenmeyen kabuk: %s",
		ShellIntegrationExists:       "⚠️ Kabuk entegrasyonu zaten mevcut.",
		ShellConfigAccessError:       "Kabuk yapılandırma dosyasına erişilemiyor: %w",
		ShellConfigWriteError:        "Kabuk yapılandırma dosyasına yazılamıyor: %w",
		ShellIntegrationAdded:        "✅ Kabuk entegrasyonu eklendi: %s",
		AliasesLoading:               "Aliaslar yükleniyor...",
		InitFailedWarning:            "Aliasları yüklerken hata oluştu (qq init): %v",
		RestartTerminalHint:          "QuickAlias ayarlarını geçerli kılmak için terminalinizi yeniden başlatmanız veya '%s. ~/.bashrc%s', '%s. ~/.zshrc%s' veya '%ssource ~/.config/fish/config.fish%s' komutunu çalıştırmanız gerekebilir.",
		RestartTerminalCmdHint:       "kaynak komutu çalıştırın",
		ResetConfigConfirmation:      "Yapılandırmayı sıfırlamak istediğinize emin misiniz? (Tüm ayarlar ve shell entegrasyon durumu sıfırlanır) [e/H]: ",
		BackupsNotFound:              "Hiç yedekleme bulunamadı.",
		AvailableBackups:             "MEVCUT YEDEKLEMELER:",
		ExportDataProcessingError:    "Dışa aktarma verileri işlenirken hata oluştu: %w",
		ExportFileWriteError:         "Dışa aktarma dosyasına yazılırken hata oluştu: %w",
		ExportConfigSuccess:          "Alias'lar başarıyla dışa aktarıldı",
		ImportFileReadError:          "İçe aktarma dosyası okunurken hata oluştu: %w",
		ImportFileParseError:         "İçe aktarma dosyası ayrıştırılırken hata oluştu: %w",
		ImportConfirmation:           "Bu işlem mevcut tüm kullanıcı ve global aliaslarınızı %d yeni alias ile DEĞİŞTİRECEKTİR. Devam etmek istiyor musunuz? [e/H]: ",
		ImportSuccess:                "Toplam %d alias başarıyla içe aktarıldı (Kullanıcı: %d, Global: %d).",
		ImportUserCount:              "Kullanıcı",
		ImportGlobalCount:            "Global",
		UsageTitle:                   "QUICKALIAS (qq) - Hızlı Alias Yönetimi",
		UsageAliasManagement:         "Alias Yönetimi:",
		UsageListingSearching:        "Listeleme ve Arama:",
		UsageSystem:                  "Sistem:",
		UsageConfiguration:           "Yapılandırma:",
		UsageOther:                   "Diğer:",
		TipsHeader:                   "İPUÇLARI:",
		TipRunSetupFirst:             "İlk çalıştırmada `qq setup` komutunu çalıştırın.",
		TipUserOverridesGlobal:       "Kullanıcı seviyesi alias'lar, aynı isimdeki global alias'ları geçersiz kılar.",
		TipUseSudoGlobal:             "Global alias'ları (`set`, `unset`) yönetmek için `sudo` kullanmanız gerekebilir.",
		// Config-specific messages
		ErrorUserConfigDirNotFound:   "kullanıcı yapılandırma dizini bulunamadı: %w",
		ErrorGlobalConfigDirNotFound: "global yapılandırma dizini bulunamadı: %w",
		ErrorUserAliasFileReset:      "kullanıcı alias dosyası sıfırlanamadı: %w",
		ErrorGlobalAliasFileReset:    "global alias dosyası sıfırlanamadı: %w",
		ErrorMainConfigFileReset:     "ana yapılandırma dosyası sıfırlanamadı: %w",
		ConfigResetSuccess:           "Yapılandırma sıfırlandı",
		ConfigResetCancelled:         "İşlem iptal edildi.",
		UserAliasFileRemovePrompt:    "Kullanıcı alias dosyasını silmek için sudo yetkisi gerekiyor:",
		GlobalAliasFileRemovePrompt:  "Global alias dosyasını silmek için sudo yetkisi gerekiyor:",
		MainConfigFileRemovePrompt:   "Ana yapılandırma dosyasını silmek için sudo yetkisi gerekiyor:",
		ErrorFileSudoRemovalFailed:   "dosya sudo ile silinemedi: %s: %w",
		ErrorFileRemovalFailed:       "dosya silinemedi: %s: %w",
	}
}

func loadEnglishMessages() *messages {
	return &messages{
		AddAliasUsage:                "Usage: qq add <alias> \"<command>\"",
		SetAliasUsage:                "Usage: qq set <alias> \"<command>\" (Add global alias)",
		RemoveAliasUsage:             "Usage: qq remove <alias>",
		UnsetAliasUsage:              "Usage: qq unset <alias> (Remove global alias)",
		SearchAliasUsage:             "Usage: qq search <keyword>",
		ConfigSubcommandRequired:     "Subcommand required for config command.",
		UnknownConfigSubcommand:      "Unknown config subcommand: %s",
		QuickAliasNotSetup:           "QuickAlias doesn't seem to be set up!",
		RunSetupTip:                  "Please run '%sqq setup%s' to configure QuickAlias.",
		UnknownCommand:               "Unknown command",
		ErrorInitializingQA:          "Error initializing QuickAlias: %v",
		ErrorCreatingUserConfigDir:   "Error creating user config directory: %w",
		ErrorCreatingBackupDir:       "Error creating backup directory: %w",
		ErrorGettingCurrentUser:      "Error getting current user: %w",
		ErrorProcessingAliasData:     "Error processing alias data: %w",
		ErrorWritingAliasFile:        "Error writing alias file: %w",
		ErrorProcessingConfigData:    "Error processing config data: %w",
		ErrorWritingConfigFile:       "Error writing config file: %w",
		ErrorProcessingBackupData:    "Error processing backup data: %w",
		ErrorWritingBackupFile:       "Error writing backup file: %w",
		AccessDeniedGlobalAlias:      "Insufficient permissions! Global aliases require administrator privileges.",
		AttemptingAsAdmin:            "Attempting command as administrator...",
		WarningAliasExists:           "⚠️ An alias named '%s' already exists at %s level. Type 'y' to overwrite: ",
		OperationCancelled:           "Operation cancelled.",
		AliasNotFound:                "Alias '%s' not found at %s level.",
		AliasAddedSuccess:            "Alias '%s' (%s level) successfully added/updated.",
		UserAliasRemovedGlobalActive: "User alias '%s' removed. Global alias '%s' -> '%s' is now active.",
		AliasRemovedSuccess:          "Alias '%s' (%s level) successfully removed.",
		GlobalAliasesHeader:          "GLOBAL ALIASES:",
		NoGlobalAliases:              "No global aliases yet.",
		UserAliasesHeader:            "USER ALIASES:",
		NoUserAliases:                "No user aliases yet.",
		TotalAliasesFound:            "Total %d aliases found.",
		SearchResults:                "Search Results",
		NoResultsFound:               "No results found for '%s'.",
		TotalResultsFound:            "Total %d results found.",
		QuickAliasStatus:             "QUICKALIAS STATUS",
		UserAliasesCount:             " %sUser Aliases:%s %d",
		GlobalAliasesCount:           " %sGlobal Aliases:%s %d",
		UserGlobalConflicts:          " %sUser-Global Conflicts:%s %d",
		ConflictsHint:                "(Conflicting: %s)",
		ShellIntegrationStatus:       " %sShell Integration:%s %s",
		StatusActive:                 "Active ✅",
		StatusNotActive:              "Not Active ❌",
		ConflictPrecedenceHint:       "💡 Note: User-level aliases take precedence over global aliases.",
		SetupStarting:                "STARTING QUICKALIAS SETUP...",
		ShellDetected:                "Detected shell type:",
		ErrorUnsupportedShell:        "Unsupported shell: %s",
		ShellIntegrationExists:       "⚠️ Shell integration already exists.",
		ShellConfigAccessError:       "Cannot access shell configuration file: %w",
		ShellConfigWriteError:        "Cannot write to shell configuration file: %w",
		ShellIntegrationAdded:        "✅ Shell integration added: %s",
		AliasesLoading:               "Loading aliases...",
		InitFailedWarning:            "Failed to load aliases (qq init): %v",
		RestartTerminalHint:          "To apply settings of QuickAlias, you may need to restart your terminal or run '%s. ~/.bashrc%s', '%s. ~/.zshrc%s', or '%ssource ~/.config/fish/config.fish%s'.",
		RestartTerminalCmdHint:       "source command",
		ResetConfigConfirmation:      "Are you sure you want to reset the configuration? (All settings and shell integration status will be reset) [y/N]: ",
		BackupsNotFound:              "No backups found.",
		AvailableBackups:             "AVAILABLE BACKUPS:",
		ExportDataProcessingError:    "Error processing export data: %w",
		ExportFileWriteError:         "Error writing export file: %w",
		ExportConfigSuccess:          "Aliases successfully exported",
		ImportFileReadError:          "Error reading import file: %w",
		ImportFileParseError:         "Error parsing import file: %w",
		ImportConfirmation:           "This operation will OVERWRITE all your existing user and global aliases with %d new aliases. Do you want to continue? [y/N]: ",
		ImportSuccess:                "Successfully imported %d aliases (User: %d, Global: %d).",
		ImportUserCount:              "User",
		ImportGlobalCount:            "Global",
		UsageTitle:                   "QUICKALIAS (qq) - Quick Alias Management",
		UsageAliasManagement:         "Alias Management:",
		UsageListingSearching:        "Listing and Searching:",
		UsageSystem:                  "System:",
		UsageConfiguration:           "Configuration:",
		UsageOther:                   "Other:",
		TipsHeader:                   "TIPS:",
		TipRunSetupFirst:             "Run `qq setup` first.",
		TipUserOverridesGlobal:       "User-level aliases override global aliases with the same name.",
		TipUseSudoGlobal:             "You might need to use `sudo` to manage global aliases (`set`, `unset`).",
		// Config-specific messages
		ErrorUserConfigDirNotFound:   "user config directory not found: %w",
		ErrorGlobalConfigDirNotFound: "global config directory not found: %w",
		ErrorUserAliasFileReset:      "user alias file could not be reset: %w",
		ErrorGlobalAliasFileReset:    "global alias file could not be reset: %w",
		ErrorMainConfigFileReset:     "main config file could not be reset: %w",
		ConfigResetSuccess:           "Configuration reset",
		ConfigResetCancelled:         "Operation cancelled.",
		UserAliasFileRemovePrompt:    "Sudo privileges required to remove user alias file:",
		GlobalAliasFileRemovePrompt:  "Sudo privileges required to remove global alias file:",
		MainConfigFileRemovePrompt:   "Sudo privileges required to remove main config file:",
		ErrorFileSudoRemovalFailed:   "file could not be removed with sudo: %s: %w",
		ErrorFileRemovalFailed:       "file could not be removed: %s: %w",
	}
}
