package ui

import "fmt"

// ShowUsage prints the command-line usage instructions for QuickAlias.
// It requires access to the global Msg and Color constants from the ui package.
func ShowUsage() {
	fmt.Printf("%s%s%s\n", ColorCyan+ColorBold, Msg.UsageTitle, ColorReset)
	fmt.Println()
	fmt.Printf("%sKULLANIM:%s\n", ColorGreen+ColorBold, ColorReset) // Bu kısmı da Msg'den almalısın
	fmt.Printf("  %s%s%s\n", ColorBlue, Msg.UsageAliasManagement, ColorReset)
	fmt.Printf("    %sqq add <alias> \"<komut>\"%s       %s\n", ColorWhite, ColorReset, "Kullanıcı seviye alias ekle") // Bu açıklama Msg'den gelmeli
	fmt.Printf("    %sqq set <alias> \"<komut>\"%s       %s\n", ColorWhite, ColorReset, "Global alias ekle (sudo gerekli)")
	fmt.Printf("    %sqq remove <alias>%s              %s\n", ColorWhite, ColorReset, "Kullanıcı alias kaldır")
	fmt.Printf("    %sqq unset <alias>%s               %s\n", ColorWhite, ColorReset, "Global alias kaldır (sudo gerekli)")
	fmt.Println()
	fmt.Printf("  %s%s%s\n", ColorBlue, Msg.UsageListingSearching, ColorReset)
	fmt.Printf("    %sqq list [anahtar_kelime]%s       %s\n", ColorWhite, ColorReset, "Tüm aliasları listele veya filtrele")
	fmt.Printf("    %sqq search <anahtar_kelime>%s     %s\n", ColorWhite, ColorReset, "Alias isim ve komutlarında ara")
	fmt.Println()
	fmt.Printf("  %s%s%s\n", ColorBlue, Msg.UsageSystem, ColorReset)
	fmt.Printf("    %sqq control%s                     %s\n", ColorWhite, ColorReset, "Durum ve çakışmaları göster")
	fmt.Printf("    %sqq setup%s                       %s\n", ColorWhite, ColorReset, "Shell entegrasyonunu kur")
	fmt.Printf("    %sqq init%s                        %s\n", ColorWhite, ColorReset, "Aliasları başlat (shell tarafından kullanılır)")
	fmt.Println()
	fmt.Printf("  %s%s%s\n", ColorBlue, Msg.UsageConfiguration, ColorReset)
	fmt.Printf("    %sqq config reset%s                %s\n", ColorWhite, ColorReset, "Yapılandırmayı sıfırla")
	fmt.Printf("    %sqq config backup%s               %s\n", ColorWhite, ColorReset, "Mevcut yedeklemeleri göster")
	fmt.Printf("    %sqq config export [yol]%s         %s\n", ColorWhite, ColorReset, "Aliasları dışa aktar")
	fmt.Printf("    %sqq config import <yol>%s         %s\n", ColorWhite, ColorReset, "Aliasları içe aktar")
	fmt.Println()
	fmt.Printf("  %s%s%s\n", ColorBlue, Msg.UsageOther, ColorReset)
	fmt.Printf("    %sqq version%s                     %s\n", ColorWhite, ColorReset, "Sürümü göster")
	fmt.Printf("    %sqq help%s                        %s\n", ColorWhite, ColorReset, "Bu yardımı göster")
	fmt.Println()
	fmt.Printf("%s%s%s\n", ColorCyan, Msg.TipsHeader, ColorReset)
	fmt.Printf("  • %s%s%s\n", ColorYellow, Msg.TipRunSetupFirst, ColorReset)
	fmt.Printf("  • %s%s%s\n", ColorYellow, Msg.TipUserOverridesGlobal, ColorReset)
	fmt.Printf("  • %s%s%s\n", ColorYellow, Msg.TipUseSudoGlobal, ColorReset)
}
