🧩 QUICKALIAS (qq) - Quick Alias Management
Quickalias is a command-line tool designed to streamline your shell alias management. It provides:

User-level and global alias control

Easy listing & searching

Robust configuration options

📦 Installation & Quick Cleanup
To quickly try out quickalias and remove the files afterward:

bash
Kopyala
Düzenle
git clone https://github.com/hakangoksu/quickalias.git
cd quickalias
chmod +x install.sh
./install.sh
cd ..
rm -rf quickalias
❌ Uninstallation
To completely remove quickalias from your system:

bash
Kopyala
Düzenle
git clone https://github.com/hakangoksu/quickalias.git
cd quickalias
chmod +x install.sh
./install.sh --uninstall
cd ..
rm -rf quickalias
Alternatively, you can run qq uninstall after installation.

🚀 Usage
📁 Alias Management
bash
Kopyala
Düzenle
qq add "<name>=<command>"       # Add a user-level alias
qq set "<name>=<command>"       # Add a global alias (requires sudo)
qq remove <name>                # Remove a user-level alias
qq unset <name>                 # Remove a global alias (requires sudo)
📋 Listing & Searching
bash
Kopyala
Düzenle
qq list [keyword]               # List all aliases or filter by keyword
qq search <term>                # Search aliases by name or command
⚙️ System & Integration
bash
Kopyala
Düzenle
qq control                      # Display system status and detect conflicts
qq setup                        # Set up shell integration
qq init                         # Initialize aliases (used by the shell)
qq uninstall                    # Uninstall quickalias (same as install.sh --uninstall)
🔧 Configuration
bash
Kopyala
Düzenle
qq config reset                 # Reset configuration
qq config backup                # Show backup locations
qq config export [path]         # Export aliases to file
qq config import <file>         # Import aliases from file
ℹ️ Other
bash
Kopyala
Düzenle
qq version                      # Show current version
qq help                         # Display help message
🖥️ Compatibility
Currently supported on:

✅ Arch Linux

Support for other Linux distributions is under development.

💡 Tips
Run qq setup after installation to integrate with your shell.

User-level aliases override global aliases with the same name.

sudo may be required for managing global aliases (qq set, qq unset).
