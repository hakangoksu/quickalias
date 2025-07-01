package alias

import (
	"fmt"
	"strings"

	"quickalias/internal/ui" // ui paketini import et
)

// Alias struct'ı persist.go'da tanımlandığı için burada tekrar tanımlamaya gerek yok.
// Alias is already defined in persist.go, no need to redefine here.
// But we need to make sure alias.go can access it, which it can by being in the same package.

// GetAlias checks if an alias exists and returns the alias itself and its level ("user" or "global").
func GetAlias(name string, userAliases, globalAliases []Alias) (*Alias, string) {
	for _, a := range userAliases {
		if a.Name == name {
			return &a, "user"
		}
	}
	for _, a := range globalAliases {
		if a.Name == name {
			return &a, "global"
		}
	}
	return nil, ""
}

// RemoveAlias is a helper function to remove an alias from a slice of aliases.
func RemoveAlias(name string, aliases []Alias) []Alias {
	for i, a := range aliases {
		if a.Name == name {
			return append(aliases[:i], aliases[i+1:]...)
		}
	}
	return aliases // If not found, return the original slice.
}

// ListAliases prints all user and global aliases, optionally filtered by a keyword.
func ListAliases(userAliases, globalAliases []Alias, keyword, globalHeader, userHeader, noGlobalMsg, noUserMsg, totalFoundMsg, colorPurpleBold, colorBlueBold, colorGreen, colorReset, colorYellow, colorCyan string) {
	fmt.Printf("%s%s%s\n", colorPurpleBold, globalHeader, colorReset)
	globalCount := 0
	for _, a := range globalAliases {
		if keyword == "" || strings.Contains(a.Name, keyword) || strings.Contains(a.Command, keyword) {
			fmt.Printf("  %s%s%s  → %s%s%s\n", colorGreen+ui.ColorBold, a.Name, colorReset, colorCyan, a.Command, colorReset) // ui.ColorBold kullanıldı
			globalCount++
		}
	}
	if globalCount == 0 {
		fmt.Printf("  %s%s%s\n", colorYellow, noGlobalMsg, colorReset)
	}

	fmt.Printf("\n%s%s%s\n", colorBlueBold, userHeader, colorReset)
	userCount := 0
	for _, a := range userAliases {
		if keyword == "" || strings.Contains(a.Name, keyword) || strings.Contains(a.Command, keyword) {
			fmt.Printf("  %s%s%s  → %s%s%s\n", colorGreen+ui.ColorBold, a.Name, colorReset, colorCyan, a.Command, colorReset) // ui.ColorBold kullanıldı
			userCount++
		}
	}
	if userCount == 0 {
		fmt.Printf("  %s%s%s\n", colorYellow, noUserMsg, colorReset)
	}

	if keyword != "" {
		fmt.Printf("\n%s%s%s\n", colorGreen, fmt.Sprintf(totalFoundMsg, globalCount+userCount), colorReset)
	}
}

// SearchAliases searches for aliases containing the given keyword in their name or command.
func SearchAliases(userAliases, globalAliases []Alias, keyword, searchResultsMsg, globalHeader, userHeader, noResultsMsg, totalResultsMsg, colorCyanBold, colorRed, colorGreen, colorReset, colorPurple, colorBlue, colorCyan string) {
	fmt.Printf("%s%s: '%s'%s\n", colorCyanBold, searchResultsMsg, keyword, colorReset)

	fmt.Printf("%s%s%s\n", colorPurple, globalHeader, colorReset)
	globalCount := 0
	for _, a := range globalAliases {
		if strings.Contains(a.Name, keyword) || strings.Contains(a.Command, keyword) {
			fmt.Printf("  %s%s%s  → %s%s%s\n", colorGreen+ui.ColorBold, a.Name, colorReset, colorCyan, a.Command, colorReset) // ui.ColorBold kullanıldı
			globalCount++
		}
	}

	fmt.Printf("\n%s%s%s\n", colorBlue, userHeader, colorReset)
	userCount := 0
	for _, a := range userAliases {
		if strings.Contains(a.Name, keyword) || strings.Contains(a.Command, keyword) {
			fmt.Printf("  %s%s%s  → %s%s%s\n", colorGreen+ui.ColorBold, a.Name, colorReset, colorCyan, a.Command, colorReset) // ui.ColorBold kullanıldı
			userCount++
		}
	}

	totalFound := globalCount + userCount
	if totalFound == 0 {
		fmt.Printf("\n%s❌ %s%s\n", colorRed, fmt.Sprintf(noResultsMsg, keyword), colorReset)
	} else {
		fmt.Printf("\n%s✅ %s%s\n", colorGreen, fmt.Sprintf(totalResultsMsg, totalFound), colorReset)
	}
}

// ShowStatus displays the current status of QuickAlias, including alias counts and conflicts.
func ShowStatus(userAliasCount, globalAliasCount int, conflicts []string, initialized bool, quickAliasStatusMsg, userAliasesCountMsg, globalAliasesCountMsg, userGlobalConflictsMsg, conflictsHintMsg, shellIntegrationStatusMsg, statusActiveMsg, statusNotActiveMsg, conflictPrecedenceHintMsg, colorCyanBold, colorBlue, colorPurple, colorGreen, colorYellow, colorReset, colorWhite string) error {
	fmt.Printf("%s%s%s\n", colorCyanBold, quickAliasStatusMsg, colorReset)
	fmt.Printf(userAliasesCountMsg+"\n", colorBlue, ui.ColorBold, userAliasCount, colorReset)       // ui.ColorBold kullanıldı
	fmt.Printf(globalAliasesCountMsg+"\n", colorPurple, ui.ColorBold, globalAliasCount, colorReset) // ui.ColorBold kullanıldı

	conflictColor := colorGreen
	if len(conflicts) > 0 {
		conflictColor = colorYellow // Change color if conflicts exist.
	}
	fmt.Printf(userGlobalConflictsMsg, conflictColor, ui.ColorBold, len(conflicts), colorReset) // ui.ColorBold kullanıldı
	if len(conflicts) > 0 {
		fmt.Printf(" %s%s%s", colorYellow, fmt.Sprintf(conflictsHintMsg, strings.Join(conflicts, ", ")), colorReset)
	}
	fmt.Println()

	statusText := statusNotActiveMsg // Default status is not active.
	if initialized {
		statusText = statusActiveMsg // Change to active if initialized.
	}
	fmt.Printf(shellIntegrationStatusMsg+"\n", colorWhite, statusText, colorReset) // colorReset parametresi eklendi.

	if len(conflicts) > 0 {
		fmt.Printf("\n%s%s%s\n", ui.ColorCyan, conflictPrecedenceHintMsg, ui.ColorReset) // ui.ColorCyan ve ui.ColorReset kullanıldı
	}

	return nil
}

// FindConflicts identifies aliases that exist at both user and global levels.
// User-level aliases take precedence over global ones.
func FindConflicts(userAliases, globalAliases []Alias) []string {
	conflicts := []string{}
	userAliasesMap := make(map[string]bool)

	// Populate a map with user alias names for quick lookup.
	for _, a := range userAliases {
		userAliasesMap[a.Name] = true
	}

	// Check if any global alias name exists in the user aliases map.
	for _, a := range globalAliases {
		if userAliasesMap[a.Name] {
			conflicts = append(conflicts, a.Name)
		}
	}

	return conflicts
}
