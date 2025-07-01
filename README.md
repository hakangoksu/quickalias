Harika, şimdi kaldırma (uninstall) seçeneğini ekleyip, Arch Linux desteği ile ilgili notları da içeren güncellenmiş bir README hazırlayalım.

QUICKALIAS (qq) - Quick Alias Management
QUICKALIAS (or qq) is a command-line tool designed to streamline your shell alias management, offering both user-level and global alias control, easy listing, searching, and robust configuration options.

Installation & Quick Cleanup
To quickly try quickalias and then clean up the downloaded files, use the following command sequence:

Bash

git clone https://github.com/hakangoksu/quickalias.git
cd quickalias
chmod +x install.sh
./install.sh
cd ..
rm -rf quickalias
Uninstallation
To remove quickalias from your system, navigate into the cloned directory and run the install.sh script with the --uninstall flag:

Bash

git clone https://github.com/hakangoksu/quickalias.git
cd quickalias
chmod +x install.sh
./install.sh --uninstall
cd ..
rm -rf quickalias
Usage
Alias Management
qq add <alias> "<command>": Add a user-level alias.

qq set <alias> "<command>": Add a global alias (requires sudo).

qq remove <alias>: Remove a user-level alias.

qq unset <alias>: Remove a global alias (requires sudo).

Listing and Searching
qq list [keyword]: List all aliases or filter by a keyword.

qq search <keyword>: Search for aliases by name or command.

System
qq control: Display system status and potential conflicts.

qq setup: Set up shell integration.

qq init: Initialize aliases (used by the shell).

qq uninstall: Uninstall quickalias from your system. (This also performs the uninstall, same as install.sh --uninstall).

Configuration
qq config reset: Reset configuration.

qq config backup: Show existing backups.

qq config export [path]: Export aliases to a file.

qq config import <path>: Import aliases from a file.

Other
qq version: Display the tool version.

qq help: Show this help message.

Compatibility
Currently, quickalias officially supports Arch Linux. We are actively working on expanding support to other operating systems in future releases.

Tips
Run qq setup first to integrate qq with your shell.

User-level aliases take precedence over global aliases with the same name.

Managing global aliases (set, unset) may require sudo privileges.
